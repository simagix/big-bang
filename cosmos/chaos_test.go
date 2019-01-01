// Copyright 2018 Kuei-chun Chen. All rights reserved.

package cosmos

import (
	"context"
	"testing"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func TestBigBangString(t *testing.T) {
	chaos := Create("testdata/bigbang-string.json")
	chaos.SetVerbose(true)
	chrono := chaos.BigBang()
	if chrono.Error() != nil {
		t.Fatal(chrono.Error())
	}
}

func TestBigBangNumber(t *testing.T) {
	chaos := Create("testdata/bigbang-number.json")
	chaos.SetVerbose(true)
	chrono := chaos.BigBang()
	if chrono.Error() != nil {
		t.Fatal(chrono.Error())
	}
}

func TestGetFields(t *testing.T) {
	var err error
	var client *mongo.Client
	chaos := Create("testdata/bigbang-string.json")
	chaos.SetVerbose(true)
	ctx := context.Background()
	if client, err = mongo.Connect(ctx, "mongodb://localhost/keyhole?replicaSet=replset"); err != nil {
		t.Fatal(err)
	}
	c := client.Database("keyhole").Collection("dealers")
	var doc bson.M
	var list []interface{}
	if err = c.FindOne(context.Background(), nil).Decode(&doc); err != nil {
		t.Fatal(err)
	}
	list, err = chaos.getFields(doc, "_id", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 10 {
		t.Fatal(err)
	}
}
