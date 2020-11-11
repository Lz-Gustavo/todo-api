package main

import (
	"context"
	"testing"
)

func TestDBCon(t *testing.T) {
	inst := &MongoInstance{}
	ctx := context.TODO()

	inst.Connect(ctx)
	err := inst.Disconnect(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
