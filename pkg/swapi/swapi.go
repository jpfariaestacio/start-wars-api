package swapi

import (
	"encoding/json"
	"errors"
	fmt "fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/golang/protobuf/jsonpb"
)

const (
	planetsEndPoint = "https://swapi.co/api/planets"
)

var (
	timeout            = time.Duration(5 * time.Second)
	planetMatchPattern = regexp.MustCompile(".+/planets/([0-9]+)")
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

func (ne netTimeoutError) Timeout() bool { return true }

// If you don't set RoundTrip on TimeoutTransport, this will always timeout at 0
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

// Client requests information from Swapi/planets and parses the JSON response to a proto response
func Client(uri string) (interface{}, error) {
	url, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &TimeoutTransport{
			Transport: http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					return net.Dial(netw, addr)
				},
			},
			RoundTripTimeout: time.Second * 10,
		},
	}
	resp, err := client.Get(url.String())
	if err != nil {
		return nil, err
	}
	log.Println(uri)
	defer resp.Body.Close()
	if planetMatchPattern.Match([]byte(uri)) && resp.StatusCode == 404 {
		return nil, fmt.Errorf("the request to %s returned error: NotFound", uri)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Requesting to %s returned %v not 200", uri, resp.StatusCode)
	}
	assertedType, err := assertResponseType(url.String(), json.NewDecoder(resp.Body))
	if err != nil {
		return nil, err
	}
	return assertedType, nil
}

func assertResponseType(uri string, resp *json.Decoder) (interface{}, error) {
	match := planetMatchPattern.Match([]byte(uri))
	if match { // a request to a planet id was used
		sr := SwapiPlanetResponse{}
		if err := jsonpb.UnmarshalNext(resp, &sr); err != nil {
			return nil, err
		}
		id, err := getPlanetID(uri)
		if err != nil {
			return nil, err
		}
		sr.Id = int32(id)
		return sr, nil
	}
	sr := SwapiResponse{}
	if err := jsonpb.UnmarshalNext(resp, &sr); err != nil {
		return nil, err
	}
	return sr, nil
}

func getPlanetID(uri string) (int, error) {
	subgroups := planetMatchPattern.FindStringSubmatch(uri)
	if len(subgroups) <= 1 {
		return 0, fmt.Errorf("the uri %s does not have a planet id", uri)
	}
	id, err := strconv.Atoi(subgroups[1])
	if err != nil {
		return 0, err
	}
	return id, nil
}

// RetrivePlanet requests a planet by its id from Swapi/Planets returning its data
func RetrivePlanet(id int) (SwapiPlanetResponse, error) {
	resp, err := Client(buildPlanetRequest(id))
	if err != nil {
		return resp.(SwapiPlanetResponse), err
	}
	return resp.(SwapiPlanetResponse), nil
}

// RetriveAllPlanets requests all pages from Swapi/planets returning its data
func RetriveAllPlanets(numberOfPages int) []SwapiResponse {
	c := make(chan SwapiResponse)
	srs := []SwapiResponse{}
	for i := 1; i <= numberOfPages; i++ {
		log.Println(i)
		go func(page int) {
			resp, err := Client(buildRequestWithPage(page))
			if err != nil {
				log.Fatal(err)
			}
			log.Println(page)
			c <- resp.(SwapiResponse)
		}(i)
	}
	go func(srs *[]SwapiResponse) {
		for {
			if len(*srs) == numberOfPages {
				close(c)
				return
			}
		}
	}(&srs)
	for resp := range c {
		log.Println(1)
		srs = append(srs, resp)
	}
	return srs
}

func buildRequestWithPage(page int) string {
	return fmt.Sprintf("%s/?page=%s", planetsEndPoint, strconv.Itoa(page))
}

func buildPlanetRequest(id int) string {
	return fmt.Sprintf("%s/%s", planetsEndPoint, strconv.Itoa(id))
}

// GetTotalPages returns the number os pages necessary to get all planets.
func GetTotalPages() (int, error) {
	firstPage, err := Client(planetsEndPoint)
	if err != nil {
		return 0, err
	}
	count := int(firstPage.(SwapiResponse).Count)
	switch {
	case count == 0:
		return 0, errors.New("can't count the number of pages to get all data from swapi/planets as it returns 0")
	case count < 10 && count > 0:
		return 1, nil
	case count%10 == 0:
		return count / 10, nil
	default:
		return (count / 10) + 1, nil
	}
}
