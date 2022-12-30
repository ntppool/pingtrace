package handlers

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os/exec"
	"sync"

	"go.ntppool.org/pingtrace/traceroute"
)

// todo: wrap handler that gets the IP, checks for private nets and gets a queue slot

// GET /traceroute/{ip}
func (h *Handlers) TracerouteHandler(w http.ResponseWriter, req *http.Request) {

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

	cmd := exec.Command("traceroute", "-q", "2", "-w", "3", "-n", ip.String())
	// cmd := exec.Command("./slowly.sh", "5")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Could not get stdoutpipe: %s", err)
		w.WriteHeader(500)
		return
	}

	err = cmd.Start()
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

	r := bufio.NewReader(stdout)

	trp := traceroute.NewTracerouteParser()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			tr := trp.Read()
			if tr == nil {
				return
			}

			if fmtJSON {
				_, err = w.Write(tr.JSON())
			} else {
				_, err = w.Write(tr.Bytes())
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
	}()

	reading := true

	for reading {
		line, err := r.ReadString('\n')
		if err != nil {
			trp.Close()
			if err == io.EOF {
				// got to the end
				log.Println("eof:", line)
			} else {
				log.Println("Error reading from traceroute pipe: ", err)
			}
			break
		}

		trp.Add(line)
		if err != nil {
			log.Printf("Could not parse '%s': %s", line, err)
			continue
		}
	}

	cmdRV := cmd.Wait()
	if cmdRV != nil {
		log.Printf("Error finishing command: %s", err)
	}

	// make sure we read everything from the parser
	wg.Wait()

	w.(http.Flusher).Flush()
}

//
