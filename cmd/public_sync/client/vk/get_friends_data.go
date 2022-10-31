package vk

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type ResponseFriends struct {
	Count int
	Ids   [][]int
}

type GetFriendsResponse struct {
	VkFriendsData      ResponseFriends
	DownloadLimitBytes int64
}

func (r *GetFriendsResponse) ReadFrom(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(
			errors.New("got not 200 response code"),
			"got bad response code %d",
			resp.StatusCode,
		)
	}

	reader := http.MaxBytesReader(nil, resp.Body, int64(r.DownloadLimitBytes))
	if r.DownloadLimitBytes == 0 {
		reader = resp.Body
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.Wrap(err, "failed to read data from reader")
	}

	fmt.Println(string(data))

	err = json.Unmarshal(data, &r.VkFriendsData)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal data")
	}

	return nil
}
