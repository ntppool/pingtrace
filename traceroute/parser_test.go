package traceroute

import (
	"testing"
)

func TestParseTracerouteLine(t *testing.T) {

	for n, td := range Tests {
		if n == 0 {
			continue
		}

		// counting test results, not input
		i := 0

		for _, line := range td.Lines {

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

				if tr.Error() != nil {
					t.Fatalf("err: %s\n", tr.Error())
				}
				// pretty.Println("Got result:", i, tr)
				// pretty.Println("Expects", results[i])
				// pretty.Println("Got    ", tr)

				if len(td.Results) > i {
					if td.Results[i] != nil {
						if td.Results[i].String() != tr.String() {
							t.Logf("got %q, expected %q", tr.String(), td.Results[i].String())
							t.Fail()
						}
					}
				}
				if len(td.Strings) > i && len(td.Strings[i]) > 0 {
					if td.Strings[i] != tr.String() {
						t.Logf("got %q, expected %q", tr.String(), td.Strings[i])
						t.Fail()

					}
				}

				i = i + 1
			}

		}
	}

}
