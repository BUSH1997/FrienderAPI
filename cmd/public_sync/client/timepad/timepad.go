package timepad

import (
	"context"
	"encoding/json"
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/client"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	httplib "github.com/BUSH1997/FrienderAPI/internal/pkg/tools/http"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type TimePadTransportConfig struct {
	DownloadLimitBytes int64 `mapstructure:"download_limit_bytes"`
}

type HTTPTimePadClient struct {
	config TimePadTransportConfig
	client httplib.Client
}

func New(client httplib.Client) client.PublicEventsClient {
	return &HTTPTimePadClient{
		client: client,
	}
}

func (c HTTPTimePadClient) UploadPublicEvents(ctx context.Context, url string) ([]models.Event, error) {
	resp := Response{
		downloadLimitBytes: c.config.DownloadLimitBytes,
	}
	err := c.client.PerformRequest(ctx, Request{
		RequestURL: url,
	}, &resp)

	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request to timepad api")
	}

	return convertEventsToModes(resp.PublicEventsData), nil
}

func convertEventsToModes(data PublicEventsData) []models.Event {
	events := make([]models.Event, 0, len(data.Values))
	for _, value := range data.Values {
		event := models.Event{
			Uid:      strconv.Itoa(value.ID),
			Title:    value.Name,
			StartsAt: value.StartsAt.Time.Unix(),
			IsPublic: true,
		}

		events = append(events, event)
	}

	return events
}

type Request struct {
	RequestURL string
}

func (r Request) URL() string {
	return r.RequestURL
}

func (r Request) Method() string {
	return http.MethodGet
}

func (r Request) Headers() http.Header {
	// TODO: token in config(or secret)
	return map[string][]string{
		"Authorization": {"Bearer a9231732800b6f6e059b492a9ce8f96f939f2dee"},
	}
}

const ctLayout = "2006-01-02T15:04:05"

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) parse(s string) {
	signPosition := 19
	offsetStartPosition := 20

	offset, err := strconv.ParseInt(s[offsetStartPosition:offsetStartPosition+2], 10, 64)
	if err != nil {
		ct.Time = time.Time{}
	}

	ct.Time, err = time.Parse(ctLayout, s[:signPosition])
	if err != nil {
		ct.Time = time.Time{}
	}

	if s[signPosition] == '-' {
		ct.Time = ct.Time.Add(time.Duration(offset * 60 * 60 * 1000000000))
	}
	if s[signPosition] == '+' {
		ct.Time = ct.Time.Add(-time.Duration(offset * 60 * 60 * 1000000000))
	}
}

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return errors.Wrapf(err, "failed to unquote string %s", s)
	}

	ct.parse(s)

	return
}

type Value struct {
	ID          int         `json:"id"`
	StartsAt    CustomTime  `json:"starts_at"`
	Name        string      `json:"name"`
	URL         string      `json:"url"`
	PosterImage PosterImage `json:"poster_image"`
	Location    Location    `json:"location"`
	Categories  []Category  `json:"categories"`
}

type PosterImage struct {
	DefaultURL    string `json:"default_url"`
	UploadCareURL string `json:"uploadcare_url"`
}

type Coordinates []string

type Location struct {
	Country     string      `json:"country"`
	City        string      `json:"city"`
	Address     string      `json:"address"`
	Coordinates Coordinates `json:"coordinates"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PublicEventsData struct {
	Values []Value `json:"values"`
}

type Response struct {
	PublicEventsData   PublicEventsData
	downloadLimitBytes int64
}

func (r *Response) ReadFrom(resp *http.Response) error {
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

	err = json.Unmarshal(data, &r.PublicEventsData)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal data")
	}

	return nil
}
