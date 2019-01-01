// Copyright 2019 Kuei-chun Chen. All rights reserved.

package cosmos

import (
	"encoding/json"
	"io/ioutil"
)

// Config -
type Config struct {
	Source      string       `json:"source" bson:"source"`
	Target      string       `json:"target" bson:"target"`
	Collections []Collection `json:"collections" bson:"collections"`
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
