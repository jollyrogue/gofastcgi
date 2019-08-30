package main

import (
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"runtime"

	"github.com/gorilla/mux"
)

var local string
var tcp string
var unix string

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&local, "local", "",
		"serve as webserver, example: 0.0.0.0:8000")
	flag.StringVar(&tcp, "tcp", "",
		"serve as FCGI via TCP, example: 0.0.0.0:8000")
	flag.StringVar(&unix, "unix", "",
		"serve as FCGI via UNIX socket, example: /tmp/myprogram.sock")
}

func hello(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers.Set("Content-Type", "text/html")
	io.WriteString(w, "<html><head></head><body><p>Hello</p></body></html>")
}

func main() {
	flag.Parse()
	r := mux.NewRouter()

	r.HandleFunc("/", hello)

	var err error

	if local != "" { // Run as a local web server
		log.Printf("Running builtin webserver at %s...\n", local)
		err = http.ListenAndServe(local, r)
	} else if tcp != "" { // Run as FCGI via TCP
		log.Printf("Running as FastCGI TCP at %s...", tcp)
		listener, err := net.Listen("tcp", tcp)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		err = fcgi.Serve(listener, r)
	} else if unix != "" { // Run as FCGI via UNIX socket
		log.Printf("Running as FastCGI Unix socket at %s...\n", unix)
		listener, err := net.Listen("unix", unix)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		err = fcgi.Serve(listener, r)
	} else { // Run as FCGI via standard I/O
		log.Printf("Running as FastCGI standard IO...\n")
		err = fcgi.Serve(nil, r)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Shutting down.")
}

/*
 * References:
 *	- http://www.dav-muz.net/blog/2013/09/how-to-use-go-and-fastcgi/
 *	- https://discussion.dreamhost.com/t/how-to-run-go-language-programs-on-dreamhost-servers-using-fastcgi/64844
 *	- https://mwholt.blogspot.com/2013/05/writing-go-golang-web-app-with-nginx.html
 *
 */
