package main

import (
	"net/http"
	"time"
	"github.com/garyburd/redigo/redis"
	"log"
	"strconv"
	"flag"
)

var (
	httpServerPort int
)

func init() {
	const (
		defaultHttpServerPort 	= 8080
		usage       			= "http server port"
	)

	flag.IntVar(&httpServerPort, "port", defaultHttpServerPort, usage)

	redisPool = &redis.Pool{
		MaxIdle: 5,
		IdleTimeout: 240 * time.Second,
		Dial: func () (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				return nil, err
			}

			return c, err
		},
	}
}

func main() {
	log.Println("api.qpcrbox.com")

	http.HandleFunc("/v1/qpcr/", qpcrHandler)
	http.HandleFunc("/v1/experiment/", experimentHandler)
	http.HandleFunc("/v1/rate-limit", rateLimitHandler)
	http.HandleFunc("/v1/status", statusHandler)
	http.ListenAndServe(":" + strconv.Itoa(httpServerPort), nil)
}
