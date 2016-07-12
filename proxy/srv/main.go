package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kazoup/platform/proxy/srv/proxy"
	"github.com/kazoup/platform/proxy/srv/server"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/selector"
)

var (
	re = regexp.MustCompile("^[a-zA-Z0-9]+$")
	// Default address to bind to
	Address = ":8000"
	// The namespace to serve
	// Example:
	// Namespace + /[Service]/foo/bar
	// Host: Namespace.Service Endpoint: /foo/bar
	Namespace = "go.micro.web"
	// Base path sent to web service.
	// This is stripped from the request path
	// Allows the web service to define absolute paths
	BasePathHeader = "X-Micro-Web-Base-Path"
	statsURL       string
)

type srv struct {
	*mux.Router
}

func (s *srv) proxy() http.Handler {
	sel := selector.NewSelector(
		selector.Registry((*cmd.DefaultOptions().Registry)),
	)
	sel.Init()
	director := func(r *http.Request) {
		kill := func() {
			r.URL.Host = ""
			r.URL.Path = ""
			r.URL.Scheme = ""
			r.Host = ""
			r.RequestURI = ""
		}

		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 2 {
			kill()
			return
		}
		if !re.MatchString(parts[1]) {
			kill()
			return
		}
		next, err := sel.Select(Namespace + "." + parts[1])
		if err != nil {
			kill()
			return
		}

		s, err := next()
		if err != nil {
			kill()
			return
		}
		r.Header.Set(BasePathHeader, "/"+parts[1])
		r.URL.Host = fmt.Sprintf("%s:%d", s.Address, s.Port)
		r.URL.Path = "/" + strings.Join(parts[2:], "/")
		r.URL.Scheme = "http"
		r.Host = r.URL.Host
	}

	return &proxy.Proxy{
		Default:  &httputil.ReverseProxy{Director: director},
		Director: director,
	}
}

func main() {
	cmd.Init()
	var h http.Handler
	r := mux.NewRouter()
	s := &srv{r}
	h = s

	s.PathPrefix("/{service:[a-zA-Z0-9]+}").Handler(s.proxy())

	var opts []server.Option

	srv := server.NewServer(Address)
	srv.Init(opts...)
	srv.Handle("/", h)

	// Initialise Server
	service := micro.NewService(
		micro.Name("go.micro.srv.proxy"),
		micro.Version("latest"),
	)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
	// Init service
	service.Init()
	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

	if err := srv.Stop(); err != nil {
		log.Fatal(err)
	}
}
