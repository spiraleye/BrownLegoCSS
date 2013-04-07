/*
 * BrownLegoCSS
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
 * This work is derived from work that was originally done by Yahoo! inc:
 * http://developer.yahoo.com/yui/compressor/
 * Author: Julien Lecomte -  http://www.julienlecomte.net/
 * Author: Isaac Schlueter - http://foohack.com/
 * Author: Stoyan Stefanov - http://phpied.com/
 * Copyright (c) 2011 Yahoo! Inc.  All rights reserved.
 * The copyrights embodied in the content of this file are licensed
 * by Yahoo! Inc. under the BSD (revised) open source license.
 *
 * HOWEVER, THE DERIVED WORK IS NOT ENDORSED BY THE ABOVE ORGANISATION AND/OR
 * INDIVIDUALS IN ANY WAY, SHAPE, OR FORM.
 *
 */

package BrownLegoCSS

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type CssCompressor struct {
	Css             []byte
	preservedTokens []string
	comments        []string
}

func (c *CssCompressor) extractDataUris() {
	re, _ := regexp.Compile("url\\(\\s*([\"']?)data\\:")
	re2, _ := regexp.Compile("\\s+")

	var sb bytes.Buffer
	appendIndex := 0
	maxIndex := len(c.Css) - 1

	indexes := re.FindAllIndex(c.Css, -1)
	submatches := re.FindAllStringSubmatch(string(c.Css), -1)

	for counter, i := range indexes {
		startIndex := i[0] + 4 // "url(".length()
		endIndex := i[1] - 1

		terminator := submatches[counter][1] // ', " or empty (not quoted)
		if len(terminator) == 0 {
			terminator = ")"
		}
		foundTerminator := false

		for foundTerminator == false && endIndex+1 <= maxIndex {
			endIndex = bytes.IndexByte(c.Css[endIndex+1:], terminator[0]) + len(c.Css[:endIndex]) + 1
			if (endIndex > 0) && (string(c.Css[endIndex-1]) != "\\") {
				foundTerminator = true
				if terminator != ")" {
					endIndex = bytes.IndexByte(c.Css[endIndex:], ')') + len(c.Css[:endIndex])
				}
			}
		}

		// Enough searching, start moving stuff over to the buffer
		sb.WriteString(string(c.Css[appendIndex:i[0]]))

		if foundTerminator {
			var token string = string(c.Css[startIndex:endIndex])
			token = re2.ReplaceAllString(token, "")
			c.preservedTokens = append(c.preservedTokens, token)

			preserver := "url(___YUICSSMIN_PRESERVED_TOKEN_" + strconv.Itoa(len(c.preservedTokens)-1) + "___)"
			sb.WriteString(preserver)

			appendIndex = endIndex + 1
		} else {
			// No end terminator found, re-add the whole match. Should we throw/warn here?
			sb.WriteString(string(c.Css[i[0]:i[1]]))
			appendIndex = i[1]
		}
	}
	if appendIndex > 0 {
		sb.WriteString(string(c.Css[appendIndex:]))
	}
	// Our string buffer is not empty, so something must have changed.
	if sb.Len() > 0 {
		c.Css = sb.Bytes()
	}
}

func (c *CssCompressor) extractComments() {
	var sb bytes.Buffer
	startIndex := 0
	endIndex := 0

	tmpCss := c.Css
	for startIndex = bytes.Index(tmpCss, []byte("/*")); startIndex >= 0; {
		sb.WriteString(string(tmpCss[:startIndex]))

		endIndex = bytes.Index(tmpCss[startIndex+2:], []byte("*/"))
		if endIndex < 0 {
			endIndex = len(tmpCss)
		}
		c.comments = append(c.comments, string(tmpCss[startIndex+2:endIndex+startIndex+2]))
		sb.WriteString(
			string("/*___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_" +
				(strconv.Itoa(len(c.comments) - 1)) +
				"___*/"))

		tmpCss = tmpCss[startIndex+2+endIndex+2:]
		startIndex = bytes.Index(tmpCss, []byte("/*"))
	}
	sb.WriteString(string(tmpCss))
	c.Css = sb.Bytes()
}

func (c *CssCompressor) extractStrings() {
	re, _ := regexp.Compile("(\"([^\\\\\"]|\\\\.|\\\\)*\")|('([^\\\\']|\\\\.|\\\\)*')")
	re2, _ := regexp.Compile("(?i)progid:DXImageTransform.Microsoft.Alpha\\(Opacity=")

	var sb bytes.Buffer
	previousIndex := 0

	indexes := re.FindAllIndex(c.Css, -1)
	tokens := re.FindAllStringSubmatch(string(c.Css), -1)

	for counter, i := range indexes {
		if i[0] > 0 {
			sb.WriteString(string(c.Css[previousIndex:i[0]]))
		}
		token := tokens[counter][0]
		quote := token[0]
		token = token[1 : len(token)-1]

		// maybe the string contains a comment-like substring?
		// one, maybe more? put'em back then
		if strings.Index(token, "___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_") >= 0 {
			max := len(c.comments)
			for j := 0; j < max; j += 1 {
				token = strings.Replace(token, "___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_"+strconv.Itoa(j)+"___", c.comments[j], -1)
			}
		}

		// minify alpha opacity in filter strings
		token = re2.ReplaceAllString(token, "alpha(opacity=")

		c.preservedTokens = append(c.preservedTokens, token)
		preserver := string(quote) + "___YUICSSMIN_PRESERVED_TOKEN_" + strconv.Itoa(len(c.preservedTokens)-1) + "___" + string(quote)
		sb.WriteString(preserver)

		previousIndex = i[1]
	}
	if previousIndex > 0 {
		sb.WriteString(string(c.Css[previousIndex:]))
	}
	// Our string buffer is not empty, so something must have changed.
	if sb.Len() > 0 {
		c.Css = sb.Bytes()
	}
}

func (c *CssCompressor) processComments() {
	max := len(c.comments)

	for i := 0; i < max; i += 1 {
		token := c.comments[i]
		placeholder := "___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_" + strconv.Itoa(i) + "___"

		// ! in the first position of the comment means preserve
		// so push to the preserved tokens while stripping the !
		if strings.Index(token, "!") == 0 {
			c.preservedTokens = append(c.preservedTokens, token)
			c.Css = bytes.Replace(c.Css, []byte(placeholder), []byte("___YUICSSMIN_PRESERVED_TOKEN_"+strconv.Itoa(len(c.preservedTokens)-1)+"___"), -1)
			continue
		}

		// \ in the last position looks like hack for Mac/IE5/Opera
		// shorten that to /*\*/ and the next one to /**/
		// TODO: this doesn't seem to be working as intended, even in the Java version.
		if token != "" && strings.LastIndex(token, "\\") == len(token)-1 {
			c.preservedTokens = append(c.preservedTokens, "\\")
			c.Css = bytes.Replace(c.Css, []byte(placeholder), []byte("___YUICSSMIN_PRESERVED_TOKEN_"+strconv.Itoa(len(c.preservedTokens)-1)+"___"), -1)

			i = i + 1 // attn: advancing the loop
			c.preservedTokens = append(c.preservedTokens, "")
			c.Css = bytes.Replace(c.Css, []byte("___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_"+strconv.Itoa(i)+"___"), []byte("___YUICSSMIN_PRESERVED_TOKEN_"+strconv.Itoa(len(c.preservedTokens)-1)+"___"), -1)

			continue
		}

		// keep empty comments after child selectors (IE7 hack)
		// e.g. html >/**/ body
		if len(token) == 0 {
			startIndex := bytes.Index(c.Css, []byte(placeholder))
			if startIndex > 2 {
				if c.Css[startIndex-3] == '>' {
					c.preservedTokens = append(c.preservedTokens, "")
					c.Css = bytes.Replace(c.Css, []byte(placeholder), []byte("___YUICSSMIN_PRESERVED_TOKEN_"+strconv.Itoa(len(c.preservedTokens)-1)+"___"), -1)
				}
			}
		}

		// in all other cases kill the comment
		c.Css = bytes.Replace(c.Css, []byte("/*"+placeholder+"*/"), []byte(""), -1)
	}
}

func (c *CssCompressor) performGeneralCleanup() {
	// This function does a lot, ok?

	// Remove the spaces before the things that should not have spaces before them.
	// But, be careful not to turn "p :link {...}" into "p:link{...}"
	// Swap out any pseudo-class colons with the token, and then swap back.
	re, _ := regexp.Compile("(^|\\})(([^\\{:])+:)+([^\\{]*\\{)")
	var sb bytes.Buffer
	previousIndex := 0
	indexes := re.FindAllIndex(c.Css, -1)
	groups := re.FindAllStringSubmatch(string(c.Css), -1)
	for counter, i := range indexes {
		if i[0] > 0 {
			sb.WriteString(string(c.Css[previousIndex:i[0]]))
		}
		s := groups[counter][0]
		s = strings.Replace(s, ":", "___YUICSSMIN_PSEUDOCLASSCOLON___", -1)
		s = strings.Replace(s, "\\\\", "\\\\\\\\", -1)
		s = strings.Replace(s, "\\$", "\\\\\\$", -1)
		sb.WriteString(s)
		previousIndex = i[1]
	}
	if previousIndex > 0 {
		sb.WriteString(string(c.Css[previousIndex:]))
	}
	if sb.Len() > 0 {
		c.Css = sb.Bytes()
	}

	// Remove spaces before the things that should not have spaces before them.
	re, _ = regexp.Compile("\\s+([!{};:>+\\(\\)\\],])")
	c.Css = re.ReplaceAll(c.Css, []byte("$1"))
	// bring back the colon
	c.Css = bytes.Replace(c.Css, []byte("___YUICSSMIN_PSEUDOCLASSCOLON___"), []byte(":"), -1)

	// retain space for special IE6 cases
	re, _ = regexp.Compile(":first\\-(line|letter)(\\{|,)")
	c.Css = re.ReplaceAll(c.Css, []byte(":first-$1 $2"))

	// no space after the end of a preserved comment
	c.Css = bytes.Replace(c.Css, []byte("*/ "), []byte("*/"), -1)

	// If there is a @charset, then only allow one, and push to the top of the file.
	re, _ = regexp.Compile("^(.*)(@charset \"[^\"]*\";)")
	c.Css = re.ReplaceAll(c.Css, []byte("$2$1"))
	re, _ = regexp.Compile("^(\\s*@charset [^;]+;\\s*)+")
	c.Css = re.ReplaceAll(c.Css, []byte("$1"))

	// Put the space back in some cases, to support stuff like
	// @media screen and (-webkit-min-device-pixel-ratio:0){
	re, _ = regexp.Compile("\\band\\(")
	c.Css = re.ReplaceAll(c.Css, []byte("and ("))

	// Remove the spaces after the things that should not have spaces after them.
	re, _ = regexp.Compile("([!{}:;>+\\(\\[,])\\s+")
	c.Css = re.ReplaceAll(c.Css, []byte("$1"))

	// remove unnecessary semicolons
	re, _ = regexp.Compile(";+}")
	c.Css = re.ReplaceAll(c.Css, []byte("}"))

	// Replace 0(px,em,%) with 0.
	re, _ = regexp.Compile("([\\s:])(0)(px|em|%|in|cm|mm|pc|pt|ex)")
	c.Css = re.ReplaceAll(c.Css, []byte("$1$2"))

	// Replace 0 0 0 0; with 0.
	re, _ = regexp.Compile(":0 0 0 0(;|})")
	re2, _ := regexp.Compile(":0 0 0(;|})")
	re3, _ := regexp.Compile(":0 0(;|})")
	c.Css = re.ReplaceAll(c.Css, []byte(":0$1"))
	c.Css = re2.ReplaceAll(c.Css, []byte(":0$1"))
	c.Css = re3.ReplaceAll(c.Css, []byte(":0$1"))

	// Replace background-position:0; with background-position:0 0;
	// same for transform-origin
	sb.Reset()
	re, _ = regexp.Compile("(?i)(background-position|transform-origin|webkit-transform-origin|moz-transform-origin|o-transform-origin|ms-transform-origin):0(;|})")
	previousIndex = 0
	indexes = re.FindAllIndex(c.Css, -1)
	groups = re.FindAllStringSubmatch(string(c.Css), -1)
	for counter, i := range indexes {
		if i[0] > 0 {
			sb.WriteString(string(c.Css[previousIndex:i[0]]))
		}
		s := strings.ToLower(groups[counter][1]) + ":0 0" + groups[counter][2]
		sb.WriteString(s)
		previousIndex = i[1]
	}
	if previousIndex > 0 {
		sb.WriteString(string(c.Css[previousIndex:]))
	}
	if sb.Len() > 0 {
		c.Css = sb.Bytes()
	}

	// Replace 0.6 to .6, but only when preceded by : or a white-space
	re, _ = regexp.Compile("(:|\\s)0+\\.(\\d+)")
	c.Css = re.ReplaceAll(c.Css, []byte("$1.$2"))

	// Shorten colors from rgb(51,102,153) to #336699
	// This makes it more likely that it'll get further compressed in the next step.
	sb.Reset()
	previousIndex = 0
	re, _ = regexp.Compile("rgb\\s*\\(\\s*([0-9,\\s]+)\\s*\\)")
	indexes = re.FindAllIndex(c.Css, -1)
	groups = re.FindAllStringSubmatch(string(c.Css), -1)
	for counter, i := range indexes {
		if i[0] > 0 {
			sb.WriteString(string(c.Css[previousIndex:i[0]]))
		}
		rgbcolors := strings.Split(groups[counter][1], ",")
		var hexcolor bytes.Buffer
		hexcolor.WriteString("#")
		for _, colour := range rgbcolors {
			val, _ := strconv.Atoi(colour)
			if val < 16 {
				hexcolor.WriteString("0")
			}
			hexcolor.WriteString(fmt.Sprintf("%x", val))
		}
		sb.WriteString(hexcolor.String())
		previousIndex = i[1]
	}
	if previousIndex > 0 {
		sb.WriteString(string(c.Css[previousIndex:]))
	}
	if sb.Len() > 0 {
		c.Css = sb.Bytes()
	}

	// Shorten colors from #AABBCC to #ABC. Note that we want to make sure
	// the color is not preceded by either ", " or =. Indeed, the property
	//     filter: chroma(color="#FFFFFF");
	// would become
	//     filter: chroma(color="#FFF");
	// which makes the filter break in IE.
	// We also want to make sure we're only compressing #AABBCC patterns inside { }, not id selectors ( #FAABAC {} )
	// We also want to avoid compressing invalid values (e.g. #AABBCCD to #ABCD)
	sb.Reset()
	re, _ = regexp.Compile("(\\=\\s*?[\"']?)?" + "#([0-9a-fA-F])([0-9a-fA-F])([0-9a-fA-F])([0-9a-fA-F])([0-9a-fA-F])([0-9a-fA-F])" + "(:?\\}|[^0-9a-fA-F{][^{]*?\\})")
	previousIndex = 0

	for match := re.Find(c.Css[previousIndex:]); match != nil; match = re.Find(c.Css[previousIndex:]) {
		index := re.FindIndex(c.Css[previousIndex:])
		submatches := re.FindStringSubmatch(string(c.Css[previousIndex:]))
		submatchIndexes := re.FindSubmatchIndex(c.Css[previousIndex:])

		sb.WriteString(string(c.Css[previousIndex : index[0]+len(c.Css[:previousIndex])]))

		//boolean isFilter = (m.group(1) != null && !"".equals(m.group(1)));
		// I hope the below is the equivalent of the above :P
		isFilter := submatches[1] != "" && submatchIndexes[1] != -1

		if isFilter {
			// Restore, as is. Compression will break filters
			sb.WriteString(submatches[1] + "#" + submatches[2] + submatches[3] + submatches[4] + submatches[5] + submatches[6] + submatches[7])
		} else {
			if strings.ToLower(submatches[2]) == strings.ToLower(submatches[3]) &&
				strings.ToLower(submatches[4]) == strings.ToLower(submatches[5]) &&
				strings.ToLower(submatches[6]) == strings.ToLower(submatches[7]) {
				// #AABBCC pattern
				sb.WriteString("#" + strings.ToLower(submatches[3]+submatches[5]+submatches[7]))
			} else {
				// Non-compressible color, restore, but lower case.
				sb.WriteString("#" + strings.ToLower(submatches[2]+submatches[3]+submatches[4]+submatches[5]+submatches[6]+submatches[7]))
			}
		}

		// The "+ 4" below is a crazy hack which will come back to haunt me later.
		// For now, it makes everything work 100%.
		previousIndex = submatchIndexes[7] + len(c.Css[:previousIndex]) + 4

	}
	if previousIndex > 0 {
		sb.WriteString(string(c.Css[previousIndex:]))
	}
	if sb.Len() > 0 {
		c.Css = sb.Bytes()
	}

	// border: none -> border:0
	re, _ = regexp.Compile("(?i)(border|border-top|border-right|border-bottom|border-left|outline|background):none(;|})")
	sb.Reset()
	previousIndex = 0
	indexes = re.FindAllIndex(c.Css, -1)
	groups = re.FindAllStringSubmatch(string(c.Css), -1)
	for counter, i := range indexes {
		if i[0] > 0 {
			sb.WriteString(string(c.Css[previousIndex:i[0]]))
		}
		s := strings.ToLower(groups[counter][1]) + ":0" + groups[counter][2]
		sb.WriteString(s)
		previousIndex = i[1]
	}
	if previousIndex > 0 {
		sb.WriteString(string(c.Css[previousIndex:]))
	}
	if sb.Len() > 0 {
		c.Css = sb.Bytes()
	}

	// shorter opacity IE filter
	re, _ = regexp.Compile("(?i)progid:DXImageTransform.Microsoft.Alpha\\(Opacity=")
	c.Css = re.ReplaceAll(c.Css, []byte("alpha(opacity="))

	// Remove empty rules.
	re, _ = regexp.Compile("[^\\}\\{/;]+\\{\\}")
	c.Css = re.ReplaceAll(c.Css, []byte(""))
}

func (c *CssCompressor) Compress() []byte {
	c.extractDataUris()
	c.extractComments()

	// preserve strings so their content doesn't get accidentally minified
	c.extractStrings()

	// strings are safe, now wrestle the comments
	c.processComments()

	// Normalize all whitespace strings to single spaces. Easier to work with that way.
	re, _ := regexp.Compile("\\s+")
	c.Css = re.ReplaceAll(c.Css, []byte(" "))

	// Do lots and lots and lots of fun things
	c.performGeneralCleanup()

	// Replace multiple semi-colons in a row with a single one
	re, _ = regexp.Compile(";;+")
	c.Css = re.ReplaceAll(c.Css, []byte(";"))

	// restore preserved comments and strings
	for i, t := range c.preservedTokens {
		c.Css = bytes.Replace(c.Css, []byte("___YUICSSMIN_PRESERVED_TOKEN_"+strconv.Itoa(i)+"___"), []byte(t), -1)
	}

	// Trim the final string (for any leading or trailing white spaces)
	c.Css = bytes.TrimSpace(c.Css)

	// Hooray, we're done!
	return c.Css
}
