package model

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Model interface {
	Close()
	FamilyInterface
	MusterInterface
	BallotInterface
	TokenInterface
}

type model struct {
	database *mongo.Database
	context  context.Context
	cancel   context.CancelFunc
}

//GetModel returns a mongo model which will be used in controller
//to help call for the functions that work on mongoDB
func GetModel() *model {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	mongoModel := &model{
		database: getMongoDataBase(ctx),
		context:  ctx,
		cancel:   cancel,
	}

	return mongoModel
}
