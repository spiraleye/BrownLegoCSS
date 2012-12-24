package BrownLegoCSS

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	//"strings"
)

type CssCompressor struct {
	Css             []byte
	preservedTokens []string
}

func (c *CssCompressor) extractDataUris() {
	re, _ := regexp.Compile("url\\(\\s*([\"']?)data\\:")
	re2, _ := regexp.Compile("\\s+")

	maxIndex := len(c.Css) - 1
	appendIndex := 0

	tmpCss := c.Css
	var sb bytes.Buffer

	for match := re.Find(tmpCss); match != nil; match = re.Find(tmpCss) {
		index := re.FindIndex(tmpCss)
		startIndex := index[0] + 4 // length of string "url("
		endIndex := index[1] - 1

		// The below is (hopefully) the equivalent of Java's Matcher m.group(1)
		terminator := re.FindStringSubmatch(string(tmpCss))[1] // ', " or empty (not quoted)

		if len(terminator) == 0 {
			terminator = ")"
		}

		foundTerminator := false
		for foundTerminator == false && endIndex+1 <= maxIndex {
			endIndex = bytes.IndexByte(tmpCss[endIndex+1:], terminator[0]) + len(tmpCss[:endIndex]) + 1
			if (endIndex > 0) && (string(tmpCss[endIndex-1]) != "\\") {
				foundTerminator = true
				if terminator != ")" {
					endIndex = bytes.IndexByte(tmpCss[endIndex:], ')') + len(tmpCss[:endIndex])
				}
			}
		}

		// Enough searching, start moving stuff over to the buffer
		sb.WriteString(string(tmpCss[appendIndex:index[0]]))

		if foundTerminator {
			var token string = string(tmpCss[startIndex:endIndex])
			token = re2.ReplaceAllString(token, "")
			c.preservedTokens = append(c.preservedTokens, token)

			preserver := "url(___YUICSSMIN_PRESERVED_TOKEN_" + strconv.Itoa(len(c.preservedTokens)-1) + "___)"
			sb.WriteString(preserver)

			appendIndex = endIndex + 1
		} else {
			// No end terminator found, re-add the whole match. Should we throw/warn here?
			sb.WriteString(string(tmpCss[index[0]:index[1]]))
			appendIndex = index[1]
		}

		sb.WriteString(string(tmpCss[appendIndex:]))

		tmpCss = tmpCss[appendIndex:]
	}

	// Our string buffer is not empty, so something must have changed.
	if sb.Len() > 0 {
		c.Css = sb.Bytes()
	}

	return
}

func (c *CssCompressor) extractComments() {
	c.Css = []byte("/* This is a comment yo */ body { margin: 0px; } /* giggle */")

	re, _ := regexp.Compile("\\/\\*")

	match := re.FindAll(c.Css, -1)
	index := re.FindAllIndex(c.Css, -1)
	fmt.Printf("%s", match)
	fmt.Printf("%s", index)

	//tmpCss := string(c.Css)

	//for startIndex := strings.Index(tmpCss[startIndex:], "/*"); startIndex >= 0; startIndex = strings.Index(tmpCss[startIndex:], "/*") {
	//endIndex = strings.Index(tmpCss[startIndex+2:], "*/")
	//fmt.Println(endIndex)
	//startIndex += 2
	//}
}

func (c *CssCompressor) Compress() string {
	c.extractDataUris()
	c.extractComments()
	return "\n"
}
