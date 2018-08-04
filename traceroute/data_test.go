package traceroute

type TestData struct {
	Lines   []string
	Results []*TraceRouteLine
	Strings []string
}

var Tests = []TestData{
	TestData{
		Lines: []string{
			" 1  ns.unitedwifi.com (172.19.248.1)  3.353 ms  0.959 ms  0.858 ms",
			" 9  209.85.249.5 (209.85.249.5)  445.505 ms",
			"    209.85.249.5 (209.85.249.5)  50.725 ms  2026.411 ms",

			// icmp error codes
			" 8 * * *",
			"15  lax1.ntppool.net (207.171.3.2)  67.841 ms !X  67.827 ms !X  67.811 ms !X",

			// multiple IPs on a line (from linux traceroute)
			" 7  72.14.239.48 (72.14.239.48)  49.941 ms 216.239.46.53 (216.239.46.53)  50.084 ms",
			// " 7 bbisp-gw1-example.com (10.0.74.100) 11.379 ms bbisp-gw2-ae256vla100.example.com (10.0.74.106) 10.554 ms",
		},
		Results: []*TraceRouteLine{
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
		},
		Strings: []string{
			" 1 ns.unitedwifi.com (172.19.248.1)  3.353  0.959  0.858",
			" 9 (209.85.249.5) AS15169  445.505",
			" 0 (209.85.249.5) AS15169  50.725  2026.411",
			" 8  *  *  *",
			"15 lax1.ntppool.net (207.171.3.2) AS7012  67.841 !X  67.827 !X  67.811 !X",
		},
	},
	TestData{
		Lines: []string{
			"traceroute to 207.171.3.4 (207.171.3.4), 64 hops max, 52 byte packets",
			" 1  192.168.0.1  5.742 ms  1.796 ms  1.945 ms",
			" 2  10.198.6.113  38.144 ms  28.064 ms  28.551 ms",
			" 3  10.170.214.10  45.213 ms  26.221 ms  29.326 ms",
			" 4  10.164.162.196  34.272 ms  31.933 ms  36.128 ms",
			" 5  10.164.165.25  37.344 ms  32.390 ms  38.076 ms",
			" 6  208.185.86.21  34.511 ms  38.841 ms  31.178 ms",
			" 7  64.125.14.30  39.902 ms  40.939 ms  39.429 ms",
			" 8  129.250.2.229  49.507 ms",
			"    129.250.3.59  39.635 ms",
			"    129.250.2.229  42.812 ms",
			" 9  129.250.3.26  29.561 ms  37.871 ms  40.026 ms",
			"10  129.250.4.151  51.863 ms  38.242 ms  51.610 ms",
			"11  129.250.4.107  43.122 ms  37.377 ms  44.303 ms",
			"12  198.172.90.74  39.952 ms  40.684 ms  48.832 ms",
			"13  207.171.30.62  42.627 ms  145.465 ms  35.547 ms",
			"14  207.171.3.4  41.351 ms !Z  38.477 ms !Z  41.645 ms !Z",
		},
		Results: []*TraceRouteLine{
			nil, nil, nil, nil,
			nil, nil, nil, nil,
			nil, nil, nil, nil,
			nil, nil, nil,
			&TraceRouteLine{14, "207-171-3-4.ntppool.net", "207.171.3.4", 7012,
				[]Latency{
					Latency{41.351, "!Z"},
					Latency{38.477, "!Z"},
					Latency{41.645, "!Z"},
				},
				nil,
			},
		},
	},
}
