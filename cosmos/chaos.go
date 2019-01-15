// Copyright 2019 Kuei-chun Chen. All rights reserved.

package cosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
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
	if chrono.targetClient, err = mdb.NewMongoClient(c.config.Target.URI, c.config.Target.CAFile, c.config.Target.ClientPEM); err != nil {
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
			if v.ForeignField == "_id" && v.NumSeeds < v.Total {
				v.NumSeeds = v.Total
			}
			if templ1, err = util.GetDocByTemplate(v.Template, true); err != nil {
				if c.verbose == true {
					log.Println("getTemplateFromFile", v.Template, "failed: ", err)
				}
				chrono.err = err
				return &chrono
			}
			if templ2, err = util.GetDocByTemplate(conf.Template, true); err != nil {
				if c.verbose == true {
					log.Println("getTemplateFromFile", conf.Template, "failed: ", err)
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
			if chrono.seedsMap[conf.Name] == nil {
				doc2 := bson.M{"$template": templ2, "$total": conf.Total, v.LocalField: list}
				chrono.seedsMap[conf.Name] = doc2
			} else {
				m := chrono.seedsMap[conf.Name].(primitive.M)
				m[v.LocalField] = list
				chrono.seedsMap[conf.Name] = m
			}
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

func (c *Chaos) getTemplateFromCollection(client *mongo.Client, collection string) (bson.M, error) {
	var err error
	cs, err := connstring.Parse(client.ConnectionString())
	if err != nil {
		return nil, err
	}
	if c.verbose {
		log.Println("getTemplateFromCollection", client.ConnectionString(), cs.Database, collection)
	}

	var doc bson.M
	coll := client.Database(cs.Database).Collection(collection)
	err = coll.FindOne(context.Background(), bson.D{{}}).Decode(&doc)
	return doc, err
}

func (c *Chaos) getFields(doc bson.M, field string, num int64) ([]interface{}, error) {
	if c.verbose {
		log.Println("{field, num}", field, num)
	}
	var err error
	var list []interface{}
	var xmap = map[string]bson.M{}
	var total = 1000
	for len(xmap) < int(num) && total < 2024 {
		total++
		var f interface{}
		value := doc[field]
		_, ok := value.(string)
		if ok {
			value = fmt.Sprintf("%s_%d", value, total)
			b, _ := json.Marshal(bson.M{field: value})
			json.Unmarshal(b, &f)
		} else {
			b, _ := json.Marshal(bson.M{field: total})
			json.Unmarshal(b, &f)
		}
		v := make(map[string]interface{})
		util.RandomizeDocument(&v, f, false)
		xmap[fmt.Sprintf("%v", v[field])] = v
	}

	for _, v := range xmap {
		list = append(list, v[field])
	}
	return list, err
}
