package swapi

import "testing"

var buildRequestWithPageTable = []struct {
	in  int
	out string
}{
	{1, "https://swapi.co/api/planets/?page=1"},
	{2, "https://swapi.co/api/planets/?page=2"},
	{0, "https://swapi.co/api/planets/?page=0"},
	{1000000, "https://swapi.co/api/planets/?page=1000000"},
	{465165, "https://swapi.co/api/planets/?page=465165"},
}

func TestBuildRequestWithPage(t *testing.T) {
	for _, tt := range buildRequestWithPageTable {
		t.Run(tt.in, func(t *testing.T) {
			r := buildPlanetRequest(tt.in)
			if r != tt.out {
				t.Errorf("got %q, want %q", s, tt.out)
			}
		})
	}
}
