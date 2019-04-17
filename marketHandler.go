package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func marketHandler(res http.ResponseWriter, req *http.Request) {
	findParam := regexp.MustCompile("^/market/([a-zA-Z0-9]+)$")
	param := findParam.FindStringSubmatch(req.URL.Path)
	if len(param) == 0 {
		servFile(res, req)
		return
	}
	p := param[1]
	client, err := mongo.NewClient(options.Client().ApplyURI(MONGOCONST.mongoUrl))
	if err != nil {
		log.Fatal(err)
	}
	ctx, can := context.WithTimeout(context.Background(), 10*time.Second)
	defer can() // close connection
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	col := client.Database(MONGOCONST.dbName).Collection(MONGOCONST.marketGroupsCol)
	var pipeline []bson.D
	if p == "null" {
		pipeline = mongo.Pipeline{
			{{"$match", bson.D{{"parentGroupID", bson.D{{"$exists", false}}}}}},
			{{"$sort", bson.D{{"marketGroupName.en-us", 1}}}},
		}
	} else {
		p, _ := strconv.Atoi(p)
		pipeline = mongo.Pipeline{
			{{"$match", bson.D{{"parentGroupID", p}}}},
			{{"$sort", bson.D{{"marketGroupName.en-us", 1}}}},
		}

	}
	cur, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}

	firstChunk := true
	res.Header().Set("Content-Type", "application/json")
	res.Write([]byte("["))
	defer cur.Close(context.Background())
	for {
		var result bson.M
		if cur.Next(context.Background()) == false {
			res.Write([]byte("]"))
			client.Disconnect(ctx)
			break
		}
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		b, err := json.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		if firstChunk {
			res.Write(b)
			firstChunk = false
		} else {
			res.Write([]byte(","))
			res.Write(b)
		}
	}
}
