package vk

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type GetEventsDataRequest struct {
	RequestURL string
}

func (r GetEventsDataRequest) URL() string {
	return r.RequestURL
}

func (r GetEventsDataRequest) Method() string {
	return http.MethodPost
}

func (r GetEventsDataRequest) Headers() http.Header {
	return map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
	}
}

type GetEventsDataRequestWithBody struct {
	GetEventsDataRequest
	FormData map[string]string
}

func (rb GetEventsDataRequestWithBody) Body() ([]byte, error) {
	form := url.Values{}
	for k, v := range rb.FormData {
		form.Add(k, v)
	}

	return []byte(form.Encode()), nil
}

type GetEventsDataResponse struct {
	VKEventsData       VKEventsData
	downloadLimitBytes int64
}

type VKEventsData struct {
	VKEventsData       []VKEventData `json:"response,omitempty"`
	downloadLimitBytes int64
}

type VKEventData struct {
	ID           int64     `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	ScreenName   string    `json:"screen_name,omitempty"`
	Description  string    `json:"description,omitempty"`
	IsClosed     int       `json:"is_closed,omitempty"`
	Addresses    Addresses `json:"addresses"`
	Type         string    `json:"type,omitempty"`
	IsAdmin      int       `json:"is_admin,omitempty"`
	IsMember     int       `json:"is_member,omitempty"`
	IsAdvertiser int       `json:"is_advertiser,omitempty"`
	Photo50      string    `json:"photo_50,omitempty"`
	Photo100     string    `json:"photo_100,omitempty"`
	Photo200     string    `json:"photo_200,omitempty"`
	StartDate    int64     `json:"start_date,omitempty"`
	FinishDate   int64     `json:"finish_date,omitempty"`
	Category     string    `json:"-"`
	Place        Place     `json:"place,omitempty"`
}

type Place struct {
	Id        int     `json:"id,omitempty"`
	Title     string  `json:"title,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Type      string  `json:"type,omitempty"`
	Country   int     `json:"country,omitempty"`
	City      int     `json:"city,omitempty"`
	Addresses string  `json:"addresses,omitempty"`
}

type Addresses struct {
	MainAddress MainAddress `json:"main_address,omitempty"`
}

type MainAddress struct {
	Address string  `json:"address,omitempty"`
	City    City    `json:"city,omitempty"`
	Country Country `json:"country,omitempty"`
	Title   string  `json:"title"`
}

type City struct {
	Title string `json:"title"`
}

type Country struct {
	Title string `json:"title"`
}

const SourceTypeEventVK = "vk_event"

type EventInfo struct {
	Category string
	Source   string
}

func (r *GetEventsDataResponse) ReadFrom(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(
			errors.New("got not 200 response code"),
			"got bad response code %d",
			resp.StatusCode,
		)
	}

	reader := http.MaxBytesReader(nil, resp.Body, int64(r.downloadLimitBytes))
	if r.downloadLimitBytes == 0 {
		reader = resp.Body
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.Wrap(err, "failed to read data from reader")
	}

	fmt.Println(string(data))

	err = json.Unmarshal(data, &r.VKEventsData)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal data")
	}

	return nil
}
