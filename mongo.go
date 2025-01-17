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
	maxDocs   = 50
	defaultDB = "tasks"
)

// MongoInstance stores a single and shared instance established during service initilization.
// A mongo.Client already implements a pool of connections to a MongoDB deployment by default,
// beign safe for concurrent use by multiple goroutines.
type MongoInstance struct {
	con *mongo.Client
}

// Connect establishes connection to MongoDB, testing connection with a 5sec timeout, returning
// an error if failed on ping.
func (s *MongoInstance) Connect(ctx context.Context) error {
	clientOptions := options.Client().ApplyURI(mongoURL)
	var err error
	s.con, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// check the connection with 5s timeout
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

// Insert one document on the 'col' collection from defaultDB.
func (s *MongoInstance) Insert(ctx context.Context, col string, doc interface{}) error {
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

// Update a document matching a specific 'id'. Both 'desc' and 'status' fields must always be updated.
func (s *MongoInstance) Update(ctx context.Context, col string, taskID int, desc, status string) error {
	filter := bson.D{{Key: "id", Value: taskID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "desc", Value: desc},
		}},
		{Key: "$set", Value: bson.D{
			{Key: "status", Value: status},
		}},
	}

	db := s.con.Database(defaultDB).Collection(col)
	_, err := db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// UpdateFinished updates the 'finished' field of a task matching a specific 'id'.
func (s *MongoInstance) UpdateFinished(ctx context.Context, col string, taskID int, finished bool) error {
	filter := bson.D{{Key: "id", Value: taskID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "finished", Value: finished},
		}},
	}

	db := s.con.Database(defaultDB).Collection(col)
	_, err := db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// Count returns the number of documents on a certain collection.
func (s *MongoInstance) Count(ctx context.Context, col string) (int, error) {
	db := s.con.Database(defaultDB).Collection(col)
	n, err := db.CountDocuments(ctx, bson.D{{}})
	if err != nil {
		return 0, err
	}

	// will never overflow in this simple api
	return int(n), nil
}

// List returns 'maxDocs' documents from 'col' collection with a certain status tag.
func (s *MongoInstance) List(ctx context.Context, col string, status []string) ([]Task, error) {
	findOptions := options.Find()
	findOptions.SetLimit(maxDocs)
	db := s.con.Database(defaultDB).Collection(col)

	results := make([]Task, 0)
	filter := bson.D{
		{Key: "status", Value: bson.D{
			{Key: "$in", Value: status},
		}},
	}

	cur, err := db.Find(ctx, filter, findOptions)
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

// Delete at most one task matching a certain 'id' from 'col' collection.
func (s *MongoInstance) Delete(ctx context.Context, col string, taskID int) error {
	db := s.con.Database(defaultDB).Collection(col)
	_, err := db.DeleteOne(ctx, bson.D{{Key: "id", Value: taskID}})
	if err != nil {
		return err
	}
	return nil
}

// DeleteAll deletes all documents from 'col' collection.
func (s *MongoInstance) DeleteAll(ctx context.Context, col string) error {
	db := s.con.Database(defaultDB).Collection(col)
	_, err := db.DeleteMany(ctx, bson.D{{}})
	if err != nil {
		return err
	}
	return nil
}
