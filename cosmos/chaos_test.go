// Copyright 2018 Kuei-chun Chen. All rights reserved.

package cosmos

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var sourceURI = "mongodb://localhost/keyhole?replicaSet=replset&authSource=admin"
var sourceDB = "keyhole"

func TestBigBangNumber(t *testing.T) {
	chaos := Create("testdata/bigbang-number.json")
	chaos.SetVerbose(true)
	chrono := chaos.BigBang()
	if chrono.Error() != nil {
		t.Fatal(chrono.Error())
	}
}

func TestBigBangObjectID(t *testing.T) {
	var err error
	var client *mongo.Client
	ctx := context.Background()
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI(sourceURI)); err != nil {
		t.Fatal(err)
	}
	cityID := primitive.NewObjectID()
	cityDoc := bson.D{{Key: "_id", Value: cityID}, {Key: "name", Value: "Atlanta"}}
	userDoc := bson.D{{Key: "name", Value: "simagix"}, {Key: "city_id", Value: cityID}}
	client.Database(sourceDB).Collection("users").InsertOne(ctx, userDoc)
	client.Database(sourceDB).Collection("cities").InsertOne(ctx, cityDoc)

	chrono := Create("testdata/bigbang-oid.json").BigBang()
	if chrono.Error() != nil {
		t.Fatal(chrono.Error())
	}
	client.Database(sourceDB).Collection("users").DeleteOne(ctx, bson.D{{Key: "city_id", Value: cityID}})
	client.Database(sourceDB).Collection("cities").DeleteOne(ctx, bson.D{{Key: "_id", Value: cityID}})
}

func TestBigBangString(t *testing.T) {
	chaos := Create("testdata/bigbang-string.json")
	chaos.SetVerbose(true)
	chrono := chaos.BigBang()
	if chrono.Error() != nil {
		t.Fatal(chrono.Error())
	}
}

func TestGetTemplateFromCollection(t *testing.T) {
	var err error
	var client *mongo.Client
	chaos := Create("testdata/bigbang-string.json")
	chaos.SetVerbose(true)
	ctx := context.Background()
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI(sourceURI)); err != nil {
		t.Fatal(err)
	}

	if _, err = chaos.getTemplateFromCollection(client, "dealers"); err != nil {
		t.Fatal(err)
	}
}

func TestGetFields(t *testing.T) {
	var err error
	var client *mongo.Client
	chaos := Create("testdata/bigbang-string.json")
	chaos.SetVerbose(true)
	ctx := context.Background()
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI(sourceURI)); err != nil {
		t.Fatal(err)
	}
	c := client.Database("keyhole").Collection("dealers")
	var doc bson.M
	var list []interface{}
	if err = c.FindOne(context.Background(), bson.D{{}}).Decode(&doc); err != nil {
		t.Fatal(err)
	}
	list, err = chaos.getFields(doc, "_id", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) > 10 {
		t.Fatal(err)
	}
}
