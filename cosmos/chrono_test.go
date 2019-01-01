// Copyright 2018 Kuei-chun Chen. All rights reserved.

package cosmos

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/bson"
)

func TestExpand(t *testing.T) {
	var err error
	chaos := Create("testdata/bigbang.json")
	var result = bson.M{}
	if err = chaos.BigBang().Expand().Exec(&result); err != nil {
		t.Fatal(err)
	}
	if result["expand"] == nil || result["expand"] == false {
		t.Fatal("expected true, but got", result["seed"])
	}
}

func TestSeed(t *testing.T) {
	var err error
	chaos := Create("testdata/bigbang.json")
	var result = bson.M{}
	if err = chaos.BigBang().Seed().Exec(&result); err != nil {
		t.Fatal(err)
	}
	if result["seed"] == nil || result["seed"] == false {
		t.Fatal("expected true, but got", result["seed"])
	}
}
