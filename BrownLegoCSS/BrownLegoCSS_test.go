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
			"body { margin: 0px; } /* giggle */ test",
			"body { margin: 0px; } /*___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_0___*/ test",
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
		{"123 /*! Preserve comment */", "123 /*! Preserve comment */"},
		{"123 /* Multiple */ comment /* Comments */", "123  comment "},
		// Not quite sure if the below is correct...
		{"/* Hack comment \\*/ /* asdf */", "/* Hack comment \\*/ "},
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

func Test(t *testing.T) {
	testExtractDataUris(t)
	testExtractComments(t)
	testExtractStrings(t)
	testProcessComments(t)
}
