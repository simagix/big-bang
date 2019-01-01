// Copyright 2018 Kuei-chun Chen. All rights reserved.

package cosmos

import (
	"testing"
)

func TestGetFields(t *testing.T) {
	mongo := Source("mongodb://localhost/keyhole?replicaSet=replset")
	list, err := mongo.getFields("dealers", "_id", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 10 {
		t.Fatal(err)
	}
}
