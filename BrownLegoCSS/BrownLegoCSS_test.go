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

func Test(t *testing.T) {
	testExtractDataUris(t)
}
