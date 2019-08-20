package swapi

import (
	"encoding/json"
	fmt "fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/golang/protobuf/jsonpb"
)

const (
	planetsEndPoint = "https://swapi.co/api/planets"
)

// Client requests information from Swapi/planets and parses the JSON response to a proto response
func Client(uri string) (SwapiResponse, error) {
	sr := SwapiResponse{}
	resp, err := http.Get(uri)
	if err != nil {
		return sr, err
	} else if resp.StatusCode != 200 {
		return sr, fmt.Errorf("Requesting to %s returned %v not 200", uri, resp.StatusCode)
	}
	defer resp.Body.Close()
	if err := jsonpb.UnmarshalNext(json.NewDecoder(resp.Body), &sr); err != nil {
		return sr, err
	}
	return sr, nil
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
			c <- resp
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

// GetTotalPages returns the number os pages necessary to get all planets.
func GetTotalPages() (int, error) {
	firstPage, err := Client(planetsEndPoint)
	if err != nil {
		return 0, err
	}
	return int(firstPage.GetCount() / 10), nil
}
