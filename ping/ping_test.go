package ping

import (
	"fmt"
	"net"
	"testing"

	"github.com/abh/pingtrace/cmdparser"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
)

func TestParsePingLine(t *testing.T) {

	testLines := [][]string{
		[]string{
			"64 bytes from 8.8.8.8: icmp_seq=2 ttl=46 time=34.540 ms",
			"From 64.233.174.41 icmp_seq=8 Time to live exceeded",
			"Request timeout for icmp_seq 387",
		},
		[]string{
			"--- 8.8.8.8 ping statistics ---",
			"58 packets transmitted, 58 packets received, 0.0% packet loss",
			"round-trip min/avg/max/stddev = 34.750/40.406/70.660/5.502 ms",
		},

		[]string{
			"--- 192.0.2.1 ping statistics ---",
			"11 packets transmitted, 0 received, +7 errors, 100% packet loss, time 10215ms",
		},

		[]string{
			"64 bytes from 8.8.4.4: icmp_seq=4 ttl=52 time=26.7 ms",
			"64 bytes from 8.8.4.4: icmp_seq=5 ttl=52 time=26.7 ms",
			"",
			"--- 8.8.4.4 ping statistics ---",
			"5 packets transmitted, 5 received, 0% packet loss, time 1231ms",
			"rtt min/avg/max/mdev = 26.768/26.783/26.793/0.146 ms",
		},

		[]string{
			"--- 4.4.2.2 ping statistics ---",
			"10 packets transmitted, 0 received, 100% packet loss, time 3000ms",
			"",
		},
	}
	results := []cmdparser.ParserOutput{
		// &PingResponse{2, 64, "8.8.8.8", 15169, 46, Latency{34.540, ""}, nil},
		// &PingResponse{8, 0, "64.233.174.41", 15169, 0, Latency{0, "Time to live exceeded"}, nil},
		// &PingResponse{387, 0, "", 0, 0, Latency{0, "Request timeout"}, nil},
		&PingResponse{"64 bytes from 8.8.8.8: icmp_seq=2 ttl=46 time=34.540 ms", nil},
		&PingResponse{"From 64.233.174.41 icmp_seq=8 Time to live exceeded", nil},
		&PingResponse{"Request timeout for icmp_seq 387", nil},

		&PingSummary{net.ParseIP("8.8.8.8"), 58, 58, 0, 0.0, 34.75, 40.406, 70.66, 5.502, nil},
		&PingSummary{net.ParseIP("192.0.2.1"), 11, 0, 7, 100.0, 0, 0, 0, 0, nil},

		&PingResponse{"64 bytes from 8.8.4.4: icmp_seq=4 ttl=52 time=26.7 ms", nil},
		&PingResponse{"64 bytes from 8.8.4.4: icmp_seq=5 ttl=52 time=26.7 ms", nil},
		&PingSummary{net.ParseIP("8.8.4.4"), 5, 5, 0, 0.0, 26.768, 26.783, 26.793, 0.146, nil},

		&PingSummary{net.ParseIP("4.4.2.2"), 10, 0, 0, 100.0, 0, 0, 0, 0, nil},
	}
	strings := []string{
		"64 bytes from 8.8.8.8: icmp_seq=2 ttl=46 time=34.540 ms",
		"From 64.233.174.41 icmp_seq=8 Time to live exceeded",
		"Request timeout for icmp_seq 387",
		"--- 8.8.8.8 ping statistics ---\n" +
			"58 packets transmitted, 58 packets received, 0.0% packet loss\n" +
			"round-trip min/avg/max/stddev = 34.750/40.406/70.660/5.502 ms\n",
	}

	// counting test results, not input
	i := 0

	for _, group := range testLines {

		prp := NewPingParser()

		go func() {
			for _, line := range group {
				prp.Add(line)
			}
			prp.Close()
		}()

		for {
			pr := prp.Read()
			if pr == nil {
				break
			}
			assert.Nil(t, pr.Error())
			if pr.Error() != nil {
				assert.Fail(t, fmt.Sprintf("err: %s\n", pr.Error()))
				continue
			}
			pretty.Println("Got result:", i, pr)
			pretty.Println("Expects", results[i])
			pretty.Println("Got    ", pr)
			assert.Equal(t, results[i], pr)

			if len(strings) > i && len(strings[i]) > 0 {
				assert.Equal(t, strings[i], pr.String())
			}

			i = i + 1
		}
	}

}
