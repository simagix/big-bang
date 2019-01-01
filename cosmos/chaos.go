// Copyright 2019 Kuei-chun Chen. All rights reserved.

package cosmos

import (
	"github.com/mongodb/mongo-go-driver/bson"
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
	if c.err != nil {
		return &Chrono{err: c.err}
	}
	chrono := Chrono{config: c.config, verbose: c.verbose, seeds: bson.M{}, numSeed: 10}
	for _, conf := range c.config.Collections {
		for _, v := range conf.Lookup {
			localKey := conf.Name + "/" + v.LocalField
			foreignKey := v.From + "/" + v.ForeignField
			mongo := Source(c.config.Source)
			var list []bson.M
			if list, err = mongo.getFields(v.From, v.ForeignField, 10); err != nil {
				chrono.err = err
				return &chrono
			}
			chrono.seeds[foreignKey] = list
			if list, err = mongo.getFields(conf.Name, v.LocalField, 10); err != nil {
				chrono.err = err
				return &chrono
			}
			chrono.seeds[localKey] = list
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
