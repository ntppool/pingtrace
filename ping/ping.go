package ping

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/abh/globe/cmdparser"
	"github.com/kr/pretty"
)

// ping -t 4 -i 0.2 -c 10 some-ip

type PingResponse struct {
	Ping string
	// Seq     int
	// Size    int
	// IP      string
	// ASN     int
	// TTL     int
	// Latency Latency
	err error `json:"Error,omitempty"`
}

// --- 8.8.8.8 ping statistics ---
// 58 packets transmitted, 58 packets received, 0.0% packet loss
// round-trip min/avg/max/stddev = 34.750/40.406/70.660/5.502 ms

// --- 10.0.200.141 ping statistics ---
// 10 packets transmitted, 10 packets received, 0.0% packet loss
// round-trip min/avg/max/stddev = 0.998/2.143/4.191/1.068 ms

// 11 packets transmitted, 0 received, +7 errors, 100% packet loss, time 10215ms

type PingSummary struct {
	IP              net.IP
	Sent            int
	Received        int
	Errors          int
	Loss            float64
	RoundTripMin    float64
	RoundTripAvg    float64
	RoundTripMax    float64
	RoundTripStdDev float64
	err             error `json:"Error,omitempty"`
}

type PingParser struct {
	out       chan cmdparser.ParserOutput
	in        chan string
	inSummary bool
	sum       *PingSummary
}

func NewPingParser() *PingParser {
	pp := new(PingParser)
	pp.in = make(chan string, 10)
	pp.out = make(chan cmdparser.ParserOutput, 10)
	go pp.run()
	return pp
}

func (pp *PingParser) Add(line string) {
	pp.in <- line
}

func (pp *PingParser) Close() {
	close(pp.in)
}

func (pp *PingParser) Read() cmdparser.ParserOutput {
	pr, ok := <-pp.out
	if !ok {
		return nil
	}
	return pr
}

func (pr *PingResponse) JSON() []byte {
	b, _ := json.Marshal(pr)
	return b
}

func (pr *PingResponse) Error() error {
	return pr.err
}

func (pr *PingResponse) Bytes() []byte {
	// var b bytes.Buffer
	// fmt.Fprintf(&b, "%d bytes", pr.Size)
	// fmt.Fprintf(&b, " from %s", pr.IP)
	// if pr.ASN > 0 {
	// 	fmt.Fprintf(&b, " (AS%d)", pr.ASN)
	// }
	// fmt.Fprint(&b, "  ", pr.Latency.String())

	return []byte(pr.Ping)
}

func (pr *PingResponse) String() string {
	return pr.Ping
}

func (ps *PingSummary) Error() error {
	return ps.err
}

func (ps *PingSummary) Bytes() []byte {

	// --- 8.8.8.8 ping statistics ---
	// 58 packets transmitted, 58 packets received, 0.0% packet loss
	// round-trip min/avg/max/stddev = 34.750/40.406/70.660/5.502 ms

	// --- 10.0.200.141 ping statistics ---
	// 10 packets transmitted, 10 packets received, 0.0% packet loss
	// round-trip min/avg/max/stddev = 0.998/2.143/4.191/1.068 ms

	// 11 packets transmitted, 0 received, +7 errors, 100% packet loss, time 10215ms

	var b bytes.Buffer
	fmt.Fprintf(&b, "--- %s ping statistics ---\n", ps.IP)
	fmt.Fprintf(&b, "%d packets transmitted, %d packets received",
		ps.Sent, ps.Received)

	if ps.Errors > 0 {
		fmt.Fprintf(&b, " (+%d errors)", ps.Errors)
	}

	fmt.Fprintf(&b, ", %.1f%% packet loss\n", ps.Loss)

	if ps.Received > 0 {
		fmt.Fprintf(&b, "round-trip min/avg/max/stddev = %.3f/%.3f/%.3f/%.3f ms\n",
			ps.RoundTripMin, ps.RoundTripAvg, ps.RoundTripMax, ps.RoundTripStdDev,
		)
	}

	return b.Bytes()
}

func (ps *PingSummary) String() string {
	b := ps.Bytes()
	return string(b)
}

func (ps *PingSummary) JSON() []byte {
	b, _ := json.Marshal(ps)
	return b
}

func (pp *PingParser) run() {
	defer close(pp.out)
	for {
		line, ok := <-pp.in
		if !ok {
			if pp.sum != nil {
				pp.out <- pp.sum
			}
			break
		}
		err := pp.parseLine(line)
		if err != nil {
			pp.out <- &PingResponse{err: err}
		}
	}
}

func (pp *PingParser) parseLine(line string) error {
	// "64 bytes from 8.8.8.8: icmp_seq=2 ttl=46 time=34.540 ms",
	// "From 64.233.174.41 icmp_seq=8 Time to live exceeded",
	// "Request timeout for icmp_seq 387",

	if !pp.inSummary {
		if strings.HasPrefix(line, "---") {
			pp.inSummary = true
			pp.sum = new(PingSummary)
			s := strings.SplitN(line, " ", 3)
			pp.sum.IP = net.ParseIP(s[1])
			return nil
		}
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			pp.out <- &PingResponse{line, nil}
		}
		return nil
	}

	// --- 8.8.8.8 ping statistics ---
	// 58 packets transmitted, 58 packets received, 0.0% packet loss
	// round-trip min/avg/max/stddev = 34.750/40.406/70.660/5.502 ms

	// --- 10.0.200.141 ping statistics ---
	// 10 packets transmitted, 10 packets received, 0.0% packet loss
	// round-trip min/avg/max/stddev = 0.998/2.143/4.191/1.068 ms

	// 11 packets transmitted, 0 received, +7 errors, 100% packet loss, time 10215ms

	p := strings.Fields(line)
	log.Printf("%#v", p)

	s := pp.sum

	var err error

	if len(p) >= 9 {
		if p[1] == "packets" {
			s.Sent, err = strconv.Atoi(p[0])
			if err != nil {
				return err
			}
			s.Received, err = strconv.Atoi(p[3])
			if err != nil {
				return err
			}

			var pctIdx int
			switch len(p) {
			case 9:
				// 10 packets transmitted, 10 packets received, 0.0% packet loss
				pctIdx = 6
			case 10:
				// "5 packets transmitted, 5 received, 0% packet loss, time 1231ms",
				pctIdx = 5
			case 12:
				// 11 packets transmitted, 0 received, +7 errors, 100% packet loss, time 10215ms
				pctIdx = 7
				s.Errors, err = strconv.Atoi(p[5])
				if err != nil {
					return err
				}
			}

			s.Loss, err = strconv.ParseFloat(strings.TrimSuffix(p[pctIdx], "%"), 64)
			if err != nil {
				return err
			}
		}
	}

	if len(p) == 5 && (p[0] == "round-trip" || p[0] == "rtt") {
		// round-trip min/avg/max/stddev = 0.998/2.143/4.191/1.068 ms
		// "rtt min/avg/max/mdev = 26.768/26.783/26.793/0.146 ms",
		times := strings.Split(p[3], "/")
		for i, t := range times {
			ms, err := strconv.ParseFloat(t, 64)
			if err != nil {
				return err
			}
			switch i {
			case 0:
				s.RoundTripMin = ms
			case 1:
				s.RoundTripAvg = ms
			case 2:
				s.RoundTripMax = ms
			case 3:
				s.RoundTripStdDev = ms
			}
		}
	}

	for i, f := range p {
		pretty.Println(i, f)
	}

	return nil
}
