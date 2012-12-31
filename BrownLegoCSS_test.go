/*
 * Tests for "BrownLegoCSS"
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
 * The test files in the ./tests/ directory were obtained from
 * https://github.com/yui/yuicompressor/blob/master/tests/
 * and are the property of their respective creators.
 *
 */

package BrownLegoCSS

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testExtractDataUris(t *testing.T) {
	var tests = []struct {
		s, want         string
		preservedTokens []string
	}{
		{
			"background-image: url('elephant.png');",
			"background-image: url('elephant.png');",
			make([]string, 0),
		},
		{
			"background-image: url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQAQMAAAAlPW0iAAAABlBMVEUAAAD///+l2Z/dAAAAM0lEQVR4nGP4/5/h/1+G/58ZDrAz3D/McH8yw83NDDeNGe4Ug9C9zwz3gVLMDA/A6P9/AFGGFyjOXZtQAAAAAElFTkSuQmCC');",
			"background-image: url(___YUICSSMIN_PRESERVED_TOKEN_0___);",
			[]string{"'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQAQMAAAAlPW0iAAAABlBMVEUAAAD///+l2Z/dAAAAM0lEQVR4nGP4/5/h/1+G/58ZDrAz3D/McH8yw83NDDeNGe4Ug9C9zwz3gVLMDA/A6P9/AFGGFyjOXZtQAAAAAElFTkSuQmCC'"},
		},
	}

	for i, c := range tests {
		compressor := CssCompressor{}
		compressor.Css = []byte(tests[i].s)
		compressor.extractDataUris()
		got := string(compressor.Css)
		if got != c.want {
			t.Errorf("css: extractDataUris(%q) == (%q), want %q", c.s, got, c.want)
		}
		if fmt.Sprintf("%v", compressor.preservedTokens) != fmt.Sprintf("%v", c.preservedTokens) {
			t.Errorf("preservedTokens: extractDataUris(%q) == (%q), want %q", c.s, compressor.preservedTokens, c.preservedTokens)
		}
	}
}

func testExtractComments(t *testing.T) {
	var tests = []struct {
		s, want  string
		comments []string
	}{
		{
			"No comments here",
			"No comments here",
			make([]string, 0),
		},
		{
			"body { margin: 0px; } /* giggle */ test",
			"body { margin: 0px; } /*___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_0___*/ test",
			[]string{" giggle "},
		},
		{
			"/* This is a comment yo */ body { margin: 0px; } /* giggle */",
			"/*___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_0___*/ body { margin: 0px; } /*___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_1___*/",
			[]string{" This is a comment yo ", " giggle "},
		},
		{
			"/** Slightly Strange comment **/ test",
			"/*___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_0___*/ test",
			[]string{"* Slightly Strange comment *"},
		},
	}

	for i, c := range tests {
		compressor := CssCompressor{}
		compressor.Css = []byte(tests[i].s)
		compressor.extractComments()
		got := string(compressor.Css)
		if got != c.want {
			t.Errorf("css: extractComments(%q) == (%q), want %q", c.s, got, c.want)
		}
		if fmt.Sprintf("%v", compressor.comments) != fmt.Sprintf("%v", c.comments) {
			t.Errorf("comments: extractComments(%q) == (%q), want %q", c.s, compressor.comments, c.comments)
		}
	}
}

func testExtractStrings(t *testing.T) {
	var tests = []struct {
		s, want         string
		preservedTokens []string
	}{
		{
			"No strings here hooray",
			"No strings here hooray",
			make([]string, 0),
		},
		{
			"123 this \"is a test\" blah",
			"123 this \"___YUICSSMIN_PRESERVED_TOKEN_0___\" blah",
			[]string{"is a test"},
		},
		{
			"123 this \"is a \\\"test\\\"\" blah",
			"123 this \"___YUICSSMIN_PRESERVED_TOKEN_0___\" blah",
			[]string{"is a \\\"test\\\""},
		},
		{
			"123 this 'is a test' blah",
			"123 this '___YUICSSMIN_PRESERVED_TOKEN_0___' blah",
			[]string{"is a test"},
		},
		{
			"123 this 'is a \\'test\\'' blah",
			"123 this '___YUICSSMIN_PRESERVED_TOKEN_0___' blah",
			[]string{"is a \\'test\\'"},
		},
		{
			"mix \"them\" 'up' yo",
			"mix \"___YUICSSMIN_PRESERVED_TOKEN_0___\" '___YUICSSMIN_PRESERVED_TOKEN_1___' yo",
			[]string{"them", "up"},
		},
		{
			"\"progid:DXImageTransform.Microsoft.Alpha(Opacity=250\"",
			"\"___YUICSSMIN_PRESERVED_TOKEN_0___\"",
			[]string{"alpha(opacity=250"},
		},
	}
	for i, c := range tests {
		compressor := CssCompressor{}
		compressor.Css = []byte(tests[i].s)
		compressor.extractStrings()
		got := string(compressor.Css)
		if got != c.want {
			t.Errorf("css: extractStrings(%q) == (%q), want %q", c.s, got, c.want)
		}
		if fmt.Sprintf("%v", compressor.preservedTokens) != fmt.Sprintf("%v", c.preservedTokens) {
			t.Errorf("preservedTokens: extractStrings(%q) == (%q), want %q", c.s, compressor.preservedTokens, c.preservedTokens)
		}
	}
}

func testProcessComments(t *testing.T) {
	var tests = []struct {
		s, want string
	}{
		{"No comment", "No comment"},
		{"/* Normal comment */ test", " test"},
		{"/** Slightly Strange comment **/ test", " test"},
		{"123 /*! Preserve comment */", "123 /*___YUICSSMIN_PRESERVED_TOKEN_0___*/"},
		{"123 /* Multiple */ comment /* Comments */", "123  comment "},
		{"/* Hack comment \\*/ kek /**/ /* asdf */", "/*___YUICSSMIN_PRESERVED_TOKEN_0___*/ kek /*___YUICSSMIN_PRESERVED_TOKEN_1___*/ "},
		{"/* Hack comment \\*/ kek /* asdf */", "/*___YUICSSMIN_PRESERVED_TOKEN_0___*/ kek /*___YUICSSMIN_PRESERVED_TOKEN_1___*/"},
	}
	for i, c := range tests {
		compressor := CssCompressor{}
		compressor.Css = []byte(tests[i].s)
		compressor.extractComments()
		compressor.processComments()
		got := string(compressor.Css)
		if got != c.want {
			t.Errorf("css: processComments(%q) == (%q), want %q", c.s, got, c.want)
		}
	}
}

var testFiles []string

func testEverything(t *testing.T) {
	filepath.Walk("./tests/", visit)

	for _, f := range testFiles {
		minFile := "./tests/" + f + ".min"
		if _, err := os.Stat(minFile); os.IsNotExist(err) {
			continue
		}
		t.Logf("Now Testing %s ...\n", f)

		testContents, _ := ioutil.ReadFile("./tests/" + f)
		compareContents, _ := ioutil.ReadFile(minFile)

		compareContents = bytes.TrimSpace(compareContents)

		compressor := CssCompressor{Css: testContents}
		results := compressor.Compress()
		if results != string(compareContents) {
			t.Logf("%s\n", minFile)
			t.Logf("Attempting to compare\n%s\nwith\n%s\n...\n", results, compareContents)
			t.Errorf("testEverything: %q's contents do not match the results", f)
		}
	}
}

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		name := f.Name()
		if strings.LastIndex(name, ".css") == len(name)-4 {
			testFiles = append(testFiles, name)
		}
	}
	return nil
}

func Test(t *testing.T) {
	// First we test the individual extraction functions...
	testExtractDataUris(t)
	testExtractComments(t)
	testExtractStrings(t)
	testProcessComments(t)

	// ...then we test the full compilation using various test-cases
	// from https://github.com/yui/yuicompressor/blob/master/tests/
	testEverything(t)
}
