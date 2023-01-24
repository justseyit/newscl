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
		}else{
			
			return false, err
		}
		
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
		}else{
			
			return false, err
		}
	}
	return true, nil
}

func (m *MongoDB) GetProvidersFromDB() ([]model.Provider, error) {
	collection := m.Client.Database("newscl").Collection("providers")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var providers []model.Provider
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return providers, err
	}
	cursor.All(ctx, &providers)
	return providers, nil
}

func (m *MongoDB) GetLanguagesFromDB() ([]model.Language, error) {
	collection := m.Client.Database("newscl").Collection("languages")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var languages []model.Language
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return languages, err
	}
	cursor.All(ctx, &languages)
	return languages, nil
}

func (m *MongoDB) GetCategoriesFromDB() ([]model.Category, error) {
	collection := m.Client.Database("newscl").Collection("categories")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var categories []model.Category
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return categories, err
	}
	cursor.All(ctx, &categories)
	return categories, nil
}


func (m *MongoDB) AggregateNews(){
	//Update all of the news,
	/*
	* 1. Get all of the news
	* 2. Add a new field called "sourceID"
	* 3. If news's sourceName is BBC, then sourceID is 1
	* 4. If news's sourceName is Reuters, then sourceID is 0
	* 5. Update the news
	*/

	//1. Get all of the news
	collection := m.Client.Database("newscl").Collection("news")
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	var newsList model.NewsClNewsList
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	cursor.All(ctx, &newsList.News)

	//2. Add a new field called "sourceID"
	for i := 0; i < len(newsList.News); i++ {
		newsList.News[i].SourceID = ""
	}

	//3. If news's sourceName is BBC, then sourceID is 1
	for i := 0; i < len(newsList.News); i++ {
		if newsList.News[i].SourceName == "BBC" {
			newsList.News[i].SourceID = "1"
		}
	}

	//4. If news's sourceName is Reuters, then sourceID is 0
	for i := 0; i < len(newsList.News); i++ {
		if newsList.News[i].SourceName == "Reuters" {
			newsList.News[i].SourceID = "0"
		}
	}

	//5. Update the news
	for i := 0; i < len(newsList.News); i++ {
		filter := bson.D{{Key: "_id", Value: newsList.News[i].ID}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "sourceID", Value: newsList.News[i].SourceID}}}}
		_, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Fatal(err)
		}
	}
}