// Copyright 2018 Kuei-chun Chen. All rights reserved.

package cosmos

import (
	"testing"
)

func TestBigBang(t *testing.T) {
	chaos := Create("testdata/bigbang.json")
	chaos.SetVerbose(true)
	chrono := chaos.BigBang()
	if chrono.Error() != nil {
		t.Fatal(chrono.Error())
	}
}
