package handlers

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/netip"

	"go4.org/netipx"
)

// Limiter manages the rate limiting
type Limiter interface {
	Check() bool
	Allowed() bool
	Done()
}

// Handlers have some common functions for the http handlers
type Handlers struct {
	PrivateNets []netip.Prefix
}

// NewHandlers returns Handlers
func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) getIP(prefix string, req *http.Request) (*netip.Addr, error) {

	ipparam := ""

	if len(req.URL.Path) > len(prefix) {
		ipparam = req.URL.Path[len(prefix):]
	}

	if len(ipparam) == 0 {
		return nil, nil
	}

	ip, err := netip.ParseAddr(ipparam)
	if err != nil || !ip.IsValid() {
		ips, err := net.LookupIP(ipparam)
		if err != nil {
			log.Printf("could not lookup %q: %s", ipparam, err)
			return nil, errors.New("could not find")
		}
		if len(ips) > 0 {
			ip, _ = netipx.FromStdIP(ips[0])
		} else {
			log.Printf("unspecified ip '%s'", ipparam)
			return nil, errors.New("unspecified ip")
		}
	}

	for _, p := range h.PrivateNets {
		if p.Contains(ip) {
			log.Printf("private ip '%s'", ip.String())
			return nil, errors.New("private ip")
		}
	}

	return &ip, nil

}
