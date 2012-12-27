package BrownLegoCSS

import (
	"fmt"
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
			"body { margin: 0px; } /* giggle */",
			"body { margin: 0px; } /*___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_0___*/",
			[]string{" giggle "},
		},
		{
			"/* This is a comment yo */ body { margin: 0px; } /* giggle */",
			"/*___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_0___*/ body { margin: 0px; } /*___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_1___*/",
			[]string{" This is a comment yo ", " giggle "},
		},
	}

	for i, c := range tests {
		compressor := CssCompressor{}
		compressor.Css = []byte(tests[i].s)
		compressor.extractComments()
		got := string(compressor.Css)
		if got != c.want {
			t.Errorf("css: testExtractComments(%q) == (%q), want %q", c.s, got, c.want)
		}
		if fmt.Sprintf("%v", compressor.comments) != fmt.Sprintf("%v", c.comments) {
			t.Errorf("comments: testExtractComments(%q) == (%q), want %q", c.s, compressor.comments, c.comments)
		}
	}
}

func Test(t *testing.T) {
	testExtractDataUris(t)
	testExtractComments(t)
}
