package mg

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnType - type of cluster to connect to
type ConnType string

// Single - No cluster, single instance
const Single ConnType = "Single"

// ReplicaSet - A replica set cluster
const ReplicaSet ConnType = "ReplicaSet"

// Sharded - Sharded database
const Sharded ConnType = "Sharded"

// ConnOpts - options for connecting to a mongodb instance
type ConnOpts struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"dbName"`
}

func (co *ConnOpts) toUri() string {

	if co.User != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
			co.User, co.Password, co.Host, co.Port, co.DbName)
	}
	return fmt.Sprintf("mongodb://%s:%d/%s", co.Host, co.Port, co.DbName)
}

// store - holds mongodb connection handle and information
type store struct {
	client *mongo.Client
	dbName string
}

var mongoStore *store

// var defaultDB = "teak"

// Connect - connects to a mongoDB instance or a cluster based on the
// the options provided
func Connect(gtx context.Context, uri string, defaultDB string) error {
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(gtx, clientOpts)
	if err != nil {
		return err
	}

	mongoStore = &store{
		client: client,
		dbName: defaultDB,
	}
	return nil
}

func ConnectWithOpts(gtx context.Context, opts *ConnOpts) error {
	return Connect(gtx, opts.toUri(), opts.DbName)
}

// C - get a handle to collection in the default database, single letter name
// to have nice way to transition from mgo
func C(collectionName string) *mongo.Collection {
	return CollectionWithDB(mongoStore.dbName, collectionName)
}

// CollectionWithDB - gives a reference to a collection in given database
func CollectionWithDB(db, coll string) *mongo.Collection {
	return mongoStore.client.Database(db).Collection(coll)
}
