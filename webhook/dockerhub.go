package webhook

import (
	"fmt"
	"net/http"
)

// DockerHubHandler handlers webhook push data from docker hub
func DockerHubHandler() http.Handler {
	dockerHubMux := http.NewServeMux()
	dockerHubMux.HandleFunc("/dockerhub", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})
	return dockerHubMux
}
