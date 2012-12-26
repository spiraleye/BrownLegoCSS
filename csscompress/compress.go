package main

import (
	"flag"
	"fmt"
	"github.com/spiraleye/BrownLegoCSS/BrownLegoCSS"
	//"io/ioutil"
	//"os"
)

func main() {
	var filename string
	flag.StringVar(&filename, "input", "", "Input CSS file for compression")
	flag.StringVar(&filename, "i", "", "Input CSS file for compression (shorthand)")
	flag.Parse()

	/*if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("%s does not exist: %s\n", filename, err)
		os.Exit(1)
	}
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}*/
	//compressor := BrownLegoCSS.CssCompressor{Css: contents}
	compressor := BrownLegoCSS.CssCompressor{}
	fmt.Printf(compressor.Compress())

	//fmt.Printf("%s\n", strings.LastIndex("background-image: url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQAQMAAAAlPW0iAAAABlBMVEUAAAD///+l2Z/dAAAAM0lEQVR4nGP4/5/h/1+G/58ZDrAz3D/McH8yw83NDDeNGe4Ug9C9zwz3gVLMDA/A6P9/AFGGFyjOXZtQAAAAAElFTkSuQmCC');", "'"))
}
