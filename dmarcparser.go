package main

import (
	"flag"
	"fmt"
	"github.com/gc1code/dmarcparser/dmarc"
	"os"
)

func main() {

	var filename = flag.String("filename", "none", "Path to DMARC Aggregate Report XML")

	flag.Parse()

	if *filename == "none" {
		fmt.Printf("XML File required.\n")
		return
	}

	xmlFile, err := os.Open(*filename) // For read access.
	if err != nil {
		fmt.Printf("os error: %v\n", err)
		return
	}
	defer xmlFile.Close()

	feedbackReport := dmarc.ParseReader(xmlFile)
	fmt.Printf("XMLName: %#v\n", feedbackReport)
}
