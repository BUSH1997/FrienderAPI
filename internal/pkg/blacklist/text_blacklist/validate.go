package text_blacklist

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/blacklist"
	"github.com/pkg/errors"
	"strings"
)

func (tb TextBlackLister) Validate(data blacklist.RowData) error {
	text := data.CheckData.(string)
	text = strings.ToLower(text)

	for _, regExp := range tb.regExpArray {
		res := regExp.Find([]byte(text))
		if res != nil {
			return errors.New("failed to validate input data")
		}
	}

	return nil
}
