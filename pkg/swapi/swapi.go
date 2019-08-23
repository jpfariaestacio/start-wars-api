package swapi

import (
	"encoding/json"
	"errors"
	fmt "fmt"
	"net"
	"net/http"
	"time"
)

const (
	planetsEndPoint = "https://swapi.co/api/"
)

var (
	timeout = time.Duration(5 * time.Second)
)

type TimeoutTransport struct {
	http.Transport
	RoundTripTimeout time.Duration
}

type respAndErr struct {
	resp *http.Response
	err  error
}

type netTimeoutError struct {
	error
}

// RoundTrip handles the timeout while requesting to swapi
func (t *TimeoutTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	timeout := time.After(t.RoundTripTimeout)
	resp := make(chan respAndErr, 1)
	go func() {
		r, e := t.Transport.RoundTrip(req)
		resp <- respAndErr{
			resp: r,
			err:  e,
		}
	}()

	select {
	case <-timeout: // A round trip timeout has occurred.
		t.Transport.CancelRequest(req)
		return nil, netTimeoutError{
			error: fmt.Errorf("timed out after %s", t.RoundTripTimeout),
		}
	case r := <-resp: // Success!
		return r.resp, r.err
	}
}

type SwapiClient struct {
	client *http.Client
}

func NewSwapiClient() *SwapiClient {
	sc := new(SwapiClient)
	sc.client = &http.Client{
		Transport: &TimeoutTransport{
			Transport: http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					return net.Dial(netw, addr)
				},
			},
			RoundTripTimeout: timeout,
		},
	}
	return sc
}

// Request requests information from Swapi/planets and parses the JSON response from the buffer
func (sc *SwapiClient) Request(url string) (*json.Decoder, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Connection", "keep-alive")
	resp, err := sc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("the request to %s returned error: NotFound", url)
	}
	return json.NewDecoder(resp.Body), nil
}

type SwapiProtoObj interface {
	Convert(j *json.Decoder) error
}

// NewSwapiProtoObj takes the JSON response from swapi server and converts it to the correct protobuf type
func NewSwapiProtoObj(j *json.Decoder) (SwapiProtoObj, error) {
	var JSON map[string]interface{}
	switch {
	case planetsPattern.Match(JSON["url"].([]byte)):
		if _, ok := JSON["name"]; ok {
			spr := new(SwapiPlanetResponse)
			if err := spr.Convert(j); err != nil {
				return nil, err
			}
			return spr, nil
		} else if _, ok := JSON["count"]; ok {

			sr := new(SwapiResponse)
			if err := sr.Convert(j); err != nil {
				return nil, err
			}
			return sr, nil
		} else {
			return nil, errors.New("the json data does not match any of the expected returns to the planet API")
		}
	}
	return nil, errors.New("the JSON response from https://swapi.co/ does not match any expected schema")
}
