// Copyright 2019 Kuei-chun Chen. All rights reserved.

package cosmos

import "github.com/globalsign/mgo/bson"

// Chaos contains configurations
type Chaos struct {
	err      error
	filename string
	config   bson.M
	verbose  bool
}

func (c *Chaos) Error() error {
	return c.err
}

// SetVerbose sets verbose mode
func (c *Chaos) SetVerbose(verbose bool) {
	c.verbose = verbose
}

// BigBang builds relationship among collections
func (c *Chaos) BigBang() *Chrono {
	if c.err != nil {
		return &Chrono{err: c.err}
	}
	return &Chrono{}
}
