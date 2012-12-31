/*
 * Application for "BrownLegoCSS"
 * https://github.com/spiraleye/BrownLegoCSS
 * Author: Johan Meiring
 *
 * The copyrights embodied in the content of this file are licensed
 * by Spiraleye Studios under the 3-Clause BSD open source license,
 * as follows:
 *
 * Copyright (c) 2012, Spiraleye Studios
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *     * Neither the name of Spiraleye Studios nor the
 *       names of its contributors may be used to endorse or promote products
 *       derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL Spiraleye Studios BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

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

	// Dynamic gofunc to compress multiple files' contents at the same time.
	go func() {
		for _, infile := range flag.Args() {
			contents, err := ioutil.ReadFile(infile)
			if err != nil {
				fmt.Printf("Error reading file %s: %s\n", infile, err)
			} else {
				compressor := BrownLegoCSS.CssCompressor{Css: contents}
				compressedString := compressor.Compress()
				dataChan <- compressedString
			}
		}
		close(dataChan)
	}()

	// Grab data from the channel and append it to our slice.
	for tmpString := range dataChan {
		results = append(results, tmpString)
	}

	outputString := strings.Join(results, "")
	ioutil.WriteFile(outfile, []byte(outputString), 0644)
	os.Exit(0)
}
