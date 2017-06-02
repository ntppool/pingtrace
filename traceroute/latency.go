package traceroute

import "fmt"

type Latency struct {
	Ms    float64
	Error string `json:",omitempty"`
}

func (l Latency) String() string {
	if len(l.Error) == 0 {
		// no error
		return fmt.Sprintf("%.3f", l.Ms)
	}
	if l.Ms > 0 {
		// error + timestamp
		return fmt.Sprintf("%.3f %s", l.Ms, l.Error)
	}
	// just an error
	return fmt.Sprintf("%s", l.Error)
}
