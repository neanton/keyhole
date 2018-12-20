// Copyright 2018 Kuei-chun Chen. All rights reserved.

package keyhole

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Transaction -
type Transaction struct {
	C      string   `json:"c"`
	Filter bson.M   `json:"filter"`
	Op     bson.M   `json:"op"`
	Pipe   []bson.M `json:"pipe"`
}

// TransactionDoc -
type TransactionDoc struct {
	Transactions []Transaction `json:"transactions" bson:"transactions"`
	Indexes      []bson.M      `json:"indexes" bson:"indexes"`
}

// GetTransactions -
func GetTransactions(filename string) TransactionDoc {
	if filename == "" {
		return TransactionDoc{}
	}
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return TransactionDoc{}
	}

	var doc TransactionDoc
	json.Unmarshal(bytes, &doc)
	return doc
}

func execTXForDemo(c *mgo.Collection, doc bson.M) int {
	schema := FavoritesSchema{}
	results := []bson.M{}
	change := bson.M{"$set": bson.M{"ts": time.Now()}}
	bytes, _ := json.Marshal(doc)
	json.Unmarshal(bytes, &schema)
	city := schema.FavoriteCity
	book := schema.FavoriteBook
	movie := schema.FavoriteMovie
	txCount := 0

	c.Insert(doc)
	txCount++
	// c.Find(bson.M{"favoriteCity": city}).Sort("favoriteCity").Limit(512).All(&results)
	// txCount++
	c.Find(bson.M{"favoriteCity": city}).Limit(20).All(&results)
	txCount++
	c.Find(bson.M{"favoriteCity": city, "favoriteBook": book}).One(&results)
	txCount++
	c.Update(bson.M{"_id": doc["_id"]}, change)
	txCount++
	// c.Find(bson.M{"favoriteCity": city, "favoriteBook": book, "FavoriteMovie": movie}).One(&results)
	c.Find(bson.M{"favoritesList": bson.M{"$elemMatch": bson.M{"movie": movie}}}).One(&results)
	txCount++
	// c.Find(bson.M{"favoritesList": bson.M{"$elemMatch": bson.M{"book": book}}}).Limit(100).All(&results)
	// txCount++
	return txCount
}

func execTXByTemplateAndTX(c *mgo.Collection, doc bson.M, transactions []Transaction) int {
	results := []bson.M{}
	var op = make(map[string]interface{})
	var pipe []map[string]interface{}
	var pipeline []map[string]interface{}

	for _, tx := range transactions {
		if tx.C == "insert" {
			c.Insert(doc)
		} else {
			bytes, _ := json.Marshal(tx.Filter)
			cmd := make(map[string]interface{})
			filter := make(map[string]interface{})
			json.Unmarshal(bytes, &cmd)
			RandomizeDocument(&filter, cmd, false)

			if tx.C == "find" {
				c.Find(filter).Limit(20).All(&results)
			} else if tx.C == "findOne" {
				c.Find(filter).One(&results)
			} else if tx.C == "update" {
				bytes, _ = json.Marshal(tx.Op)
				json.Unmarshal(bytes, &op)
				RandomizeDocument(&filter, op, false)
				c.Update(filter, op)
			} else if tx.C == "updateAll" {
				bytes, _ = json.Marshal(tx.Op)
				json.Unmarshal(bytes, &op)
				RandomizeDocument(&filter, op, false)
				c.UpdateAll(filter, op)
			} else if tx.C == "remove" {
				c.Remove(filter)
			} else if tx.C == "removeAll" {
				c.RemoveAll(filter)
			} else if tx.C == "aggregate" {
				bytes, _ := json.Marshal(tx.Pipe)
				json.Unmarshal(bytes, &pipe)
				for _, p := range pipe {
					for k, v := range p {
						if k == "$match" {
							q := make(map[string]interface{})
							RandomizeDocument(&q, v, false)
							pipeline = append(pipeline, bson.M{"$match": q})
						} else {
							pipeline = append(pipeline, p)
						}
					}
				}
				c.Pipe(pipeline).All(&results)
			}
		}
	}

	return len(transactions)
}

func execTXByTemplate(c *mgo.Collection, doc bson.M) int {
	results := []bson.M{}
	change := bson.M{"$set": bson.M{"ts": time.Now()}}
	_id := doc["_id"]
	txCount := 0
	c.Insert(doc)
	txCount++
	c.Find(bson.M{"_id": _id}).One(&results)
	txCount++
	c.Update(bson.M{"_id": _id}, change)
	txCount++
	c.Find(bson.M{"_search": doc["_search"]}).Sort("_search").Limit(10).All(&results)
	txCount++
	c.Remove(bson.M{"_id": _id})
	txCount++
	return txCount
}