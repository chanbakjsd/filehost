package main

import (
	"net"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

//requestLimit and sizeLimit are the limiters. limitMutex makes sure access to them are synchronous.
var (
	requestLimit = make(map[string]*rate.Limiter)
	sizeLimit    = make(map[string]*rate.Limiter)
	limitMutex   sync.Mutex
)

//hasHitRequestLimit returns if the requested remote address has reached the request limit.
func hasHitRequestLimit(remoteAddr string) bool {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return true
	}

	limitMutex.Lock()
	defer limitMutex.Unlock()

	// Create if limiters don't exist yet.
	if requestLimit[ip] == nil {
		requestLimit[ip] = rate.NewLimiter(rate.Limit(1.0/12), 5)
		sizeLimit[ip] = rate.NewLimiter(rate.Limit(1024*1024), 100*1024*1024)
	}

	return !requestLimit[ip].Allow()
}

//hasHitSizeLimit returns if the requested remote address has reached the size limit.
func hasHitSizeLimit(remoteAddr string, size int64) bool {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return true
	}

	limitMutex.Lock()
	defer limitMutex.Unlock()

	return !sizeLimit[ip].AllowN(time.Now(), int(size))
}