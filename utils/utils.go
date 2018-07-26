package utils

import (
	"errors"
	"regexp"
	"sync"
	"time"

	"gopkg.in/birkirb/loggers.v1/log"

	"github.com/toolkits/net/httplib"
)

// Pushover implements send messages to pushover
func Pushover(userKeys, apiToken, message string) error {
	const pushoverAPI = "https://api.pushover.net/1/messages.json"
	re := regexp.MustCompile("[,;]")
	userKeySlice := re.Split(userKeys, -1)
	errUserKeyChan := make(chan string, len(userKeySlice))
	errChan := make(chan error, 1)

	go func() {
		errUserKeys := make([]string, 0)
		for errUserKey := range errUserKeyChan {
			errUserKeys = append(errUserKeys, errUserKey)
		}
		if len(errUserKeys) > 0 {
			errChan <- errors.New("pushover request error")
			return
		}
		errChan <- nil
	}()

	var wg sync.WaitGroup
	for _, userKey := range userKeySlice {
		wg.Add(1)
		go func(userKey, apiToken, message string) {
			defer wg.Done()
			req := httplib.Post(pushoverAPI).SetTimeout(3*time.Second, 5*time.Second)
			req.Param("user", userKey)
			req.Param("token", apiToken)
			req.Param("message", message)
			resp, err := req.Response()
			if err != nil {
				log.Errorf("pushover to %s error: %v", userKey, err)
				errUserKeyChan <- userKey
				return
			}
			if resp.StatusCode != 200 {
				log.Errorf("pushover to %s error: %v", userKey, resp.Status)
				errUserKeyChan <- userKey
			}
		}(userKey, apiToken, message)
	}
	wg.Wait()
	close(errUserKeyChan)
	return <-errChan
}
