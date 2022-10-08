package models

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

type Event struct {
	Uid         string
	Title       string
	Description string
	Members     []int
	Images      []string
	TimeCreated time.Time
	TimeUpdated time.Time
	GeoData     Geo
	Author      int
	StartsAt    time.Time
	IsGroup     bool
	IsPublic    bool
	Category    Category
}

type Category string

const (
	Sport Category = "SPORT"
)

type Geo struct {
	Longitude float64
	Latitude  float64
}

func (e Event) GetEtag() string {
	timeStamp := e.StartsAt.Unix()

	s := []byte(e.Uid + e.Title + strconv.Itoa(int(timeStamp)))

	hasher := sha256.New()
	hasher.Write(s)

	return hex.EncodeToString(hasher.Sum(nil))
}
