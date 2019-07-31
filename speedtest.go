package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cavaliercoder/grab"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/schollz/progressbar"
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
	fmt.Println("\n", speed, metric)
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
	barCheck := 0
	progress := 0

	var bar *progressbar.ProgressBar
	if !hideBar {
		bar = progressbar.New64(resp.Size())
	}

	for {
		progress = int(resp.BytesComplete())

		if barCheck != progress {
			barCheck = progress
			if !hideBar {
				bar.Add(barCheck)
			}
		}

		if resp.IsComplete() {
			break
		}
	}

	if !hideBar {
		bar.Finish()
	}

	err = resp.Err()
	if err != nil {
		return 0
	}

	mbps := int(resp.BytesPerSecond() * 0.000008)
	os.Remove(downloadDest)

	return mbps
}
