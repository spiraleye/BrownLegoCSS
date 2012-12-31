package main

import (
	"flag"
	"fmt"
	"github.com/spiraleye/BrownLegoCSS"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	var outfile string
	flag.StringVar(&outfile, "o", "", "Target file for minified output")
	flag.Parse()

	// Let the concurrency magic begin.
	dataChan := make(chan string)
	results := []string{}

	go func() {
		for _, infile := range flag.Args() {
			contents, err := ioutil.ReadFile(infile)
			if err != nil {
				fmt.Printf("Error reading file %s: %s\n", infile, err)
			} else {
				compressor := BrownLegoCSS.CssCompressor{Css: contents}
				compressed := compressor.Compress()
				dataChan <- compressed
			}
		}
		close(dataChan)
	}()

	for stringy := range dataChan {
		results = append(results, stringy)
	}

	outputString := strings.Join(results, "")

	ioutil.WriteFile(outfile, []byte(outputString), 0644)

	os.Exit(0)
}
