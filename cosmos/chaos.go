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
	"github.com/simagix/keyhole/mdb"
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
		chrono.err = c.err
		return &chrono
	}
	if chrono.sourceClient, err = mdb.NewMongoClient(c.config.Source.URI, c.config.Source.CAFile, c.config.Source.ClientPEM); err != nil {
		if c.verbose == true {
			log.Println("connecting to", c.config.Source.URI, " failed: ", err)
		}
		return &Chrono{err: err}
	}
	if chrono.targetClient, err = mdb.NewMongoClient(c.config.Target.URI, c.config.Source.CAFile, c.config.Source.ClientPEM); err != nil {
		if c.verbose == true {
			log.Println("connecting to", c.config.Target.URI, " failed: ", err)
		}
		return &Chrono{err: err}
	}

	for _, conf := range c.config.Collections {
		for _, v := range conf.Lookup {
			var list []interface{}
			var templ1 bson.M
			var templ2 bson.M
			if templ1, err = c.getTemplate(chrono.sourceClient, v.From); err != nil {
				if c.verbose == true {
					log.Println("getTemplate from collection", v.From, "failed: ", err)
				}
				chrono.err = err
				return &chrono
			}
			if templ2, err = c.getTemplate(chrono.sourceClient, conf.Name); err != nil {
				if c.verbose == true {
					log.Println("getTemplate from collection", conf.Name, "failed: ", err)
				}
				chrono.err = err
				return &chrono
			}
			if list, err = c.getFields(templ2, v.LocalField, v.NumSeeds); err != nil {
				chrono.err = err
				if c.verbose {
					log.Println("getFields from", v.LocalField, v.NumSeeds, "failed", err)
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
	if c.verbose {
		log.Println(client.ConnectionString(), cs.Database, collection)
	}
	if err != nil {
		return nil, err
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
