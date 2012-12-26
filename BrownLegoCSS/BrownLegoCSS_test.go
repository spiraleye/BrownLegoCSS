package BrownLegoCSS

import "testing"

func testExtractDataUris(t *testing.T) {
	compressor := CssCompressor{}

	var tests = []struct {
		s, want string
	}{
		{"background-image: url('elephant.png');", "background-image: url('elephant.png');"},
		{"background-image: url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQAQMAAAAlPW0iAAAABlBMVEUAAAD///+l2Z/dAAAAM0lEQVR4nGP4/5/h/1+G/58ZDrAz3D/McH8yw83NDDeNGe4Ug9C9zwz3gVLMDA/A6P9/AFGGFyjOXZtQAAAAAElFTkSuQmCC');", "background-image: url(___YUICSSMIN_PRESERVED_TOKEN_0___);"},
	}

	for i, c := range tests {
		compressor.Css = []byte(tests[i].s)
		compressor.extractDataUris()
		got := string(compressor.Css)
		if got != c.want {
			t.Errorf("extractDataUris(%q) == (%q), want %q", c.s, got, c.want)
		}
	}
}

func testExtractComments(t *testing.T) {
	compressor := CssCompressor{}

	var tests = []struct {
		s, want string
	}{
		{"No comments here", "No comments here"},
		{"body { margin: 0px; } /* giggle */", "body { margin: 0px; } /*___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_0___*/"},
		{"/* This is a comment yo */ body { margin: 0px; } /* giggle */", "/*___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_0___*/ body { margin: 0px; } /*___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_1___*/"},
	}

	for i, c := range tests {
		compressor.Css = []byte(tests[i].s)
		compressor.extractComments()
		got := string(compressor.Css)
		if got != c.want {
			t.Errorf("testExtractComments(%q) == (%q), want %q", c.s, got, c.want)
		}
	}
}

func Test(t *testing.T) {
	testExtractDataUris(t)
	testExtractComments(t)
}
