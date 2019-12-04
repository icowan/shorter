/**
 * @Time : 2019-11-15 14:41
 * @Author : solacowa@gmail.com
 * @File : repository
 * @Software: GoLand
 */

package mongodb

import (
	"context"
	"github.com/icowan/shorter/pkg/service"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func (m *mongoRepository) Exists(has string) (exists bool, err error) {
	return
}

func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return client, nil
}

func NewMongoRepository(mongoURL, mongoDB string, mongoTimeout int) (service.Repository, error) {
	repo := &mongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}
	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.mongodb.NewMongoRepository")
	}
	repo.client = client

	return repo, nil
}

func (m *mongoRepository) Find(code string) (redirect *service.Redirect, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	collection := m.client.Database(m.database).Collection("redirects")
	filter := bson.M{"code": code}

	if err = collection.FindOne(ctx, filter).Decode(&redirect); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(err, "repository.Redirect.Find")
		}
	}

	return
}

func (m *mongoRepository) Store(redirect *service.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	collection := m.client.Database(m.database).Collection("redirects")
	_, err := collection.InsertOne(ctx, bson.M{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	})

	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}

	return nil
}
