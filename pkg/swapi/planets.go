package swapi

import (
	"encoding/json"
	"regexp"

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
	return nil
}

// Convert passes a json from the buffer to the SwapiRespose
func (sr *SwapiResponse) Convert(j *json.Decoder) error {
	if err := jsonpb.UnmarshalNext(j, sr); err != nil {
		return err
	}
	return nil
}
