package models

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

type Event struct {
	Uid      int
	Title    string
	StartsAt time.Time
	IsPublic bool
}

func (e Event) GetEtag() string {
	timeStamp := e.StartsAt.Unix()

	s := []byte(strconv.Itoa(e.Uid) + e.Title + strconv.Itoa(int(timeStamp)))

	hasher := sha256.New()
	hasher.Write(s)

	return hex.EncodeToString(hasher.Sum(nil))
}
