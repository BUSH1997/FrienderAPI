package stammer

import (
	"github.com/goodsign/snowball"
	"github.com/pkg/errors"
	"strings"
)

func GetStammersFromTitle(title string) ([]string, error) {
	stammers, err := GetStammers(strings.Split(title, " "))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stammers from title")
	}

	return stammers, nil
}

func GetStammers(words []string) ([]string, error) {
	stemmer, err := snowball.NewWordStemmer("ru", "UTF_8")
	if err != nil {
		return nil, errors.Wrap(err, "failed to init stammer")
	}

	defer stemmer.Close()

	terms, err := getTerms(words, stemmer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get terms")
	}

	return terms, nil
}

func getTerms(words []string, stemmer *snowball.WordStemmer) ([]string, error) {
	var ret []string
	for _, word := range words {
		term, err := stemmer.Stem([]byte(strings.ToLower(word)))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get term of word %s", word)
		}

		ret = append(ret, string(term))
	}

	return ret, nil
}
