// Copyright 2019 Kuei-chun Chen. All rights reserved.

package cosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/network/connstring"
	"github.com/simagix/keyhole/sim/util"
)

// Chaos contains configurations
type Chaos struct {
	config   Config
	err      error
	filename string
	verbose  bool
}

// BigBang builds relationships among collections
func (c *Chaos) BigBang() *Chrono {
	var err error
	chrono := Chrono{config: c.config, verbose: c.verbose, seedsMap: bson.M{}}

	if c.err != nil {
		chrono.err = err
		return &chrono
	}
	ctx := context.Background()
	if chrono.sourceClient, err = mongo.Connect(ctx, c.config.Source); err != nil {
		return &Chrono{err: c.err}
	}
	if chrono.targetClient, err = mongo.Connect(ctx, c.config.Target); err != nil {
		return &Chrono{err: c.err}
	}

	for _, conf := range c.config.Collections {
		for _, v := range conf.Lookup {
			var list []interface{}
			var templ1 bson.M
			var templ2 bson.M
			if templ1, err = c.getTemplate(chrono.sourceClient, v.From); err != nil {
				chrono.err = err
				return &chrono
			}
			if templ2, err = c.getTemplate(chrono.sourceClient, conf.Name); err != nil {
				chrono.err = err
				return &chrono
			}
			if list, err = c.getFields(templ2, v.LocalField, v.NumSeeds); err != nil {
				chrono.err = err
				if c.verbose {
					log.Println("get source field error", err)
				}
				return &chrono
			}
			doc1 := bson.M{"$template": templ1, "$total": v.Total, v.ForeignField: list}
			chrono.seedsMap[v.From] = doc1
			doc2 := bson.M{"$template": templ2, "$total": conf.Total, v.LocalField: list}
			chrono.seedsMap[conf.Name] = doc2
			// log.Println(mdb.Stringify(chrono.seedsMap, "", "  "))

		}
	}
	return &chrono
}

func (c *Chaos) Error() error {
	return c.err
}

// SetVerbose sets verbose mode
func (c *Chaos) SetVerbose(verbose bool) {
	c.verbose = verbose
}

func (c *Chaos) getTemplate(client *mongo.Client, collection string) (bson.M, error) {
	var err error
	cs, err := connstring.Parse(client.ConnectionString())
	if err != nil {
		return nil, err
	}

	if c.verbose {
		log.Println(client.ConnectionString(), cs.Database, collection)
	}

	var doc bson.M
	coll := client.Database(cs.Database).Collection(collection)
	err = coll.FindOne(context.Background(), nil).Decode(&doc)
	return doc, err
}

func (c *Chaos) getFields(doc bson.M, field string, num int) ([]interface{}, error) {
	if c.verbose {
		log.Println("{field, num}", field, num)
	}
	var err error
	var list []interface{}
	var xmap = map[string]bson.M{}
	for len(xmap) < num {
		var f interface{}
		b, _ := json.Marshal(bson.M{field: doc[field]})
		json.Unmarshal(b, &f)
		v := make(map[string]interface{})
		util.RandomizeDocument(&v, f, false)
		xmap[fmt.Sprintf("%v", v[field])] = v
	}

	for _, v := range xmap {
		list = append(list, v[field])
	}
	return list, err
}
