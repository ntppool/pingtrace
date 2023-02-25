package handlers

import (
	"log"
	"net/http"
	"net/netip"

	"go.ntppool.org/pingtrace/traceroute"
)

// todo: wrap handler that gets the IP, checks for private nets and gets a queue slot

// GET /traceroute/{ip}
func (h *Handlers) TracerouteHandler(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var fmtJSON bool

	if s := req.URL.Query().Get("json"); len(s) > 0 {
		if s != "0" {
			fmtJSON = true
			w.Header().Set("Content-Type", "application/json")
		}
	}

	ip, err := h.getIP("/traceroute/", req)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	if ip == nil {
		w.WriteHeader(404)
		return
	}

	tr, err := traceroute.New(*ip, netip.Addr{})
	if err != nil {
		log.Printf("traceroute: %s", err)
		w.WriteHeader(400)
		return
	}

	err = tr.Start(ctx)
	if err != nil {
		log.Printf("Could not start traceroute command: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")

	if !fmtJSON {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}

	w.WriteHeader(200)

	if !fmtJSON {
		w.Write([]byte("Traceroute to " + ip.String() + "\n"))
	}

	for {
		trl, err := tr.Read()
		if err != nil {
			log.Printf("read error: %s", err)
			w.WriteHeader(500)
			break
		}
		if trl == nil {
			break
		}

		if fmtJSON {
			_, err = w.Write(trl.JSON())
		} else {
			_, err = w.Write(trl.Bytes())
		}
		if err != nil {
			log.Printf("Error writing: %s", err)
		}
		_, err = w.Write([]byte("\n"))
		if err != nil {
			log.Printf("Error writing: %s", err)
		}
		w.(http.Flusher).Flush()
	}

	// final flush
	w.(http.Flusher).Flush()
}

//
