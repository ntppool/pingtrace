package traceroute

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTracerouteLine(t *testing.T) {

	testLines := []string{
		" 1  ns.unitedwifi.com (172.19.248.1)  3.353 ms  0.959 ms  0.858 ms",
		" 9  209.85.249.5 (209.85.249.5)  445.505 ms",
		"    209.85.249.5 (209.85.249.5)  50.725 ms  2026.411 ms",

		// icmp error codes
		" 8 * * *",
		"15  lax1.ntppool.net (207.171.3.2)  67.841 ms !X  67.827 ms !X  67.811 ms !X",

		// multiple IPs on a line (from linux traceroute)
		" 7  72.14.239.48 (72.14.239.48)  49.941 ms 216.239.46.53 (216.239.46.53)  50.084 ms",
		// " 7 bbisp-gw1-example.com (10.0.74.100) 11.379 ms bbisp-gw2-ae256vla100.example.com (10.0.74.106) 10.554 ms",
	}
	results := []*TraceRouteLine{
		&TraceRouteLine{1, "ns.unitedwifi.com", "172.19.248.1", 0,
			[]Latency{Latency{3.353, ""}, Latency{0.959, ""}, Latency{0.858, ""}}, nil,
		},
		&TraceRouteLine{9, "", "209.85.249.5", 15169, []Latency{Latency{445.505, ""}}, nil},
		&TraceRouteLine{0, "", "209.85.249.5", 15169,
			[]Latency{Latency{50.725, ""}, Latency{2026.411, ""}},
			nil,
		},
		&TraceRouteLine{8, "", "", 0,
			[]Latency{
				Latency{0, "*"},
				Latency{0, "*"},
				Latency{0, "*"},
			},
			nil,
		},
		&TraceRouteLine{15, "lax1.ntppool.net", "207.171.3.2", 7012,
			[]Latency{
				Latency{67.841, "!X"},
				Latency{67.827, "!X"},
				Latency{67.811, "!X"},
			},
			nil,
		},
		&TraceRouteLine{7, "", "72.14.239.48", 15169, []Latency{Latency{49.941, ""}}, nil},
		&TraceRouteLine{7, "", "216.239.46.53", 15169, []Latency{Latency{50.084, ""}}, nil},
		// &TraceRouteLine{7, "bbisp-gw1-ae256vla100.example.com", "10.0.74.100", 714, []Latency{Latency{11.379, ""}}, nil},
		// &TraceRouteLine{7, "bbisp-gw2-ae256vla100.example.com", "10.0.74.106", 714, []Latency{Latency{10.554, ""}}, nil},
	}
	strings := []string{
		" 1 ns.unitedwifi.com (172.19.248.1)  3.353  0.959  0.858",
		" 9 (209.85.249.5) AS15169  445.505",
		" 0 (209.85.249.5) AS15169  50.725  2026.411",
		" 8  *  *  *",
		"15 lax1.ntppool.net (207.171.3.2) AS7012  67.841 !X  67.827 !X  67.811 !X",
	}

	// counting test results, not input
	i := 0

	for _, line := range testLines {

		trp := NewTracerouteParser()

		go func() {
			trp.Add(line)
			trp.Close()
		}()

		for {
			tr := trp.Read()
			if tr == nil {
				break
			}
			assert.Nil(t, tr.Error())
			if tr.Error() != nil {
				fmt.Printf("err: %s\n", tr.Error())
				continue
			}
			// pretty.Println("Got result:", i, tr)
			// pretty.Println("Expects", results[i])
			// pretty.Println("Got    ", tr)

			if len(results) > i {
				assert.Equal(t, results[i], tr)
			}
			if len(strings) > i && len(strings[i]) > 0 {
				assert.Equal(t, strings[i], tr.String())
			}

			i = i + 1
		}

	}

}
