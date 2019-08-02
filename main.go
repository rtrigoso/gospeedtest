package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cavaliercoder/grab"
	uuid "github.com/nu7hatch/gouuid"
)

var err error

var downloadPath = "/tmp"

type jsonOutput struct {
	Speed  int    `json:"speed"`
	Metric string `json:"metric"`
}

func handleError() {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unique id creation failed: %v\n", err)
		os.Exit(1)
	}
}

func outputText(speed int, metric string) {
	fmt.Println(speed, metric)
}

func outputJSON(speed int, metric string) {
	defer handleError()

	resp := &jsonOutput{Speed: speed, Metric: metric}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		return
	}

	println(string(respJSON))
}

func downloadFile(url string, hideBar bool) int {
	defer handleError()

	id, err := uuid.NewV4()
	if err != nil {
		return 0
	}

	fileExt := filepath.Ext(testFile)
	downloadDest := downloadPath + "/" + id.String() + fileExt

	client := grab.NewClient()
	req, err := grab.NewRequest(downloadDest, testFile)
	if err != nil {
		return 0
	}

	resp := client.Do(req)
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

	done := false
	for {
		select {
		case <-resp.Done:
			done = true
		}

		if done {
			break
		}
	}

	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}

	mbps := int(resp.BytesPerSecond() * 0.000008)
	os.Remove(downloadDest)

	return mbps
}
