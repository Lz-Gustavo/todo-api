package main

import (
	"context"
	"testing"
	"time"
)

func TestDBCon(t *testing.T) {
	inst := &MongoInstance{}
	ctx, cn := context.WithTimeout(context.TODO(), time.Second)
	defer cn()

	inst.Connect(ctx)
	err := inst.Disconnect(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
