// Copyright 2019 Kuei-chun Chen. All rights reserved.

package cosmos

import (
	"context"
	"encoding/json"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/network/connstring"
	"github.com/simagix/keyhole/sim/util"
)

// Mongo -
type Mongo struct {
	uri string
}

// Source -
func Source(uri string) *Mongo {
	return &Mongo{uri: uri}
}

func (m *Mongo) getFields(collection string, field string, num int) ([]bson.M, error) {
	var err error
	var client *mongo.Client
	var list = []bson.M{}
	var xmap = map[string]bson.M{}
	ctx := context.Background()
	if client, err = mongo.Connect(ctx, m.uri); err != nil {
		return list, err
	}
	defer client.Disconnect(ctx)
	cs, _ := connstring.Parse(m.uri)
	c := client.Database(cs.Database).Collection(collection)
	var doc bson.M
	if err = c.FindOne(ctx, nil).Decode(&doc); err != nil {
		return list, err
	}
	for len(xmap) < num {
		var f interface{}
		b, _ := json.Marshal(bson.M{field: doc[field]})
		json.Unmarshal(b, &f)
		v := make(map[string]interface{})
		util.RandomizeDocument(&v, f, false)
		xmap[v[field].(string)] = v
	}

	for _, v := range xmap {
		list = append(list, v)
	}
	return list, err
}
