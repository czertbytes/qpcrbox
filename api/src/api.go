package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"io/ioutil"
	"log"
	"encoding/json"
)

const (
	rateLimit = 50
)

type ComputationResponse struct {
	ExpiresAt, ExperimentId string
}

type ConsumerRateLimit struct {
	Exceeded bool
	Limit, Current int
	RetryAfter time.Time
}

func qpcrHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[handler|qpcr] request '%s'\n", r.URL)
	log.Printf("[handler|qpcr] headers: %+v\n", r.Header)

	if r.Method == "OPTIONS" {
		log.Println("[handler|qpcr] options")
		w.WriteHeader(http.StatusOK)
		return
	}

	var consumerRateLimit ConsumerRateLimit
	var err error

	if consumerRateLimit, err = checkIPAddressRateLimit(r); err != nil {
		log.Printf("[handler|qpcr] checking consumer rate limit failed with error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if consumerRateLimit.Exceeded {
		w.Header().Add("Retry-After", fmt.Sprintf("%s", consumerRateLimit.RetryAfter))
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.WriteHeader(429)
		return
	}

	if r.Method != "POST" {
		log.Printf("[handler|qpcr] method '%s' is not POST!\n", r.Method)
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	urlPath := strings.Split(r.URL.Path[1:], "/")
	if len(urlPath) != 3 {
		log.Printf("[handler|qpcr] path '%s' is not valid!\n", r.URL.Path[1:])
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	bodyContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("[handler|qpcr] body parameter is not valid!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var expComputer ExperimentComputer
	switch urlPath[2] {
	case "ab7300":
		mock := r.FormValue("mock")
		if len(mock) == 0 {
			log.Println("[handler|qpcr|ab7300] missing mock query parameter!")
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		log.Println("[handler|qpcr] experiment computer set to ab7300")
		expComputer = &AB7300{Content: string(bodyContent), Mock: mock}
	default:
		log.Printf("[handler|qpcr] experiment computer type '%s' is not valid!", urlPath[1])
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	expId := doExperimentComputation(w, expComputer)
	if len(expId) == 0 {
		return
	}

	w.Header().Add("Location", "http://api.fastqpcr.com/experiment/" + expId)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	computationResponse := ComputationResponse{}
	computationResponse.ExpiresAt = fmt.Sprintf("%s", time.Now().Add(time.Duration(30) * time.Minute))
	computationResponse.ExperimentId = expId

	response, err := json.Marshal(computationResponse)
	if err != nil {
		log.Printf("[handler|qpcr] marshalling computationResponse failed with error '%s'\n", err)
	}
	w.Write(response)
}

func doExperimentComputation(w http.ResponseWriter, expComputer ExperimentComputer) string {
	e, err := expComputer.Compute()
	if err != nil {
		log.Printf("[handler|qpcr] experiment computation failed with error '%s'!\n", err)
		http.Error(w, "", http.StatusBadRequest)
		return ""
	}

	expId, err := SaveExperiment(e)
	if err != nil {
		log.Printf("[handler|qpcr] experiment computation persistence failed with error '%s'!\n", err)
		http.Error(w, "", http.StatusInternalServerError)
		return ""
	}

	log.Printf("[handler|qpcr] experiment computed, key %s \n", expId)

	return expId
}

func experimentHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[handler|experiment] request  %s\n", r.URL)
	log.Printf("[handler|experiment] headers: %+v\n", r.Header)

	if r.Method == "OPTIONS" {
		log.Println("[handler|experiment] options")
		w.WriteHeader(http.StatusOK)
		return
	}

	var consumerRateLimit ConsumerRateLimit
	var err error

	if consumerRateLimit, err = checkIPAddressRateLimit(r); err != nil {
		log.Printf("[handler|qpcr] checking consumer rate limit failed with error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if consumerRateLimit.Exceeded {
		w.Header().Add("Retry-After", fmt.Sprintf("%s", consumerRateLimit.RetryAfter))
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.WriteHeader(429)
		return
	}

	if r.Method != "GET" {
		log.Printf("[handler|experiment] method '%s' is not GET!\n", r.Method)
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	var ex Exporter
	accept := r.Header.Get("Accept")
	if accept == "" {
		log.Println("[handler|experiment] accept type is not set, using application/json")
		accept = "application/json"
	}

	switch strings.ToLower(accept) {
	case "application/json":
		log.Println("[handler|experiment] exporter set to json")
		ex = &JSONExport{}
	case "application/xml":
		log.Println("[handler|experiment] exporter set to xml")
		ex = &XMLExport{}
	case "text/csv":
		log.Println("[handler|experiment] exporter set to csv")
		ex = &CSVExport{}
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		log.Println("[handler|experiment] exporter set to xlsx")
		ex = &XLSXExport{}
	case "application/vnd.oasis.opendocument.spreadsheet":
		log.Println("[handler|experiment] experiment computer set to ods")
		ex = &ODSExport{}
	default:
		log.Printf("[handler|experiment] accept type '%s' is not valid!\n", accept)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	urlPath := strings.Split(r.URL.Path[1:], "/")
	if len(urlPath) != 3 {
		log.Printf("[handler|experiment] experiment id '%s' is not valid!\n", r.URL.Path[1:])
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	content := readExperimentResults(w, urlPath[2], ex)
	if len(content) == 0 {
		return
	}

	w.Header().Add("Content-Type", ex.ContentType())
	w.Write(content)
}

func readExperimentResults(w http.ResponseWriter, expId string, ex Exporter) []byte {
	expBytes, err := GetExperiment(expId)
	if err != nil {
		log.Printf("[handler|experiment] experiment id '%s' not found!\n", expId)
		http.Error(w, "", http.StatusNotFound)
		return []byte{}
	}

	var e Experiment
	err = json.Unmarshal(expBytes, &e)
	if err != nil {
		log.Printf("[handler|experiment] parsing experiment id '%s' failed!\n", expId)
		http.Error(w, "", http.StatusInternalServerError)
		return []byte{}
	}

	content, err := ex.Export(&e)
	if err != nil {
		log.Printf("[handler|experiment] exporting experiment id '%s' failed!\n", expId)
		http.Error(w, "", http.StatusInternalServerError)
		return []byte{}
	}

	log.Println("[handler|experiment] experiment content generated")

	return content
}

func rateLimitHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[handler|ratelimit] request  %s\n", r.URL)
	log.Printf("[handler|ratelimit] headers: %+v\n", r.Header)

	if r.Method == "OPTIONS" {
		log.Println("[handler|ratelimit] options")
		w.WriteHeader(http.StatusOK)
		return
	}

	var consumerRateLimit ConsumerRateLimit
	var err error

	if consumerRateLimit, err = checkIPAddressRateLimit(r); err != nil {
		log.Printf("[handler|ratelimit] checking consumer rate limit failed with error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if consumerRateLimit.Exceeded {
		w.Header().Add("Retry-After", fmt.Sprintf("%s", consumerRateLimit.RetryAfter))
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.WriteHeader(429)
		return
	}

	if r.Method != "GET" {
		log.Printf("[handler|ratelimit] method '%s' is not GET!\n", r.Method)
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	content, err := json.Marshal(consumerRateLimit)
	if err != nil {
		log.Printf("[handler|ratelimit] marshalling consumer rate limit failed with error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(content)
}

func checkIPAddressRateLimit(r *http.Request) (ConsumerRateLimit, error) {
	var ipAddress, consumerToken string
	var counter int
	var err error

	ipAddress = r.Header.Get("X-Real-Ip")
	if consumerToken = r.FormValue("consumer-token"); len(consumerToken) == 0 {
		consumerToken = r.Header.Get("Consumer-Token")
	}

	timeNow := time.Now()
	//	retryTime is time first minute of next hour (timeNow 16:01 -> retryTime 17:00)
	retryAfter := timeNow.Add(tokenExpiresTime * time.Second).Add(time.Duration((-1) * timeNow.Minute()) * time.Minute)
	if counter, err = GetRateLimitCounter(ipAddress, timeNow); err != nil {
		log.Printf("[handler|ratelimit] getting counter for ip address '%s' failed with error: %s\n", ipAddress, err)
		return ConsumerRateLimit{}, err
	}

	if counter >= rateLimit {
		if len(consumerToken) == 0 {
			log.Printf("[handler|ratelimit] ip address '%s' exceeded the limit!\n", ipAddress)
			return ConsumerRateLimit{Exceeded: true, Current: counter, Limit: rateLimit, RetryAfter: retryAfter}, nil
		}

		var tokenExists bool
		if tokenExists, err = GetConsumerToken(consumerToken); err != nil {
			log.Printf("[handler|ratelimit] getting consumer token '%s' for ip address '%s' failed with error: %s\n", consumerToken, ipAddress, err)
			return ConsumerRateLimit{}, err
		}

		if !tokenExists {
			log.Printf("[handler|ratelimit] consumer token '%s' for ip address '%s' is not valid!\n", consumerToken, ipAddress, err)
			return ConsumerRateLimit{Exceeded: true, Current: counter, Limit: rateLimit, RetryAfter: retryAfter}, nil
		}
	}

	return ConsumerRateLimit{Exceeded: false, Current: counter, Limit: rateLimit, RetryAfter: retryAfter}, nil
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	//	TODO: add redis PING-PONG, disk space check, ...

	w.WriteHeader(http.StatusOK)
}



