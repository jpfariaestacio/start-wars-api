package swapi

import (
	"encoding/json"
	fmt "fmt"
	"regexp"
	"strconv"

	"github.com/golang/protobuf/jsonpb"
)

var (
	planetsPattern = regexp.MustCompile(`planets/?(([0-9]+)|\?page=[0-9]+)?`)
)

// Convert passes a json from the buffer to the SwapiPlanetResponse
func (spr *SwapiPlanetResponse) Convert(j *json.Decoder) error {
	if err := jsonpb.UnmarshalNext(j, spr); err != nil {
		return err
	}
	intID, err := spr.GetIDFromURL(spr.GetUrl())
	if err != nil {
		return err
	}
	spr.Id = intID
	return nil
}

func (spr *SwapiPlanetResponse) GetIDFromURL(url string) (int32, error) {
	strID := planetsPattern.FindStringSubmatch(url)[1]
	intID, err := strconv.Atoi(strID)
	if err != nil {
		return 0, fmt.Errorf("can not retrive planet ID from URL %s", url)
	}
	return int32(intID), nil
}

// Convert passes a json from the buffer to the SwapiRespose
func (sr *SwapiResponse) Convert(j *json.Decoder) error {
	if err := jsonpb.UnmarshalNext(j, sr); err != nil {
		return err
	}
	return nil
}
