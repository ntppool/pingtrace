package traceroute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"strings"
	"sync"

	"go.ntppool.org/pingtrace/cmdparser"
	"go.ntppool.org/pingtrace/netinfo"
)

type TraceRouteResult struct {
	SourceIP netip.Addr
	TargetIP netip.Addr
	Lines    []*TraceRouteLine
}

type TraceRouteLine struct {
	Hop     int
	Name    string
	IP      string
	ASN     int
	Latency []Latency
	Raw     string
	Err     error `json:",omitempty"`
}

type TracerouteParser struct {
	out chan TraceRouteLine
	in  chan string
	hop int
}

func NewTracerouteParser() *TracerouteParser {
	trp := new(TracerouteParser)
	trp.in = make(chan string, 10)
	trp.out = make(chan TraceRouteLine, 10)
	go trp.run()
	return trp
}

func (trp *TracerouteParser) Add(line string) {
	trp.in <- line
}

func (trp *TracerouteParser) Close() {
	close(trp.in)
}

func (trp *TracerouteParser) Read() cmdparser.ParserOutput {
	tr, ok := <-trp.out
	if !ok {
		return nil
	}
	return &tr
}

func (tr *TraceRouteLine) JSON() []byte {
	b, _ := json.Marshal(tr)
	return b
}

func (tr *TraceRouteLine) Bytes() []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%2d", tr.Hop)
	if len(tr.Name) > 0 {
		fmt.Fprintf(&b, " %s", tr.Name)
	}
	if len(tr.IP) > 0 {
		fmt.Fprintf(&b, " (%s)", tr.IP)
	}
	if tr.ASN > 0 {
		fmt.Fprintf(&b, " AS%d", tr.ASN)
	}
	for _, ms := range tr.Latency {
		fmt.Fprint(&b, "  ", ms.String())
	}
	if tr.Err != nil {
		fmt.Fprint(&b, tr.Err)
	}
	return b.Bytes()
}

func (tr *TraceRouteLine) String() string {
	b := tr.Bytes()
	return string(b)
}

func (tr *TraceRouteLine) Error() error {
	return tr.Err
}

func (trp *TracerouteParser) run() {
	defer close(trp.out)
	for {
		line, ok := <-trp.in
		if !ok {
			break
		}
		err := trp.parseLine(line)
		if err != nil {
			trp.out <- TraceRouteLine{Err: err, Raw: line}
		}
	}
}

func (trp *TracerouteParser) parseLine(line string) error {
	// "1  ns.unitedwifi.com (172.19.248.1)  3.353 ms  0.959 ms  0.858 ms"

	if strings.HasPrefix(line, "traceroute to ") {
		// todo: return or parse this as the status/options or something
		return nil
	}

	p := strings.Fields(line)
	// log.Printf("%#v", p)

	var err error

	tr := TraceRouteLine{Raw: line}

	index := 1

	for i, s := range p {

		if i == 0 {
			// is a domain name or IP
			if strings.Contains(s, ".") {
				index = 0
			} else {
				tr.Hop, err = strconv.Atoi(s)
				if err != nil {
					return err
				}
				trp.hop = tr.Hop
				continue
			}
		}
		if tr.Hop == 0 {
			tr.Hop = trp.hop
		}

		if s == "*" {
			tr.Latency = append(tr.Latency, Latency{0, "*"})
			continue
		}

		if i == index {
			// the next entry is the IP
			if p[index+1][0] == 'b' {
				tr.Name = s
			} else {
				tr.IP = s
			}

			continue
		}

		// we got a name before, this is the IP
		if i == index+1 && len(tr.IP) == 0 {
			// tr.IP = s[1 : len(s)-1]
			tr.IP = strings.Trim(s, "()")
			if tr.IP == tr.Name {
				tr.Name = ""
			}
			continue
		}

		if s == "ms" {
			continue
		}

		if strings.HasPrefix(s, "!") {
			// pretty.Println(tr)
			tr.Latency[len(tr.Latency)-1].Error = s
		} else if ip := net.ParseIP(s); ip != nil {
			// is another IP address
			tr.UpdateASN()
			trp.out <- tr
			tr.IP = s
			index = i
			tr.Latency = []Latency{}
		} else {
			ms, err := strconv.ParseFloat(s, 64)
			if err != nil {
				// maybe it's a hostname, so we run traceroute with -n to avoid it...
				return err
			}
			tr.Latency = append(tr.Latency, Latency{ms, ""})
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		tr.UpdateASN()
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		tr.UpdateName()
		wg.Done()
	}()
	wg.Wait()
	trp.out <- tr

	return nil
}

func (tr *TraceRouteLine) UpdateASN() error {
	if tr.ASN != 0 {
		return nil
	}
	asn, err := netinfo.GetASN(tr.IP)
	if err != nil {
		return err
	}
	tr.ASN = asn
	return nil
}

func (tr *TraceRouteLine) UpdateName() error {
	if len(tr.Name) > 0 {
		return nil
	}
	names, err := netinfo.GetNames(tr.IP)
	if err != nil {
		return err
	}
	if len(names) > 0 {
		tr.Name = names[0]
	}
	return nil
}
