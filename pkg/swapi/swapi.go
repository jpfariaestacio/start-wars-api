package swapi

import (
	"encoding/json"
	fmt "fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/golang/protobuf/jsonpb"
)

const (
	planetsEndPoint = "https://swapi.co/api/planets"
)

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
	query := url.Query()
	switch query.Get("page") {
	case "":
		sr := SwapiPlanetResponse{}
		if err := jsonpb.UnmarshalNext(json.NewDecoder(resp.Body), &SwapiPlanetResponse{}); err != nil {
			return sr, err
		}
	default:
		sr := SwapiResponse{}
		if err := jsonpb.UnmarshalNext(json.NewDecoder(resp.Body), &SwapiResponse{}); err != nil {
			return sr, err
		}
	}
	return nil, nil
}

func GetPlanet(id int) (SwapiPlanetResponse, error) {
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
			resp, err := Client(buildPage(page))
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

func buildPage(page int) string {
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
