// Copyright 2019 Kuei-chun Chen. All rights reserved.

package cosmos

import (
	"encoding/json"
	"io/ioutil"
)

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
