package main

import (
	"flag"
	"fmt"
	"github.com/spiraleye/BrownLegoCSS"
	"io/ioutil"
	"os"
)

func main() {
	var filename string
	flag.StringVar(&filename, "input", "", "Input CSS file for compression")
	flag.StringVar(&filename, "i", "", "Input CSS file for compression (shorthand)")
	flag.Parse()

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("%s does not exist: %s\n", filename, err)
		os.Exit(1)
	}
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}
	compressor := BrownLegoCSS.CssCompressor{Css: contents}
	fmt.Printf("%s", compressor.Compress())
}
