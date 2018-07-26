package utils

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/toolkits/net/httplib"
)

// Pushover implements send messages to pushover
func Pushover(userkey, apptoken, message, priority string) error {
	const pushoverAPI = "https://api.pushover.net/1/messages.json"
	req := httplib.Post(pushoverAPI).SetTimeout(3*time.Second, 5*time.Second)
	req.Param("user", userkey)
	req.Param("token", apptoken)
	req.Param("message", message)
	resp, err := req.Response()
	if err != nil {
		return fmt.Errorf("pushover send error: %v", err)
	}
	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			return fmt.Errorf("pushover send error: %s", string(body))
		}
		return fmt.Errorf("pushover send error: %s", resp.Status)
	}
	return nil
}
