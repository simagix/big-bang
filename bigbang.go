// Copyright 2019 Kuei-chun Chen. All rights reserved.

package main

import (
	"flag"

	"github.com/simagix/big-bang/cosmos"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	var err error

	filename := flag.String("config", "", "config file")
	verbose := flag.Bool("v", false, "verbose mode")
	flag.Parse()
	flagset := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { flagset[f.Name] = true })
	chaos := cosmos.Create(*filename)
	chaos.SetVerbose(*verbose)
	var result bson.M
	if err = chaos.BigBang().Expand().Exec(&result); err != nil {
		panic(err)
	}
}
