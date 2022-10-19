package model

import (
	"carrot-backyard/param"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const tokenCollectionName = "token"

func (m *model) tokenCollection() *mongo.Collection {
	return m.database.Collection(tokenCollectionName)
}

type TokenInterface interface {
	GetToken() (string, error)
}

var (
	defaultToken = "carrot"
	Token        = defaultToken
)

func (m *model) GetToken() (string, error) {
	cursor, err := m.tokenCollection().Find(m.context, bson.M{})
	if err != nil {
		return defaultToken, err
	}

	var tokens []param.TokenStruct
	if err = cursor.All(m.context, &tokens); err != nil {
		return defaultToken, err
	}
	if len(tokens) == 0 {
		return defaultToken, err
	}

	return tokens[len(tokens)-1].Token, nil
}
