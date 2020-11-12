package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURL  = "mongodb://localhost:27017"
	maxDocs   = 10
	defaultDB = "tasks"
)

// MongoInstance ...
type MongoInstance struct {
	con *mongo.Client
}

// NewMongoInstance ...
func NewMongoInstance() *MongoInstance {
	return &MongoInstance{}
}

// Connect to MongoDB.
func (s *MongoInstance) Connect(ctx context.Context) error {
	clientOptions := options.Client().ApplyURI(mongoURL)
	var err error
	s.con, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Check the connection with 5s timeout
	c, cn := context.WithTimeout(ctx, 5*time.Second)
	defer cn()

	err = s.con.Ping(c, nil)
	if err != nil {
		return err
	}
	return nil
}

// Disconnect from MongoDB.
func (s *MongoInstance) Disconnect(ctx context.Context) error {
	return s.con.Disconnect(ctx)
}

// Insert one document on the specified dbName and colName.
func (s *MongoInstance) Insert(ctx context.Context, doc interface{}, col string) error {
	db := s.con.Database(defaultDB).Collection(col)
	raw, err := bson.Marshal(doc)
	if err != nil {
		return err
	}

	_, err = db.InsertOne(ctx, raw)
	if err != nil {
		return err
	}
	return nil
}

// ListAll returns 'maxDocs' document from specific collection.
func (s *MongoInstance) ListAll(ctx context.Context, col string) ([]Task, error) {
	findOptions := options.Find()
	findOptions.SetLimit(maxDocs)
	db := s.con.Database(defaultDB).Collection(col)
	results := make([]Task, 0)

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := db.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// DeleteAll ...
func (s *MongoInstance) DeleteAll(ctx context.Context, col string) error {
	db := s.con.Database(defaultDB).Collection(col)
	_, err := db.DeleteMany(ctx, bson.D{{}})
	if err != nil {
		return err
	}
	return nil
}
