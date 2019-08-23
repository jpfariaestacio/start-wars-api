package swapi

import (
	"encoding/json"
	"errors"
	fmt "fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/golang/protobuf/jsonpb"
)

var (
	planetsPattern = regexp.MustCompile(`planets/?(([0-9]+)|\?page=[0-9]+)?`)
)

func (spr *SwapiPlanetResponse) Convert(j *json.Decoder) error {
	if err := jsonpb.UnmarshalNext(j, spr); err != nil {
		return err
	}
	return nil
}

func (sr *SwapiResponse) Convert(j *json.Decoder) error {
	if err := jsonpb.UnmarshalNext(j, sr); err != nil {
		return err
	}
	return nil
}

func getRequestedID(uri string) (int, error) {
	subgroups := planetsPattern.FindStringSubmatch(uri)
	if len(subgroups) <= 1 {
		return 0, fmt.Errorf("the uri %s does not have a planet id", uri)
	}
	id, err := strconv.Atoi(subgroups[2])
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
