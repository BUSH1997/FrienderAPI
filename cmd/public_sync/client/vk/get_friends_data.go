package vk

import (
	"encoding/json"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"io/ioutil"
	"net/http"
)

type UnmarshalResponseFriends struct {
	ResponseFriends `json:"response,omitempty"`
}
type ResponseFriends struct {
	Count int     `json:"count,omitempty"`
	Ids   []int64 `json:"items,omitempty"`
}

type GetFriendsResponse struct {
	VkFriendsData      UnmarshalResponseFriends `json:"vk_friends_data"`
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

	err = json.Unmarshal(data, &r.VkFriendsData)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal data")
	}

	return nil
}
