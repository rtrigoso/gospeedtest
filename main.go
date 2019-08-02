package main

import (
	"os"

	"github.com/integrii/flaggy"
)

var testFile = "http://ipv4.download.thinkbroadband.com/50MB.zip"

var flags struct {
	json bool
}

func main() {
	flaggy.SetName("gospeedtest")
	flaggy.SetDescription("a simple internet speed-test cli written in go")
	flaggy.SetVersion("0.0.1")
	flaggy.Bool(&flags.json, "j", "json", "output as json")
	flaggy.Parse()

	mbps := downloadFile(testFile, flags.json)

	switch {
	case flags.json:
		outputJSON(mbps, "Mbps")
	default:
		outputText(mbps, "Mbps")
	}

	os.Exit(0)
}
