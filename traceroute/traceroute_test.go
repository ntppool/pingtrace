package traceroute

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
				assert.Nil(t, tr.Error())
				if tr.Error() != nil {
					fmt.Printf("err: %s\n", tr.Error())
					continue
				}
				// pretty.Println("Got result:", i, tr)
				// pretty.Println("Expects", results[i])
				// pretty.Println("Got    ", tr)

				if len(td.Results) > i {
					if td.Results[i] != nil {
						assert.Equal(t, td.Results[i], tr)
					}
				}
				if len(td.Strings) > i && len(td.Strings[i]) > 0 {
					assert.Equal(t, td.Strings[i], tr.String())
				}

				i = i + 1
			}

		}
	}

}
