package fl

import (
	"log"
	"strings"
	"net/url"
	"net/http"

)
type RPCServer struct {
	Address string
}

func reqHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[DEBUG] RPCServer : %s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		// Some clients are adding a // at then end - this throws http server off
		if strings.HasSuffix(r.URL.Path, "//") {
			p := strings.TrimSuffix(r.URL.Path, "/")
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = p
			handler.ServeHTTP(w, r2)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}

func (s *RPCServer) Start()  {

	r := newRPCServer()

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(s.Address, reqHandler(http.DefaultServeMux)))

}
