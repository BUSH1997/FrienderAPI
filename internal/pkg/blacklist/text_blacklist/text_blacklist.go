package text_blacklist

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/blacklist"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"regexp"
)

type TextBlackLister struct {
	regExpArray []regexp.Regexp
}

type TextBlackListerInput struct {
	data string
}

func New(filters []string) (blacklist.BlackLister, error) {
	regExpArray := make([]regexp.Regexp, 0, len(filters))
	for _, filter := range filters {
		regExp, err := regexp.Compile(filter)
		if err != nil {
			return nil, errors.Wrap(err, "failed to compile regular expression")
		}

		regExpArray = append(regExpArray, *regExp)
	}
	return TextBlackLister{
		regExpArray: regExpArray,
	}, nil
}

func (d TextBlackListerInput) Type() string {
	return "text"
}
