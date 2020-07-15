package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"runtime"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

type ConfigTreeDatabase struct {
	Host string `yaml:"host"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Name string `yaml:"name"`
}

type ConfigTreeEmail struct {
	Server string `yaml:"server"`
	Port   int    `yaml:"port"`
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	Secure bool   `yaml:"secure"`
}

type ConfigTreeRoot struct {
	Key      string `yaml:"apikey"`
	Database ConfigTreeDatabase
	Email    ConfigTreeEmail
}

var configPath string
var local string
var tcp string
var unix string
var debug bool

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&configPath, "config", "config.yaml",
		"Path to the config file. ex: /path/to/config.yaml")
	flag.StringVar(&local, "local", "",
		"Start the webserver, ex: 0.0.0.0:8000")
	flag.StringVar(&tcp, "tcp", "",
		"Start a FastCGI TCP network socket, ex: 0.0.0.0:8000")
	flag.StringVar(&unix, "unix", "",
		"Start a FastCGI UNIX socket, ex: /tmp/program.sock")
	flag.BoolVar(&debug, "debug", false,
		"Turn on debugging output.")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()

	var err error

	r := mux.NewRouter()

	r.HandleFunc("/", hello)

	// Loading config file.
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("ERROR: Failure to load config file. %s\n", err)
	}

	config := ConfigTreeRoot{}
	if err = yaml.Unmarshal(configData, &config); err != nil {
		log.Fatalf("ERROR: Failed to convert config data to struct. %s", err)
	}

	if debug {
		log.Printf("DEBUG: Config settings: %+v\n", config)
	}

	if local != "" { // Run the web server
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

func hello(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers.Set("Content-Type", "text/html")
	io.WriteString(w, "<html><head></head><body><p>Hello from Go!</p><table>")
	io.WriteString(w, fmt.Sprintf("<tr><td>Method</td><td>%s</td></tr>", r.Method))
	io.WriteString(w, fmt.Sprintf("<tr><td>URL</td><td>%s</td></tr>", r.URL))
	io.WriteString(w, fmt.Sprintf("<tr><td>URL.Path</td><td>%s</td></tr>", r.URL.Path))
	io.WriteString(w, fmt.Sprintf("<tr><td>Proto</td><td>%s</td></tr>", r.Proto))
	io.WriteString(w, fmt.Sprintf("<tr><td>Host</td><td>%s</td></tr>", r.Host))
	io.WriteString(w, fmt.Sprintf("<tr><td>RemoteAddr</td><td>%s</td></tr>", r.RemoteAddr))
	io.WriteString(w, fmt.Sprintf("<tr><td>RequestURI</td><td>%s</td></tr>", r.RequestURI))
	io.WriteString(w, fmt.Sprintf("<tr><td>Header</td><td>%s</td></tr>", r.Header))
	io.WriteString(w, fmt.Sprintf("<tr><td>Body</td><td>%s</td></tr>", r.Body))
	io.WriteString(w, "</table></body></html>")
}

/*
 * References:
 *	- http://www.dav-muz.net/blog/2013/09/how-to-use-go-and-fastcgi/
 *	- https://discussion.dreamhost.com/t/how-to-run-go-language-programs-on-dreamhost-servers-using-fastcgi/64844
 *	- https://mwholt.blogspot.com/2013/05/writing-go-golang-web-app-with-nginx.html
 *  - https://github.com/bsingr/golang-apache-fastcgi/blob/master/examples/vanilla/hello_world.go
 *
 */
