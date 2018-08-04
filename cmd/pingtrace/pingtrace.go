package main

import (
	"flag"
	"fmt"

	"github.com/abh/pingtrace/handlers"

	"github.com/rs/cors"

	"log"
	"net"
	"net/http"

	"os"
)

var (
	listen = flag.String("listen", "127.0.0.1:8060", "listen address")

	privateNets []*net.IPNet
)

// A version string that can be set with
//
//     -ldflags "-X main.Version VERSION"
//
// at compile-time.
var Version string

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: pingtrace [-listen=<listen>]")
		flag.PrintDefaults()
	}
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)

	pn := []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"}
	for _, p := range pn {
		_, ipnet, err := net.ParseCIDR(p)
		if err != nil {
			panic(err)
		}
		privateNets = append(privateNets, ipnet)
	}
}

func setupMux(hdl *handlers.Handlers) http.Handler {
	mux := http.NewServeMux()

	rateh := handlers.NewRateHandler(mux, NewLimiter(10))

	mux.Handle("/ntp/", rateh.Wrap(hdl.NTPHandler))
	mux.Handle("/traceroute/", rateh.Wrap(hdl.TracerouteHandler))
	mux.HandleFunc("/traceroute/checkqueue", rateh.CheckQueue)

	return cors.New(cors.Options{AllowedOrigins: []string{"*.pool.ntp.org", "*.ntppool.org"}}).Handler(mux)

}

func main() {
	flag.Parse()

	hdl := handlers.NewHandlers()
	hdl.PrivateNets = privateNets

	h := setupMux(hdl)

	log.Printf("Listening to '%s'", *listen)
	err := http.ListenAndServe(*listen, h)
	if err != nil {
		log.Printf("listen: %s", err)
	}
}
