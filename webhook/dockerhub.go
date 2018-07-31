package webhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"gopkg.in/birkirb/loggers.v1/log"

	"github.com/WuHan0608/webhook-pushover/utils"
)

type pushData struct {
	Images   []string `json:"images"`
	PushedAt int64    `json:"pushed_at"`
	Pusher   string   `json:"pusher"`
	Tag      string   `json:"tag"`
}

type repository struct {
	CommentCount    int64  `json:"comment_count"`
	DataCreated     int64  `json:"date_created"`
	Description     string `json:"description"`
	Dockefile       string `json:"dockerfile"`
	FullDescription string `json:"full_description"`
	IsOfficial      bool   `json:"is_official"`
	IsPrivate       bool   `json:"is_private"`
	IsTrusted       bool   `json:"is_trusted"`
	Name            string `json:"name"`
	NameSpace       string `json:"namespace"`
	Owner           string `json:"owner"`
	RepoName        string `json:"repo_name"`
	RepoURL         string `json:"repo_url"`
	StarCount       int64  `json:"star_count"`
	Status          string `json:"status"`
}

type payload struct {
	CallbackURL string      `json:"callback_url"`
	PushData    *pushData   `json:"push_data"`
	Repository  *repository `json:"repository"`
}

// DockerHubHandler handlers webhook push data from docker hub
func DockerHubHandler() http.Handler {
	dockerHubMux := http.NewServeMux()
	dockerHubMux.HandleFunc("/dockerhub", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Fprintln(w, "OK")
		} else if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Errorf("read request error: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			data := payload{}
			if err := json.Unmarshal(body, &data); err != nil {
				log.Errorf("decode payload error: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var (
				trackTag = os.Getenv("TRACK_REPO_TAG")
				userKeys = os.Getenv("PUSHOVER_USER_KEYS")
				apiToken = os.Getenv("PUSHOVER_API_TOKEN")
			)
			if len(trackTag) > 0 {
				if data.PushData.Tag != trackTag {
					log.Infof("track repo tag \"%s\" only, ignore tag \"%s\"", trackTag, data.PushData.Tag)
					return
				}
			}
			if len(userKeys) == 0 || len(apiToken) == 0 {
				log.Error("pushover request error: no user keys or api token is set")
				http.Error(w, "pushover request error", http.StatusInternalServerError)
				return
			}
			pushedAtTime := time.Unix(data.PushData.PushedAt, 0)
			message := fmt.Sprintf("%s:%s is pushed by %s at %s", data.Repository.RepoName, data.PushData.Tag, data.PushData.Pusher, pushedAtTime.String())
			if err := utils.Pushover(userKeys, apiToken, message); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Allowed methods: GET, POST", http.StatusMethodNotAllowed)
		}
	})
	return dockerHubMux
}
