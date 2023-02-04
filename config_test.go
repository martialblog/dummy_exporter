package main

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	var config Config

	err := config.LoadConfig("testdata/dummy-good.json")

	if err != nil {
		t.Errorf("Error loading config %v: %v", "dummy-good.json", err)
	}
}

func TestLoadBadConfigs(t *testing.T) {
	var config Config
	testcases := []struct {
		file string
		want string
	}{
		{
			file: "testdata/dummy-invalid.json",
			want: "Error during Unmarshal: invalid character ':' after array element",
		},
		{
			file: "testdata/does-not-exists.json",
			want: "Error while opening configuration: open testdata/does-not-exists.json: no such file or directory",
		},
	}

	for _, test := range testcases {
		t.Run(test.file, func(t *testing.T) {
			got := config.LoadConfig(test.file)
			if got == nil || got.Error() != test.want {
				t.Fatalf("LoadConfig(%q) = %v; expected %q", test.file, got, test.want)
			}
		})
	}
}
