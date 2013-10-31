package main

import (
	"fmt"
	"log"
	"io"
	"time"
	"encoding/json"
	"crypto/sha256"
	"github.com/garyburd/redigo/redis"
)

var (
	redisPool *redis.Pool
)

const (
	redisKeyPrefix = "qpcrbox"
	expirimentExpiresTime = 7200 // in seconds
	tokenExpiresTime = 3600 // in seconds
)

func SaveExperiment(e *Experiment) (string, error) {
	expJsonBytes, err := json.Marshal(e)
	if err != nil {
		log.Fatalf("Marshalling experiment to JSON failed! Error: %s\n", err)
	}

	expJson := string(expJsonBytes)
	expId := getExpId(expJson)
	err = persist(expId, expJson)
	if err != nil {
		return "", err
	}

	return expId, nil
}

func GetExperiment(expId string) ([]byte, error) {
	redisConn := redisPool.Get()
	defer redisConn.Close()

	var expBytes []byte
	var err error
	key := fmt.Sprintf("%s:expid:%s", redisKeyPrefix, expId)
	if expBytes, err = redis.Bytes(redisConn.Do("GET", key)); err != nil {
		return []byte{}, err
	}

	return expBytes, nil
}

func GetRateLimitCounter(ipAddress string, timeNow time.Time) (int, error) {
	redisConn := redisPool.Get()
	defer redisConn.Close()

	var counter int
	var err error

	keyCount := fmt.Sprintf("%s:ratelimit:%02d:%s", redisKeyPrefix, timeNow.Hour(), ipAddress)
	if counter, err = redis.Int(redisConn.Do("INCR", keyCount)); err != nil {
		return -1, err
	}
	if _, err = redisConn.Do("EXPIRE", keyCount, tokenExpiresTime - (timeNow.Minute() * 60)); err != nil {
		return -1, err
	}

	log.Printf("[redis|ratelimit] keyCount '%s' has value '%d'\n", keyCount, counter)

	return counter, nil
}

func GetConsumerToken(token string) (bool, error) {
	redisConn := redisPool.Get()
	defer redisConn.Close()

	key := fmt.Sprintf("%s:token:%s", redisKeyPrefix, token)
	res, err := redis.Int(redisConn.Do("EXIST", key))
	if err != nil {
		return false, err
	}

	return res != 0, nil
}

func getExpId(s string) string {
	h := sha256.New()
	io.WriteString(h, s)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func persist(expId, value string) error {
	redisConn := redisPool.Get()
	defer redisConn.Close()

	key := fmt.Sprintf("%s:expid:%s", redisKeyPrefix, expId)
	if _, err := redisConn.Do("SET", key, value); err != nil {
		return err
	}
	if _, err := redisConn.Do("EXPIRE", key, expirimentExpiresTime); err != nil {
		return err
	}

	return nil
}
