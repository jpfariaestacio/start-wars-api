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

func GetSwapi(uri string) (SwapiResponse, error) {
	sr := SwapiResponse{}
	resp, err := http.Get(uri)
	if err != nil {
		return sr, err
	}
	defer resp.Body.Close()
	if err := jsonpb.UnmarshalNext(json.NewDecoder(resp.Body), &sr); err != nil {
		return sr, err
	}
	return sr, nil
}

func RetiveAllPlanets(numberOfPages int) []SwapiResponse {
	c := make(chan SwapiResponse)
	srs := []SwapiResponse{}
	for i := 1; i <= numberOfPages; i++ {
		go func(page int) {
			resp, err := GetSwapi(fmt.Sprintf("%s/?page=%s", planetsEndPoint, strconv.Itoa(page)))
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

func GetTotalPages() (int, error) {
	firstPage, err := GetSwapi(planetsEndPoint)
	if err != nil {
		return 0, err
	}
	return int(firstPage.GetCount() / 10), nil
}
