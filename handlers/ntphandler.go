package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/beevik/ntp"
)

// GET /ntp/{ip}
func (h *Handlers) NTPHandler(w http.ResponseWriter, req *http.Request) {

	// todo: every so often do a check against the local NTP server,
	// fail if it's not accurate enough?

	ip, err := h.getIP("/ntp/", req)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	if ip == nil {
		w.WriteHeader(404)
		return
	}

	resp, err := ntp.Query(ip.String(), 4)
	if err != nil {
		b, jerr := json.Marshal(err)
		if jerr != nil {
			log.Printf("Could not marshall error json '%s': %s", err, jerr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusGatewayTimeout)
		w.Write(b)
		return
	}

	b, jerr := json.Marshal(resp)
	if jerr != nil {
		log.Printf("Could not marshall json '%+v': %s", resp, jerr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return

}
