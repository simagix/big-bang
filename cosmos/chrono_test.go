// Copyright 2018 Kuei-chun Chen. All rights reserved.

package cosmos

import (
	"testing"

	"github.com/simagix/keyhole/mdb"
	"go.mongodb.org/mongo-driver/bson"
)

func TestExpandString(t *testing.T) {
	var err error
	chaos := Create("testdata/bigbang-string.json")
	var result = bson.M{}
	if err = chaos.BigBang().Expand().Exec(&result); err != nil {
		t.Fatal(err)
	}
	if result["expand"] == nil || result["expand"] == false {
		t.Fatal("expected true, but got", result["seed"])
	}
	t.Log(mdb.Stringify(result), "", "  ")
}

func TestExpandNumber(t *testing.T) {
	var err error
	chaos := Create("testdata/bigbang-number.json")
	var result = bson.M{}
	if err = chaos.BigBang().Expand().Exec(&result); err != nil {
		t.Fatal(err)
	}
	if result["expand"] == nil || result["expand"] == false {
		t.Fatal("expected true, but got", result["seed"])
	}
	t.Log(mdb.Stringify(result), "", "  ")
}

func TestSeedString(t *testing.T) {
	var err error
	chaos := Create("testdata/bigbang-string.json")
	var result = bson.M{}
	if err = chaos.BigBang().Seed().Exec(&result); err != nil {
		t.Fatal(err)
	}
	if result["seed"] == nil || result["seed"] == false {
		t.Fatal("expected true, but got", result["seed"])
	}
}

func TestSeedNumber(t *testing.T) {
	var err error
	chaos := Create("testdata/bigbang-number.json")
	var result = bson.M{}
	if err = chaos.BigBang().Seed().Exec(&result); err != nil {
		t.Fatal(err)
	}
	if result["seed"] == nil || result["seed"] == false {
		t.Fatal("expected true, but got", result["seed"])
	}
}
