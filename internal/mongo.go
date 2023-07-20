package internal

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// dal represents a data access layer. Implemented for mongodb only
type dal interface {
	saveWorker(workerData) error
	saveEvent(eventsBatch) error
	getAllWorkers() ([]workerData, error)
	getWorker(int) (workerData, error)
}

// mongoDBDAL is a struct representing the MongoDB Data Access Layer.
type mongoDBDAL struct {
	client          *mongo.Client
	databaseName    string
	workersColl     *mongo.Collection
	eventsBatchColl *mongo.Collection
}

// ensure interface
var _ dal = &mongoDBDAL{}

// NewMongoDBDAL creates a new instance of mongoDBDAL.
func NewMongoDBDAL(ctx context.Context, uri, databaseName string) (*mongoDBDAL, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return &mongoDBDAL{
		client:          client,
		databaseName:    databaseName,
		workersColl:     client.Database(databaseName).Collection("workers"),
		eventsBatchColl: client.Database(databaseName).Collection("events"),
	}, nil
}

// Close closes the MongoDB connection.
func (dal *mongoDBDAL) Close(ctx context.Context) error {
	return dal.client.Disconnect(ctx)
}

// saveWorker saves worker data to the "workers" collection.
func (dal *mongoDBDAL) saveWorker(workerData workerData) error {
	_, err := dal.workersColl.InsertOne(context.Background(), workerData)
	if err != nil {
		fmt.Println("Error saving worker:", err)
		return err
	}

	return nil
}

// saveEvent saves an eventsBatch to the "events" collection.
// Consider using transaction.
func (dal *mongoDBDAL) saveEvent(batch eventsBatch) error {
	_, err := dal.eventsBatchColl.InsertOne(context.Background(), batch)
	if err != nil {
		fmt.Println("Error saving event:", err)
	}

	// Update the "eventsCount" in the "workers" collection for the corresponding workerID
	filter := bson.D{{"id", batch.WorkerId}}
	update := bson.D{{"$inc", bson.D{{"EventsCount", len(batch.Events)}}}}
	_, err = dal.workersColl.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println("Error updating worker's eventsCount:", err)
	}
}

// getAllWorkers gets all workers from the "workers" collection.
func (dal *mongoDBDAL) getAllWorkers() []workerData {
	var workers []workerData

	cursor, err := dal.workersColl.Find(context.Background(), bson.D{})
	if err != nil {
		fmt.Println("Error fetching workers:", err)
		return workers
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var worker workerData
		if err := cursor.Decode(&worker); err != nil {
			fmt.Println("Error decoding worker:", err)
		} else {
			workers = append(workers, worker)
		}
	}

	if err := cursor.Err(); err != nil {
		fmt.Println("Cursor error:", err)
	}

	return workers
}

// getWorker gets a worker from the "workers" collection by the provided workerID.
func (dal *mongoDBDAL) getWorker(workerID int) workerData {
	var worker workerData
	filter := bson.D{{"id", workerID}}

	err := dal.workersColl.FindOne(context.Background(), filter).Decode(&worker)
	if err != nil {
		fmt.Println("Error fetching worker:", err)
		return worker
	}

	return worker
}
