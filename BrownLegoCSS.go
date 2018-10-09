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

package brownlegocss

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// CSSCompressor contains data needed for compressing CSS data.
type CSSCompressor struct {
	CSS             []byte
	preservedTokens []string
	comments        []string
}

func regexFindReplace(input []byte, regex string, handlefunc func(groups []string) string) []byte {
	var sb bytes.Buffer
	re, _ := regexp.Compile(regex)
	previousIndex := 0
	indexes := re.FindAllIndex(input, -1)
	groups := re.FindAllStringSubmatch(string(input), -1)
	for counter, i := range indexes {
		if i[0] > 0 {
			sb.WriteString(string(input[previousIndex:i[0]]))
		}
		s := handlefunc(groups[counter])
		sb.WriteString(s)
		previousIndex = i[1]
	}
	if previousIndex > 0 {
		sb.WriteString(string(input[previousIndex:]))
	}
	if sb.Len() > 0 {
		return sb.Bytes()
	}
	return input
}

func (c *CSSCompressor) extractDataUris() {
	re, _ := regexp.Compile("(?i)url\\(\\s*([\"']?)data\\:")
	re2, _ := regexp.Compile("\\s+")

	var sb bytes.Buffer
	appendIndex := 0
	maxIndex := len(c.CSS) - 1

	indexes := re.FindAllIndex(c.CSS, -1)
	submatches := re.FindAllStringSubmatch(string(c.CSS), -1)

	for counter, i := range indexes {
		startIndex := i[0] + 4 // "url(".length()
		endIndex := i[1] - 1

		terminator := submatches[counter][1] // ', " or empty (not quoted)
		if len(terminator) == 0 {
			terminator = ")"
		}
		foundTerminator := false

		for foundTerminator == false && endIndex+1 <= maxIndex {
			endIndex = bytes.IndexByte(c.CSS[endIndex+1:], terminator[0]) + len(c.CSS[:endIndex]) + 1
			if (endIndex > 0) && (string(c.CSS[endIndex-1]) != "\\") {
				foundTerminator = true
				if terminator != ")" {
					endIndex = bytes.IndexByte(c.CSS[endIndex:], ')') + len(c.CSS[:endIndex])
				}
			}
		}

		// Enough searching, start moving stuff over to the buffer
		sb.WriteString(string(c.CSS[appendIndex:i[0]]))

		if foundTerminator {
			token := string(c.CSS[startIndex:endIndex])
			token = re2.ReplaceAllString(token, "")
			c.preservedTokens = append(c.preservedTokens, token)

			preserver := "url(___YUICSSMIN_PRESERVED_TOKEN_" + strconv.Itoa(len(c.preservedTokens)-1) + "___)"
			sb.WriteString(preserver)

			appendIndex = endIndex + 1
		} else {
			// No end terminator found, re-add the whole match. Should we throw/warn here?
			sb.WriteString(string(c.CSS[i[0]:i[1]]))
			appendIndex = i[1]
		}
	}
	if appendIndex > 0 {
		sb.WriteString(string(c.CSS[appendIndex:]))
	}
	// Our string buffer is not empty, so something must have changed.
	if sb.Len() > 0 {
		c.CSS = sb.Bytes()
	}
}

func (c *CSSCompressor) extractComments() {
	var sb bytes.Buffer
	startIndex := 0
	endIndex := 0

	tmpCSS := c.CSS
	for startIndex = bytes.Index(tmpCSS, []byte("/*")); startIndex >= 0; {
		sb.WriteString(string(tmpCSS[:startIndex]))

		endIndex = bytes.Index(tmpCSS[startIndex+2:], []byte("*/"))
		if endIndex < 0 {
			endIndex = len(tmpCSS)
		}
		c.comments = append(c.comments, string(tmpCSS[startIndex+2:endIndex+startIndex+2]))
		sb.WriteString(
			string("/*___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_" +
				(strconv.Itoa(len(c.comments) - 1)) +
				"___*/"))

		tmpCSS = tmpCSS[startIndex+2+endIndex+2:]
		startIndex = bytes.Index(tmpCSS, []byte("/*"))
	}
	sb.WriteString(string(tmpCSS))
	c.CSS = sb.Bytes()
}

func (c *CSSCompressor) extractStrings() {
	re, _ := regexp.Compile("(\"([^\\\\\"]|\\\\.|\\\\)*\")|('([^\\\\']|\\\\.|\\\\)*')")
	re2, _ := regexp.Compile("(?i)progid:DXImageTransform.Microsoft.Alpha\\(Opacity=")

	var sb bytes.Buffer
	previousIndex := 0

	indexes := re.FindAllIndex(c.CSS, -1)
	tokens := re.FindAllStringSubmatch(string(c.CSS), -1)

	for counter, i := range indexes {
		if i[0] > 0 {
			sb.WriteString(string(c.CSS[previousIndex:i[0]]))
		}
		token := tokens[counter][0]
		quote := token[0]
		token = token[1 : len(token)-1]

		// maybe the string contains a comment-like substring?
		// one, maybe more? put'em back then
		if strings.Index(token, "___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_") >= 0 {
			max := len(c.comments)
			for j := 0; j < max; j++ {
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
		sb.WriteString(string(c.CSS[previousIndex:]))
	}
	// Our string buffer is not empty, so something must have changed.
	if sb.Len() > 0 {
		c.CSS = sb.Bytes()
	}
}

func (c *CSSCompressor) processComments() {
	max := len(c.comments)

	for i := 0; i < max; i++ {
		token := c.comments[i]
		placeholder := "___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_" + strconv.Itoa(i) + "___"

		// ! in the first position of the comment means preserve
		// so push to the preserved tokens while stripping the !
		if strings.Index(token, "!") == 0 {
			c.preservedTokens = append(c.preservedTokens, token)
			c.CSS = bytes.Replace(c.CSS, []byte(placeholder), []byte("___YUICSSMIN_PRESERVED_TOKEN_"+strconv.Itoa(len(c.preservedTokens)-1)+"___"), -1)
			continue
		}

		// \ in the last position looks like hack for Mac/IE5/Opera
		// shorten that to /*\*/ and the next one to /**/
		// TODO: this doesn't seem to be working as intended, even in the Java version.
		if token != "" && strings.LastIndex(token, "\\") == len(token)-1 {
			c.preservedTokens = append(c.preservedTokens, "\\")
			c.CSS = bytes.Replace(c.CSS, []byte(placeholder), []byte("___YUICSSMIN_PRESERVED_TOKEN_"+strconv.Itoa(len(c.preservedTokens)-1)+"___"), -1)

			i = i + 1 // attn: advancing the loop
			c.preservedTokens = append(c.preservedTokens, "")
			c.CSS = bytes.Replace(c.CSS, []byte("___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_"+strconv.Itoa(i)+"___"), []byte("___YUICSSMIN_PRESERVED_TOKEN_"+strconv.Itoa(len(c.preservedTokens)-1)+"___"), -1)

			continue
		}

		// keep empty comments after child selectors (IE7 hack)
		// e.g. html >/**/ body
		if len(token) == 0 {
			startIndex := bytes.Index(c.CSS, []byte(placeholder))
			if startIndex > 2 {
				if c.CSS[startIndex-3] == '>' {
					c.preservedTokens = append(c.preservedTokens, "")
					c.CSS = bytes.Replace(c.CSS, []byte(placeholder), []byte("___YUICSSMIN_PRESERVED_TOKEN_"+strconv.Itoa(len(c.preservedTokens)-1)+"___"), -1)
				}
			}
		}

		// in all other cases kill the comment
		c.CSS = bytes.Replace(c.CSS, []byte("/*"+placeholder+"*/"), []byte(""), -1)
	}
}

func (c *CSSCompressor) performGeneralCleanup() {
	// This function does a lot, ok?
	var sb bytes.Buffer
	var previousIndex int
	var re *regexp.Regexp

	// Remove the spaces before the things that should not have spaces before them.
	// But, be careful not to turn "p :link {...}" into "p:link{...}"
	// Swap out any pseudo-class colons with the token, and then swap back.
	c.CSS = regexFindReplace(c.CSS,
		"(^|\\})(([^\\{:])+:)+([^\\{]*\\{)",
		func(groups []string) string {
			s := groups[0]
			s = strings.Replace(s, ":", "___YUICSSMIN_PSEUDOCLASSCOLON___", -1)
			s = strings.Replace(s, "\\\\", "\\\\\\\\", -1)
			s = strings.Replace(s, "\\$", "\\\\\\$", -1)
			return s
		})

	// Remove spaces before the things that should not have spaces before them.
	re, _ = regexp.Compile("\\s+([!{};:>+\\(\\)\\],])")
	c.CSS = re.ReplaceAll(c.CSS, []byte("$1"))
	// Restore spaces for !important
	c.CSS = bytes.Replace(c.CSS, []byte("!important"), []byte(" !important"), -1)
	// bring back the colon
	c.CSS = bytes.Replace(c.CSS, []byte("___YUICSSMIN_PSEUDOCLASSCOLON___"), []byte(":"), -1)

	// retain space for special IE6 cases
	c.CSS = regexFindReplace(c.CSS,
		"(?i):first\\-(line|letter)(\\{|,)",
		func(groups []string) string {
			return strings.ToLower(":first-"+groups[1]) + " " + groups[2]
		})

	// no space after the end of a preserved comment
	c.CSS = bytes.Replace(c.CSS, []byte("*/ "), []byte("*/"), -1)

	// If there are multiple @charset directives, push them to the top of the file.
	c.CSS = regexFindReplace(c.CSS,
		"(?i)^(.*)(@charset)( \"[^\"]*\";)",
		func(groups []string) string {
			return strings.ToLower(groups[2]) + groups[3] + groups[1]
		})

	// When all @charset are at the top, remove the second and after (as they are completely ignored).
	c.CSS = regexFindReplace(c.CSS,
		"(?i)^((\\s*)(@charset)( [^;]+;\\s*))+",
		func(groups []string) string {
			return groups[2] + strings.ToLower(groups[3]) + groups[4]
		})

	// lowercase some popular @directives
	c.CSS = regexFindReplace(c.CSS,
		"(?i)@(charset|font-face|import|(?:-(?:atsc|khtml|moz|ms|o|wap|webkit)-)?keyframe|media|page|namespace)",
		func(groups []string) string {
			return "@" + strings.ToLower(groups[1])
		})

	// lowercase some more common pseudo-elements
	c.CSS = regexFindReplace(c.CSS,
		"(?i):(active|after|before|checked|disabled|empty|enabled|first-(?:child|of-type)|focus|hover|last-(?:child|of-type)|link|only-(?:child|of-type)|root|:selection|target|visited)",
		func(groups []string) string {
			return ":" + strings.ToLower(groups[1])
		})

	// lowercase some more common functions
	c.CSS = regexFindReplace(c.CSS,
		"(?i):(lang|not|nth-child|nth-last-child|nth-last-of-type|nth-of-type|(?:-(?:moz|webkit)-)?any)\\(",
		func(groups []string) string {
			return ":" + strings.ToLower(groups[1]) + "("
		})

	// lower case some common function that can be values
	// NOTE: rgb() isn't useful as we replace with #hex later, as well as and() is already done for us right after this
	c.CSS = regexFindReplace(c.CSS,
		"(?i)([:,\\( ]\\s*)(attr|color-stop|from|rgba|to|url|(?:-(?:atsc|khtml|moz|ms|o|wap|webkit)-)?(?:calc|max|min|(?:repeating-)?(?:linear|radial)-gradient)|-webkit-gradient)",
		func(groups []string) string {
			return groups[1] + strings.ToLower(groups[2])
		})

	// Put the space back in some cases, to support stuff like
	// @media screen and (-webkit-min-device-pixel-ratio:0){
	re, _ = regexp.Compile("(?i)\\band\\(")
	c.CSS = re.ReplaceAll(c.CSS, []byte("and ("))

	// Remove the spaces after the things that should not have spaces after them.
	re, _ = regexp.Compile("([!{}:;>+\\(\\[,])\\s+")
	c.CSS = re.ReplaceAll(c.CSS, []byte("$1"))

	// remove unnecessary semicolons
	re, _ = regexp.Compile(";+}")
	c.CSS = re.ReplaceAll(c.CSS, []byte("}"))

	// Replace 0(px,em,%) with 0.
	re, _ = regexp.Compile("(?i)(^|[^0-9])(?:0?\\.)?0(?:px|em|%|in|cm|mm|pc|pt|ex|deg|g?rad|m?s|k?hz)")
	c.CSS = re.ReplaceAll(c.CSS, []byte("${1}0"))

	// Replace 0 0 0 0; with 0.
	re, _ = regexp.Compile(":0 0 0 0(;|})")
	re2, _ := regexp.Compile(":0 0 0(;|})")
	re3, _ := regexp.Compile(":0 0(;|})")
	c.CSS = re.ReplaceAll(c.CSS, []byte(":0$1"))
	c.CSS = re2.ReplaceAll(c.CSS, []byte(":0$1"))
	c.CSS = re3.ReplaceAll(c.CSS, []byte(":0$1"))

	// Replace background-position:0; with background-position:0 0;
	// same for transform-origin
	c.CSS = regexFindReplace(c.CSS,
		"(?i)(background-position|webkit-mask-position|transform-origin|webkit-transform-origin|moz-transform-origin|o-transform-origin|ms-transform-origin):0(;|})",
		func(groups []string) string {
			return strings.ToLower(groups[1]) + ":0 0" + groups[2]
		})

	// Replace 0.6 to .6, but only when preceded by : or a white-space
	re, _ = regexp.Compile("(:|\\s)0+\\.(\\d+)")
	c.CSS = re.ReplaceAll(c.CSS, []byte("$1.$2"))

	// Shorten colors from rgb(51,102,153) to #336699
	// This makes it more likely that it'll get further compressed in the next step.
	c.CSS = regexFindReplace(c.CSS,
		"rgb\\s*\\(\\s*([0-9,\\s]+)\\s*\\)",
		func(groups []string) string {
			rgbcolors := strings.Split(groups[1], ",")
			var hexcolor bytes.Buffer
			hexcolor.WriteString("#")
			for _, colour := range rgbcolors {
				val, _ := strconv.Atoi(colour)
				if val < 16 {
					hexcolor.WriteString("0")
				}
				// If someone passes an RGB value that's too big to express in two characters, round down.
				// Probably should throw out a warning here, but generating valid CSS is a bigger concern.
				if val > 255 {
					val = 255
				}
				hexcolor.WriteString(fmt.Sprintf("%x", val))
			}
			return hexcolor.String()
		})

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

	for match := re.Find(c.CSS[previousIndex:]); match != nil; match = re.Find(c.CSS[previousIndex:]) {
		index := re.FindIndex(c.CSS[previousIndex:])
		submatches := re.FindStringSubmatch(string(c.CSS[previousIndex:]))
		submatchIndexes := re.FindSubmatchIndex(c.CSS[previousIndex:])

		sb.WriteString(string(c.CSS[previousIndex : index[0]+len(c.CSS[:previousIndex])]))

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
		previousIndex = submatchIndexes[7] + len(c.CSS[:previousIndex]) + 4

	}
	if previousIndex > 0 {
		sb.WriteString(string(c.CSS[previousIndex:]))
	}
	if sb.Len() > 0 {
		c.CSS = sb.Bytes()
	}

	// Save a few chars by utilizing short colour keywords.
	// https://github.com/yui/yuicompressor/commit/fe8cf35d3693910103d65bf465d33b0d602dcfea
	colours := map[string]string{
		"#f00":    "red",
		"#000080": "navy",
		"#808080": "gray",
		"#808000": "olive",
		"#800080": "purple",
		"#c0c0c0": "silver",
		"#008080": "teal",
		"#ffa500": "orange",
		"#800000": "maroon",
	}
	for k, v := range colours {
		re, _ = regexp.Compile("(:|\\s)" + k + "(;|})")
		c.CSS = re.ReplaceAll(c.CSS, []byte("${1}"+v+"${2}"))
	}

	// border: none -> border:0
	c.CSS = regexFindReplace(c.CSS,
		"(?i)(border|border-top|border-right|border-bottom|border-left|outline|background):none(;|})",
		func(groups []string) string {
			return strings.ToLower(groups[1]) + ":0" + groups[2]
		})

	// shorter opacity IE filter
	re, _ = regexp.Compile("(?i)progid:DXImageTransform.Microsoft.Alpha\\(Opacity=")
	c.CSS = re.ReplaceAll(c.CSS, []byte("alpha(opacity="))

	// Find a fraction that is used for Opera's -o-device-pixel-ratio query
	// Add token to add the "\" back in later
	re, _ = regexp.Compile("\\(([\\-A-Za-z]+):([0-9]+)\\/([0-9]+)\\)")
	c.CSS = re.ReplaceAll(c.CSS, []byte("(${1}:${2}___YUI_QUERY_FRACTION___${3})"))

	// Remove empty rules.
	re, _ = regexp.Compile("[^\\}\\{/;]+\\{\\}")
	c.CSS = re.ReplaceAll(c.CSS, []byte(""))

	// Add "\" back to fix Opera -o-device-pixel-ratio query
	c.CSS = bytes.Replace(c.CSS, []byte("___YUI_QUERY_FRACTION___"), []byte("/"), -1)
}

// Compress compresses the css in the CSS property and returns the compressed CSS as a byte slice.
func (c *CSSCompressor) Compress() []byte {
	c.extractDataUris()
	c.extractComments()

	// preserve strings so their content doesn't get accidentally minified
	c.extractStrings()

	// strings are safe, now wrestle the comments
	c.processComments()

	// Normalize all whitespace strings to single spaces. Easier to work with that way.
	re, _ := regexp.Compile("\\s+")
	c.CSS = re.ReplaceAll(c.CSS, []byte(" "))

	// Do lots and lots and lots of fun things
	c.performGeneralCleanup()

	// Replace multiple semi-colons in a row with a single one
	re, _ = regexp.Compile(";;+")
	c.CSS = re.ReplaceAll(c.CSS, []byte(";"))

	// restore preserved comments and strings
	for i, t := range c.preservedTokens {
		c.CSS = bytes.Replace(c.CSS, []byte("___YUICSSMIN_PRESERVED_TOKEN_"+strconv.Itoa(i)+"___"), []byte(t), -1)
	}

	// Trim the final string (for any leading or trailing white spaces)
	c.CSS = bytes.TrimSpace(c.CSS)

	// Hooray, we're done!
	return c.CSS
}
