package model

type Provider struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

type Category struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

type Language struct {
	ID           string `json:"id" bson:"_id"`
	Name         string `json:"name" bson:"name"`
	LanguageCode string `json:"languageCode" bson:"languageCode"`
	CountryCode  string `json:"countryCode" bson:"countryCode"`
}

type SourceRef struct {
	ID             string `json:"id" bson:"_id"`
	Name           string `json:"name" bson:"name"`
	URL            string `json:"url" bson:"url"`
	FetchFrequency uint   `json:"frequency" bson:"frequency"`
	CategoryID     string `json:"categoryID" bson:"categoryID"`
	ProviderID     string `json:"providerID" bson:"providerID"`
	LanguageID     string `json:"languageID" bson:"languageID"`
}

type Source struct {
	ID             string   `json:"id" bson:"_id"`
	Name           string   `json:"name" bson:"name"`
	URL            string   `json:"url" bson:"url"`
	Category       Category `json:"category" bson:"category"`
	Provider       Provider `json:"provider" bson:"provider"`
	Language       Language `json:"language" bson:"language"`
	FetchFrequency uint     `json:"frequency" bson:"frequency"`
}
