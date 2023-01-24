package service

import (
	"context"
	"log"
	"newscl/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func (m *MongoDB) GetSourcesFromDB() ([]model.Source, error) {
	var sourcerefs []model.SourceRef
	var sources []model.Source
	collection := m.Client.Database("newscl").Collection("sourcerefs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return sources, err
	}
	cursor.All(ctx, &sourcerefs)
	for _, sourceref := range sourcerefs {
		collection := m.Client.Database("newscl").Collection("providers")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		var provider model.Provider
		err := collection.FindOne(ctx, bson.M{"_id": sourceref.ProviderID}).Decode(&provider)
		if err != nil {
			return sources, err
		}
		collection = m.Client.Database("newscl").Collection("categories")
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		var category model.Category
		err = collection.FindOne(ctx, bson.M{"_id": sourceref.CategoryID}).Decode(&category)
		if err != nil {
			return sources, err
		}
		collection = m.Client.Database("newscl").Collection("languages")
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		var language model.Language
		err = collection.FindOne(ctx, bson.M{"_id": sourceref.LanguageID}).Decode(&language)
		if err != nil {
			return sources, err
		}
		source := model.Source{
			ID:             sourceref.ID,
			Name:           sourceref.Name,
			URL:            sourceref.URL,
			Category:       category,
			Provider:       provider,
			Language:       language,
			FetchFrequency: sourceref.FetchFrequency,
		}
		sources = append(sources, source)
	}
	return sources, nil
}

func (m *MongoDB) SaveSourcesToDB(sources []model.Source) error {
	collection := m.Client.Database("newscl").Collection("sourcerefs")
	options := options.InsertOne()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, source := range sources {
		sourceref := model.SourceRef{
			ID:             source.ID,
			Name:           source.Name,
			URL:            source.URL,
			FetchFrequency: source.FetchFrequency,
			CategoryID:     source.Category.ID,
			ProviderID:     source.Provider.ID,
			LanguageID:     source.Language.ID,
		}
		_, err := collection.InsertOne(ctx, sourceref, options)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MongoDB) RemoveSourcesFromDB(idList []string) error {
	collection := m.Client.Database("newscl").Collection("sourcerefs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, id := range idList {
		result, err := collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
		if err != nil {
			return err
		}
		log.Println(result)
	}
	return nil
}

func (m *MongoDB) GetSourceFromDB(id string) (model.Source, error) {
	var source model.Source
	collection := m.Client.Database("newscl").Collection("sourcerefs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var sourceref model.SourceRef
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&sourceref)
	if err != nil {
		return source, err
	}
	collection = m.Client.Database("newscl").Collection("providers")
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var provider model.Provider
	err = collection.FindOne(ctx, bson.M{"_id": sourceref.ProviderID}).Decode(&provider)
	if err != nil {
		return source, err
	}
	collection = m.Client.Database("newscl").Collection("categories")
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var category model.Category
	err = collection.FindOne(ctx, bson.M{"_id": sourceref.CategoryID}).Decode(&category)
	if err != nil {
		return source, err
	}
	collection = m.Client.Database("newscl").Collection("languages")
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var language model.Language
	err = collection.FindOne(ctx, bson.M{"_id": sourceref.LanguageID}).Decode(&language)
	if err != nil {
		return source, err
	}
	source = model.Source{
		ID:             id,
		Name:           sourceref.Name,
		URL:            sourceref.URL,
		Category:       category,
		Provider:       provider,
		Language:       language,
		FetchFrequency: sourceref.FetchFrequency,
	}
	return source, nil
}
