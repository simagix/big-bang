// Copyright 2019 Kuei-chun Chen. All rights reserved.

package cosmos

// Chrono contains chronical events/info
type Chrono struct {
	err     error
	verbose bool
}

// SetVerbose sets verbose mode
func (c *Chrono) SetVerbose(verbose bool) {
	c.verbose = verbose
}

func (c *Chrono) Error() error {
	return c.err
}

// Expand inflates collections with data
func (c *Chrono) Expand() *Chrono {
	if c.err != nil {
		return &Chrono{err: c.err}
	}
	return &Chrono{}
}

// Seed seeds all defined collections
func (c *Chrono) Seed() *Chrono {
	if c.err != nil {
		return &Chrono{err: c.err}
	}
	return &Chrono{}
}

// Exec execute and unmarshal results
func (c *Chrono) Exec(result interface{}) error {
	if c.err != nil {
		return c.err
	}
	return nil
}
