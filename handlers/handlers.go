package handlers

import (
	"errors"
	"log"
	"net"
	"net/http"
)

type Limiter interface {
	Check() bool
	Allowed() bool
	Done()
}

type Handlers struct {
	PrivateNets []*net.IPNet
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) getIP(prefix string, req *http.Request) (*net.IP, error) {

	ipparam := ""

	if len(req.URL.Path) > len(prefix) {
		ipparam = req.URL.Path[len(prefix):]
	}

	if len(ipparam) == 0 {
		return nil, nil
	}

	ip := net.ParseIP(ipparam)

	if ip == nil || ip.IsUnspecified() {
		ips, err := net.LookupIP(ipparam)
		if err != nil {
			log.Printf("could not lookup %q: %s", ipparam, err)
			return nil, errors.New("could not find")
		}
		if len(ips) > 0 {
			ip = ips[0]
		} else {
			log.Printf("unspecified ip '%s'", ipparam)
			return nil, errors.New("unspecified ip")
		}
	}

	for _, p := range h.PrivateNets {
		if p.Contains(ip) {
			log.Println("private ip '%s'", ip.String)
			return nil, errors.New("private ip")
		}
	}

	return &ip, nil

}
