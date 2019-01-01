// Copyright 2019 Kuei-chun Chen. All rights reserved.

package cosmos

import (
	"context"
	"encoding/json"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/x/network/connstring"
)

// Chrono contains chronical events/info
type Chrono struct {
	config       Config
	err          error
	isExpand     bool
	numSeed      int
	isSeed       bool
	seedsMap     bson.M
	sourceClient *mongo.Client
	targetClient *mongo.Client
	verbose      bool
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
	var doc = bson.M{"expand": false, "seed": false}
	doc["expand"] = c.isExpand
	doc["seed"] = c.isSeed

	if cs, err = connstring.Parse(c.targetClient.ConnectionString()); err != nil {
		return err
	}
	for k, v := range c.seedsMap {
		var data = v.(bson.M)
		var doc = data["template"].(bson.M)
		delete(doc, "_id")
		for attr, val := range data {
			if attr == "template" {
				continue
			}
			doc[attr] = val.([]interface{})[0]
		}
		if doc["_id"] == nil {
			doc["_id"] = primitive.NewObjectID()
		}
		opts := options.Replace()
		opts.SetUpsert(true)
		if _, err = c.targetClient.Database(cs.Database).Collection(k).ReplaceOne(ctx, bson.M{"_id": doc["_id"]}, doc, opts); err != nil {
			return err
		}
	}
	if c.isExpand == true {
	}
	b, _ := json.Marshal(doc)
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
