// Copyright 2019 Kuei-chun Chen. All rights reserved.

package cosmos

import (
	"encoding/json"
	"io/ioutil"
)

// Config -
type Config struct {
	Source      MongoConfig  `json:"source" bson:"source"`
	Target      MongoConfig  `json:"target" bson:"target"`
	Collections []Collection `json:"collections" bson:"collections"`
}

// MongoConfig stores mongo configuration
type MongoConfig struct {
	URI       string `json:"uri" bson:"uri"`
	CAFile    string `json:"caFile" bson:"caFile"`
	ClientPEM string `json:"clientPEM" bson:"clientPEM"`
}

// Collection -
type Collection struct {
	Lookup []Lookup `json:"lookup" bson:"lookup"`
	Name   string   `json:"name" bson:"name"`
	Total  int64    `json:"total" bson:"total"`
}

// Lookup -
type Lookup struct {
	ForeignField string `json:"foreignField" bson:"foreignField"`
	From         string `json:"from" bson:"from"`
	LocalField   string `json:"localField" bson:"localField"`
	Total        int64  `json:"total" bson:"total"`
	NumSeeds     int    `json:"nSeeds" bson:"nSeeds"`
}

// Create parses configuration string
func Create(filename string) *Chaos {
	var err error
	var b []byte
	if b, err = ioutil.ReadFile(filename); err != nil {
		return &Chaos{err: err}
	}
	chaos := Chaos{filename: filename}
	chaos.err = json.Unmarshal(b, &chaos.config)
	return &chaos
}
