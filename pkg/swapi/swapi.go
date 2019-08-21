package swapi

import (
	"encoding/json"
	fmt "fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/golang/protobuf/jsonpb"
)

const (
	planetsEndPoint = "https://swapi.co/api/planets"
)

var planetMatchPattern = regexp.MustCompile(".+/planets/([0-9]+)")

// Client requests information from Swapi/planets and parses the JSON response to a proto response
func Client(uri string) (interface{}, error) {
	url, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Requesting to %s returned %v not 200", uri, resp.StatusCode)
	}
	defer resp.Body.Close()
	assertedType, err := assertResponseType(url, json.NewDecoder(resp.Body))
	if err != nil {
		return nil, err
	}
	return assertedType, nil
}

func assertResponseType(uri *url.URL, resp *json.Decoder) (interface{}, error) {
	match := planetMatchPattern.Match([]byte(uri.String()))
	if match { // a request to a planet id was used
		sr := SwapiPlanetResponse{}
		if err := jsonpb.UnmarshalNext(resp, &sr); err != nil {
			return nil, err
		}
		id, err := getPlanetId(uri.String())
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

func getPlanetId(uri string) (int, error) {
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
		go func(page int) {
			resp, err := Client(buildRequestWithPage(page))
			if err != nil {
				log.Fatal(err)
			}
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
	return int(firstPage.(SwapiResponse).Count / 10), nil
}
