package db

import "go.mongodb.org/mongo-driver/mongo"

type DbClient interface {
}

type MongoConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Db       string
	Timeout  int
}

var _ DbClient = (*MongoClient)(nil)

type MongoClient struct {
	client *mongo.Client
}
