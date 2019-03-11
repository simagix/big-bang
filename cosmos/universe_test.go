// Copyright 2018 Kuei-chun Chen. All rights reserved.

package cosmos

import (
	"errors"
	"testing"

	"github.com/simagix/keyhole/mdb"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreate(t *testing.T) {
	chaos := Create("testdata/bigbang-string.json")
	chaos.SetVerbose(true)
	if chaos.Error() != nil {
		t.Log(mdb.Stringify(chaos.config, "", "  "))
		t.Fatal(errors.New("config is"), chaos.config)
	}
	b, _ := bson.Marshal(chaos.config)
	var doc bson.M
	bson.Unmarshal(b, &doc)
	var cfg Config
	bson.Unmarshal(b, &cfg)

	if cfg.Target.URI != doc["target"].(bson.M)["uri"] {
		t.Fatal("expected", doc["target"].(bson.M)["uri"], ", but got", cfg.Target.URI)
	}
}
