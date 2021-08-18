package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	argument := &Argument{
		Host: flag.String("h", "https://baidu.com", "主机名"),
		Port: flag.Int("p", 8082, "启动端口"),
	}
	remote, err := url.Parse(*argument.Host)
	if err != nil {
		panic(err)
	}
	proxy := NewProxy(remote)
	// use http.Handle instead of http.HandleFunc when your struct implements http.Handler interface
	http.Handle("/", &ProxyHandler{proxy})
	fmt.Printf("源站:%s  服务启动在:%d🚀", *argument.Host, *argument.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", *argument.Port), nil)

	if err != nil {
		panic(err)
	}

}

type ProxyHandler struct {
	p *httputil.ReverseProxy
}

type Argument struct {
	Host *string
	Port *int
}

func (ph *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	w.Header().Set("X-Ben", "Rad")
	ph.p.ServeHTTP(w, r)
}
func NewProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.Host = target.Host // -- 加入这句 --
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
