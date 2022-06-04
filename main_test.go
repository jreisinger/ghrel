package main

import "testing"

func TestFilename(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"https://github.com/jgm/pandoc/releases/download/2.18/pandoc-2.18-1-amd64.deb", "pandoc-2.18-1-amd64.deb"},
		{"", ""},
	}
	for _, test := range tests {
		got, err := fileName(test.url)
		if err != nil {
			t.Error(err)
		}
		if got != test.want {
			t.Errorf("got %s, want %s", got, test.want)
		}
	}
}
