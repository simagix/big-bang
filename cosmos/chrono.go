// Copyright 2019 Kuei-chun Chen. All rights reserved.

package cosmos

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/network/connstring"
	"github.com/simagix/keyhole/sim/util"
)

// Chrono contains chronical events/info
type Chrono struct {
	config         Config
	err            error
	isExpand       bool
	isSeed         bool
	seedsMap       bson.M
	targetClient   *mongo.Client
	templateSource string
	templateLookup string
	verbose        bool
}

func (c *Chrono) Error() error {
	return c.err
}

// Exec execute and unmarshal results
func (c *Chrono) Exec(result interface{}) error {
	var err error
	var ctx = context.Background()
	var cs connstring.ConnString
	if c.err != nil {
		return c.err
	}
	var res = bson.M{"expand": false, "seed": false}
	res["expand"] = c.isExpand
	res["seed"] = c.isSeed

	if cs, err = connstring.Parse(c.targetClient.ConnectionString()); err != nil {
		return err
	}

	for k, v := range c.seedsMap {
		var data = v.(bson.M)
		var total = data["$total"].(int64)
		num := total
		if c.isSeed == true {
			num = 1
		}
		var doc = data["$template"].(bson.M)
		var mdocs []interface{}
		delete(doc, "_id")
		b, _ := json.Marshal(doc)
		collection := c.targetClient.Database(cs.Database).Collection(k)
		for i := int64(0); i < num; i++ {
			var f interface{}
			json.Unmarshal(b, &f)
			v := make(map[string]interface{})
			util.RandomizeDocument(&v, f, false)
			for attr, val := range data {
				if attr == "$template" || attr == "$total" {
					continue
				}
				length := int64(len(val.([]interface{})))
				v[attr] = val.([]interface{})[i%length]
			}
			if v["_id"] == nil {
				v["_id"] = primitive.NewObjectID()
			}
			mdocs = append(mdocs, v)
			if len(mdocs) > 100 {
				if _, err = collection.InsertMany(ctx, mdocs); err != nil {
					log.Println(i, err)
				}
				mdocs = []interface{}{}
			}
		}
		if _, err = collection.InsertMany(ctx, mdocs); err != nil {
			log.Println(err)
		}
		res[cs.Database+"."+k] = num
	}
	b, _ := json.Marshal(res)
	json.Unmarshal(b, result)
	return nil
}

// Expand inflates collections with data
func (c *Chrono) Expand() *Chrono {
	if c.err != nil {
		return &Chrono{err: c.err}
	}
	c.isExpand = true
	c.isSeed = false
	return c
}

// Seed seeds all defined collections
func (c *Chrono) Seed() *Chrono {
	if c.err != nil {
		return &Chrono{err: c.err}
	}
	if c.isExpand == false {
		c.isSeed = true
	}
	return c
}

// SetVerbose sets verbose mode
func (c *Chrono) SetVerbose(verbose bool) {
	c.verbose = verbose
}
