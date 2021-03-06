// Copyright 2018 Kuei-chun Chen. All rights reserved.

package sim

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/simagix/keyhole/sim/util"
)

func TestGetTransactions(t *testing.T) {
	TransactionDoc := GetTransactions("")
	bytes, _ := json.MarshalIndent(TransactionDoc, "", "  ")
	t.Log(string(bytes))
}

func TestExecTXForDemo(t *testing.T) {
	var client *mongo.Client
	client = getMongoClient()
	defer client.Disconnect(context.Background())
	c := client.Database(SimDBName).Collection(CollectionName)
	n := execTXForDemo(c, util.GetDemoDoc())
	if n != 5 {
		t.Fatal()
	}
}

func TestExecTXByTemplate(t *testing.T) {
	var client *mongo.Client
	client = getMongoClient()
	defer client.Disconnect(context.Background())
	c := client.Database(SimDBName).Collection(CollectionName)
	n := execTXByTemplate(c, util.GetDemoDoc())
	if n != 5 {
		t.Fatal()
	}
}

func TestExecTXByTemplateAndTX(t *testing.T) {
	var filename = "testdata/transactions.json"
	var client *mongo.Client
	client = getMongoClient()
	defer client.Disconnect(context.Background())
	c := client.Database(SimDBName).Collection(CollectionName)
	tx := GetTransactions(filename)
	n := execTXByTemplateAndTX(c, util.GetDemoDoc(), tx.Transactions)
	if n != 8 {
		t.Fatal()
	}
}
