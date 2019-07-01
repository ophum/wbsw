package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-redis/redis"
)

var rconn = connectRedis("10.55.37.45:6379", "", 0)

func connectRedis(addr, pass string, db int) (r *redis.Client) {
	r = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})
	return r
}

func findDomain(host string) (nextHost string, err error) {
	val, err := rconn.Get(host).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func main() {
	errpage := "http://localhost:9000"

	director := func(req *http.Request) {
		from := req.Host
		req.URL.Scheme = "http"
		server, err := findDomain(req.Host)
		if err != nil {
			log.Println("Not found host...", from)
			server = errpage
		}
		origin, _ := url.Parse(server)
		req.URL.Host = origin.Host
		log.Println(from, "->", origin.Host)

	}

	proxy := &httputil.ReverseProxy{Director: director}

	log.Println("Server running...")
	http.ListenAndServe(":8080", proxy)
}
