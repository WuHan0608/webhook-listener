package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

const address = "127.0.0.1:8080"

func TestMain(m *testing.M) {
	os.Setenv("LISTEN_ADDRESS", address)

	go func() {
		main()
	}()

	time.Sleep(time.Second)
	os.Exit(m.Run())
}

func Test_DockerHubHandler(t *testing.T) {
	url := fmt.Sprintf("http://%s/webhook/dockerhub", address)
	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("http status code: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	if strings.TrimSpace(string(body)) != "OK" {
		t.Error("docker hub webhook handler error")
	}
}
