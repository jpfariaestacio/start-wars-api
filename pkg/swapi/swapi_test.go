package swapi

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"
)

// At the time the total number of planets were 61. Change that if necessary
func TestRetriveAllPlanets(t *testing.T) {
	totalPages, err := GetTotalPages()
	if err != nil {
		t.Error(err)
	}
	if totalPages != 7 {
		t.Errorf("the number of pages should be 7 not %v", totalPages)
	}
	resp := RetriveAllPlanets(totalPages)
	if len(resp) != 7 {
		t.Errorf("the len of the object response should be 7 not %v", len(resp))
	}
	cntr := 0
	for _, objs := range resp {
		cntr += len(objs.Results)
	}
	if cntr != 61 {
		t.Errorf("the number of planets should be 61 not %v", cntr)
	}
}

var assertResponseTypeTable = []struct {
	url          string
	jsonResponse string
	out          string
}{
	{"https://swapi.co/api/planets/61", `{
		"name": "Jakku",
		"rotation_period": "unknown",
		"orbital_period": "unknown",
		"diameter": "unknown",
		"climate": "unknown",
		"gravity": "unknown",
		"terrain": "deserts",
		"surface_water": "unknown",
		"population": "unknown",
		"residents": [],
		"films": [
			"https://swapi.co/api/films/7/"
		],
		"created": "2015-04-17T06:55:57.556495Z",
		"edited": "2015-04-17T06:55:57.556551Z",
		"url": "https://swapi.co/api/planets/61/"
	}`, "swapiPlanetResponse"},
	{"https://swapi.co/api/planets/", `{
		"count": 61,
		"next": "https://swapi.co/api/planets/?page=2",
		"previous": null,
		"results": []
		},`, "swapiResponse"},
}

func TestAssertResponseType(t *testing.T) {
	for _, tt := range assertResponseTypeTable {
		reader := strings.NewReader(tt.jsonResponse)
		resp, err := assertResponseType(tt.url, json.NewDecoder(reader))
		if err != nil {
			t.Error(err)
		}
		switch tt.out {
		case "swapiResponse":
			if _, ok := resp.(SwapiResponse); !ok {
				t.Errorf("expected request to %s to be converted to type SwapiResponse not %s", tt.url, reflect.TypeOf(resp))
			}
		case "swapiPlanetResponse":
			if _, ok := resp.(SwapiPlanetResponse); !ok {
				t.Errorf("expected request to %s to be converted to type SwapiPlanetResponse not %s", tt.url, reflect.TypeOf(resp))
			}
		}
	}
}

var matchPatternTable = []struct {
	in       string
	match    bool
	subgroup string
}{
	{"https://swapi.co/api/planets", false, ""},
	{"https://swapi.co/api/planets/", false, ""},
	{"https://swapi.co/api/planets/1", true, "1"},
	{"https://swapi.co/api/planets/111111", true, "111111"},
	{"https://swapi.co/api/planets/?page=1", false, ""},
	{"https://swapi.co/api/planets/11/12/13/14", true, "11"},
}

func TestMatchPlanetRequestPatternAndID(t *testing.T) {
	var planetMatchPattern = regexp.MustCompile(".+/planets/([0-9]+)")
	for _, tt := range matchPatternTable {
		subgroups := planetMatchPattern.FindStringSubmatch(tt.in)
		switch {
		case len(subgroups) <= 1 && tt.match:
			t.Errorf("%s should be a match and return the subgroup: %s", tt.in, tt.subgroup)
		case len(subgroups) > 1 && !tt.match:
			t.Errorf("%s shouldn't be a match but returns the subgroup: %s", tt.in, subgroups[1])
		case len(subgroups) > 1 && tt.match:
			if subgroups[1] != tt.subgroup {
				t.Errorf("got %s, want %s", subgroups[1], tt.subgroup)
			}
		}
	}
}

var getTotalPagesTable = []struct {
	name string
	in   int
	out  int
}{
	{"One", 1, 1},
	{"LessThanTen", 7, 1},
	{"Zero", 0, 0},
	{"not%10==0", 61, 7},
	{"%10==0", 70, 7},
}

func TestGetTotalPagesLogic(t *testing.T) {
	for _, tt := range getTotalPagesTable {
		switch {
		case tt.in == 0:
			r := 0
			if r != tt.out {
				t.Errorf("got %v, want %v", r, tt.out)
			}
		case tt.in < 10 && tt.in > 0:
			r := 1
			if r != tt.out {
				t.Errorf("got %v, want %v", r, tt.out)
			}
		case tt.in%10 == 0:
			r := tt.in / 10
			if r != tt.out {
				t.Errorf("got %v, want %v", r, tt.out)
			}
		default:
			r := tt.in / 10
			if r+1 != tt.out {
				t.Errorf("got %v, want %v", r, tt.out)
			}
		}
	}
}

func TestGetTotalPagesRequest(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := GetTotalPages()
			if err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
}

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
		r := buildRequestWithPage(tt.in)
		if r != tt.out {
			t.Errorf("got %q, want %q", r, tt.out)
		}
	}
}

var buildPlanetRequestTable = []struct {
	in  int
	out string
}{
	{1, "https://swapi.co/api/planets/1"},
	{2, "https://swapi.co/api/planets/2"},
	{0, "https://swapi.co/api/planets/0"},
	{1000000, "https://swapi.co/api/planets/1000000"},
	{465165, "https://swapi.co/api/planets/465165"},
}

func TestBuildPlanetRequest(t *testing.T) {
	for _, tt := range buildPlanetRequestTable {
		r := buildPlanetRequest(tt.in)
		if r != tt.out {
			t.Errorf("got %q, want %q", r, tt.out)
		}
	}
}

func TestHttpTimeout(t *testing.T) {
	http.HandleFunc("/normal", func(w http.ResponseWriter, req *http.Request) {
		// Empirically, timeouts less than these seem to be flaky
		time.Sleep(100 * time.Millisecond)
		io.WriteString(w, "ok")
	})
	http.HandleFunc("/timeout", func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(250 * time.Millisecond)
		io.WriteString(w, "ok")
	})
	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()

	numDials := 0

	client := &http.Client{
		Transport: &TimeoutTransport{
			Transport: http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					t.Logf("dial to %s://%s", netw, addr)
					numDials++                  // For testing only.
					return net.Dial(netw, addr) // Regular ass dial.
				},
			},
			RoundTripTimeout: time.Millisecond * 200,
		},
	}

	addr := ts.URL

	SendTestRequest(t, client, "1st", addr, "normal")
	if numDials != 1 {
		t.Fatalf("Should only have 1 dial at this point.")
	}
	SendTestRequest(t, client, "2st", addr, "normal")
	if numDials != 1 {
		t.Fatalf("Should only have 1 dial at this point.")
	}
	SendTestRequest(t, client, "3st", addr, "timeout")
	if numDials != 1 {
		t.Fatalf("Should only have 1 dial at this point.")
	}
	SendTestRequest(t, client, "4st", addr, "normal")
	if numDials != 2 {
		t.Fatalf("Should have our 2nd dial.")
	}

	time.Sleep(time.Millisecond * 700)

	SendTestRequest(t, client, "5st", addr, "normal")
	if numDials != 2 {
		t.Fatalf("Should still only have 2 dials.")
	}
}

func SendTestRequest(t *testing.T, client *http.Client, id, addr, path string) {
	req, err := http.NewRequest("GET", addr+"/"+path, nil)

	if err != nil {
		t.Fatalf("new request failed - %s", err)
	}

	req.Header.Add("Connection", "keep-alive")

	switch path {
	case "normal":
		if resp, err := client.Do(req); err != nil {
			t.Fatalf("%s request failed - %s", id, err)
		} else {
			result, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				t.Fatalf("%s response read failed - %s", id, err2)
			}
			resp.Body.Close()
			t.Logf("%s request - %s", id, result)
		}
	case "timeout":
		if _, err := client.Do(req); err == nil {
			t.Fatalf("%s request not timeout", id)
		} else {
			t.Logf("%s request - %s", id, err)
		}
	}
}
