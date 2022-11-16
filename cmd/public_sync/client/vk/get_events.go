package vk

import (
	"encoding/json"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type GetEventsRequest struct {
	RequestURL string
}

func (r GetEventsRequest) URL() string {
	return r.RequestURL
}

func (r GetEventsRequest) Method() string {
	return http.MethodPost
}

func (r GetEventsRequest) Headers() http.Header {
	return map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
	}
}

type GetEventsRequestWithBody struct {
	GetEventsRequest
	FormData map[string]string
}

func (rb GetEventsRequestWithBody) Body() ([]byte, error) {
	form := url.Values{}
	for k, v := range rb.FormData {
		form.Add(k, v)
	}

	return []byte(form.Encode()), nil
}

type VKEvents struct {
	VKEventsResponse VKEventsResponse `json:"response"`
}

type VKEventsResponse struct {
	Count int    `json:"count"`
	Items []Item `json:"items"`
}

type Item struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ScreenName   string `json:"screen_name"`
	IsClosed     int    `json:"is_closed"`
	Type         string `json:"type"`
	IsAdmin      int    `json:"is_admin"`
	IsMember     int    `json:"is_member"`
	IsAdvertiser int    `json:"is_advertiser"`
	Photo50      string `json:"photo_50"`
	Photo100     string `json:"photo_100"`
	Photo200     string `json:"photo_200"`
}

type GetEventsResponse struct {
	VKEvents           VKEvents
	downloadLimitBytes int64
}

func (r *GetEventsResponse) ReadFrom(resp *http.Response) error {
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

	err = json.Unmarshal(data, &r.VKEvents)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal data")
	}

	return nil
}
