package swapi

import (
	"bytes"
	"encoding/json"
	"errors"
	fmt "fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

const (
	planetsEndPoint = "https://swapi.co/api/"
)

var (
	timeout = time.Duration(20 * time.Second)
)

// TimeoutTransport sets the values for timeout
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

// Client sets the client (as timeout settings) and request type (as planets or actos)
type Client struct {
	requestType string
	client      *http.Client
}

// NewSwapiClient is a factory to the Client struct
func NewSwapiClient(requestType string) *Client {
	sc := new(Client)
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
	sc.requestType = requestType
	return sc
}

// Request requests information from Swapi/planets and parses the JSON response from the buffer
func (sc *Client) Request(url string) (ProtoObj, error) {
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
	spo, err := sc.newProtoObj(resp.Body)
	if err != nil {
		return nil, err
	}
	return spo, nil
}

// ProtoObj abstracts the conversion between response proto types
type ProtoObj interface {
	Convert(j *json.Decoder) error
}

// NewProtoObj takes the JSON response from swapi server and converts it to the correct protobuf type
func (sc *Client) newProtoObj(body io.ReadCloser) (ProtoObj, error) {
	b, fields, err := getFields(body)
	if err != nil {
		return nil, err
	}
	body = ioutil.NopCloser(bytes.NewBuffer(b))
	decoded := json.NewDecoder(body)
	switch sc.requestType {
	case "planets":
		if fieldExist("name", fields) {
			spr := new(SwapiPlanetResponse)
			if err := spr.Convert(decoded); err != nil {
				return nil, err
			}
			return spr, nil
		} else if fieldExist("count", fields) {
			sr := new(SwapiResponse)
			if err := sr.Convert(decoded); err != nil {
				return nil, err
			}
			return sr, nil
		} else {
			return nil, errors.New("the json data does not match any of the expected returns to the planet API")
		}
	}
	return nil, errors.New("the JSON response from https://swapi.co/ does not match any expected schema")
}

func fieldExist(field string, fields []string) bool {
	for _, k := range fields {
		if field == k {
			return true
		}
	}
	return false
}

func getFields(body io.ReadCloser) ([]byte, []string, error) {
	var JSON map[string]interface{}
	var fields []string
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, nil, err
	}
	json.Unmarshal(b, &JSON)
	for k := range JSON {
		fields = append(fields, k)
	}
	return b, fields, nil
}

// RetriveItem gets an item especified by ID from Swapi
func (sc *Client) RetriveItem(id int) (ProtoObj, error) {
	po, err := sc.Request(sc.buildItemRequestURL(id))
	if err != nil {
		return nil, err
	}
	return po, nil
}

// RetriveAll takes the given type and retrive all objects from swapi
func (sc *Client) RetriveAll() ([]ProtoObj, error) {
	var spo []ProtoObj
	numberOfPages := 1
	nOfObjsChan := make(chan int, 1)
	errChan := make(chan error)
	objChan := make(chan ProtoObj)
	go func() {
		for i := 1; i <= numberOfPages; i++ {
			go func(page int) {
				obj, err := sc.Request(sc.buildItemsRequestURL(page))
				if err != nil {
					errChan <- err
					return
				}
				if numberOfPages == 1 {
					nOfObjsChan <- getCount(obj)
				}
				objChan <- obj
			}(i)
			if numberOfPages == 1 {
				numberOfPages = getTotalPages(<-nOfObjsChan)
			}
		}
	}()
	for {
		select {
		case obj := <-objChan:
			spo = append(spo, obj)
		case err := <-errChan:
			return nil, err
		default:
			if len(spo) == numberOfPages {
				return spo, nil
			}
		}
	}
}

func getTotalPages(count int) int {
	switch {
	case count == 0:
		return 0
	case count < 10 && count > 0:
		return 1
	case count%10 == 0:
		return count / 10
	default:
		return (count / 10) + 1
	}
}

func getCount(obj ProtoObj) int {
	countFunc := reflect.ValueOf(obj).MethodByName("GetCount")
	in := make([]reflect.Value, countFunc.Type().NumIn())
	if countFunc.IsValid() {
		count := countFunc.Call(in)[0]
		return int(count.Int())
	}
	return 0
}

func (sc *Client) buildItemRequestURL(id int) string {
	return fmt.Sprintf("%s/%s/%s", planetsEndPoint, sc.requestType, strconv.Itoa(id))
}

func (sc *Client) buildItemsRequestURL(page int) string {
	return fmt.Sprintf("%s%s/?page=%s", planetsEndPoint, sc.requestType, strconv.Itoa(page))
}
