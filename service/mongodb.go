package service

import (
	"context"
	"log"
	"time"

	model "newscl/model"

	"go.mongodb.org/mongo-driver/bson"
	mongocl "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client *mongocl.Client
}

func NewMongoDB() *MongoDB {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://newscl:newscl@news.2zlcqmf.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongocl.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	return &MongoDB{
		Client: client,
	}
}

func (m *MongoDB) SendNews(newsList model.NewsClNewsList) error {
	collection := m.Client.Database("newscl").Collection("news")
	options := options.InsertOne()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, news := range newsList.News {
		result, err := collection.InsertOne(ctx, news, options)
		if err != nil {
			return err
		}
		log.Println(result)
	}
	return nil
}

func (m *MongoDB) GetAllNews() (model.NewsClNewsList, error) {
	collection := m.Client.Database("newscl").Collection("news")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var newsList model.NewsClNewsList
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return newsList, err
	}
	cursor.All(ctx, &newsList.News)
	return newsList, nil
}

func (m *MongoDB) GetNewsByProvider(provider model.Provider) (model.NewsClNewsList, error) {
	collection := m.Client.Database("newscl").Collection("news")
	filter := bson.D{{Key: "source_name", Value: provider}}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var newsList model.NewsClNewsList
	cursor, err := collection.Find(ctx, filter)
	cursor.All(ctx, &newsList.News)
	if err != nil {
		return newsList, err
	}
	return newsList, nil
}

func (m *MongoDB) GetNewsByLanguage(language string) (model.NewsClNewsList, error) {
	collection := m.Client.Database("newscl").Collection("news")
	filter := bson.D{{Key: "language", Value: language}}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var newsList model.NewsClNewsList
	cursor, err := collection.Find(ctx, filter)
	cursor.All(ctx, &newsList.News)
	if err != nil {
		return newsList, err
	}
	return newsList, nil
}

type newsJustID struct {
	ID string `bson:"_id" json:"id"`
}

func (m *MongoDB) FlagSucceeded(news model.NewsClNews) error {
	isAlreadyFlaggedAsSucceeded, err := m.IsSucceeded(news.ID)
	if err != nil {
		return err
	}
	if !isAlreadyFlaggedAsSucceeded {
		collection := m.Client.Database("newscl").Collection("succeeded")
		bsonDocument := bson.D{{Key: "_id", Value: news.ID}}
		options := options.InsertOne()
		options.SetBypassDocumentValidation(false)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		result, err := collection.InsertOne(ctx, bsonDocument, options)
		if err != nil {
			return err
		}
		log.Printf("Succeeded: %v", result)
	} else {
		log.Printf("Already flagged as succeeded: %v", news.ID)
	}
	return nil
}

func (m *MongoDB) FlagFailed(news model.NewsClNews) error {
	isAlreadyFlaggedAsFailed, err := m.IsFailed(news.ID)
	if err != nil {
		return err
	}
	if !isAlreadyFlaggedAsFailed {
		collection := m.Client.Database("newscl").Collection("failed")
		bsonDocument := bson.D{{Key: "_id", Value: news.ID}}
		options := options.InsertOne()
		options.SetBypassDocumentValidation(false)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		result, err := collection.InsertOne(ctx, bsonDocument, options)
		if err != nil {
			return err
		}
		log.Printf("Failed: %v", result)
	} else {
		log.Printf("Already flagged as failed: %v", news.ID)
	}
	return nil
}

func (m *MongoDB) IsSucceeded(id string) (bool, error) {
	var newsJustID newsJustID
	collection := m.Client.Database("newscl").Collection("succeeded")
	filter := bson.D{{Key: "_id", Value: id}}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(&newsJustID)
	if err != nil {
		if err == mongocl.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *MongoDB) IsFailed(id string) (bool, error) {
	var newsJustID newsJustID
	collection := m.Client.Database("newscl").Collection("failed")
	filter := bson.D{{Key: "_id", Value: id}}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(&newsJustID)
	if err != nil {
		if err == mongocl.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
