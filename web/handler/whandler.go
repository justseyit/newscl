package webhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"newscl/model"
	"newscl/repository"
	"newscl/service/auth"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func isSessionTokenValid(r *http.Request) bool {
	//Get cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Printf("webhandler.LoginHandler: %v\n", err)
	}
	if cookie != nil {
		//log.Println("Cookie found")
		//log.Println(cookie.Value)
		token := cookie.Value
		collection := repository.Mongo.Client.Database("newscl").Collection("tokens")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var dbToken auth.Token
		err := collection.FindOne(ctx, bson.M{"token": token}).Decode(&dbToken)
		if err != nil {
			//log.Println("Token not found")
			return false
		} else if dbToken.Token == token {
			//log.Println("Token found")
			return true
		}
	}
	return false
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if isSessionTokenValid(r) {
			http.Redirect(w, r, "/index", http.StatusSeeOther)
			return
		}
		t, err := template.ParseFiles("static/site/Login.html")
		if err != nil {
			log.Printf("webhandler.LoginHandler: %v\n", err)
		}
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		auth.UserLogin(w, r)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if isSessionTokenValid(r) {
			http.Redirect(w, r, "/index", http.StatusSeeOther)
			return
		}
		t, err := template.ParseFiles("static/site/Register.html")
		if err != nil {
			log.Printf("webhandler.RegisterHandler: %v\n", err)
		}
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		auth.UserSignup(w, r)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if !isSessionTokenValid(r) {
		http.Redirect(w, r, "/login", http.StatusNotFound)
		return
	}
	t, err := template.ParseFiles("static/site/index.html")
	if err != nil {
		log.Printf("webhandler.IndexHandler: %v\n", err)
	}
	t.Execute(w, nil)
}

func SourcesHandler(w http.ResponseWriter, r *http.Request) {
	if !isSessionTokenValid(r) {
		http.Redirect(w, r, "/login", http.StatusNotFound)
		return
	}
	sources, err := repository.Mongo.GetSourcesFromDB()
	if err != nil {
		log.Printf("webhandler.SourcesHandler: %v\n", err)
	}
	t, err := template.ParseFiles("static/html/list/getSources")

	if err != nil {
		log.Printf("webhandler.SourcesHandler: %v\n", err)
	}
	data := struct {
		Title   string
		Sources []model.Source
	}{
		Title:   "Sources",
		Sources: sources,
	}
	t.Execute(w, data)
}

func AddSourceHandler(w http.ResponseWriter, r *http.Request) {
	if !isSessionTokenValid(r) {
		http.Redirect(w, r, "/login", http.StatusNotFound)
		return
	}
	t, err := template.ParseFiles("static/site/Add-Source.html")

	if err != nil {
		log.Printf("webhandler.AddSourceHandler: %v\n", err)
	}
	categories, err := repository.Mongo.GetCategoriesFromDB()
	if err != nil {
		log.Printf("webhandler.AddSourceHandler: %v\n", err)
	}
	providers, err := repository.Mongo.GetProvidersFromDB()
	if err != nil {
		log.Printf("webhandler.AddSourceHandler: %v\n", err)
	}
	languages, err := repository.Mongo.GetLanguagesFromDB()
	if err != nil {
		log.Printf("webhandler.AddSourceHandler: %v\n", err)
	}
	data := struct {
		Categories []model.Category
		Providers  []model.Provider
		Languages  []model.Language
	}{
		Categories: categories,
		Providers:  providers,
		Languages:  languages,
	}

	t.Execute(w, data)
}

func AddLanguageHandler(w http.ResponseWriter, r *http.Request) {
	if !isSessionTokenValid(r) {
		http.Redirect(w, r, "/login", http.StatusNotFound)
		return
	}
	t, err := template.ParseFiles("static/site/Add-Language.html")

	if err != nil {
		log.Printf("webhandler.AddLanguageHandler: %v\n", err)
	}
	t.Execute(w, nil)
}

func AddProviderHandler(w http.ResponseWriter, r *http.Request) {
	if !isSessionTokenValid(r) {
		http.Redirect(w, r, "/login", http.StatusNotFound)
		return
	}
	t, err := template.ParseFiles("static/site/Add-Provider.html")

	if err != nil {
		log.Printf("webhandler.AddProviderHandler: %v\n", err)
	}
	t.Execute(w, nil)
}

func AddCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if !isSessionTokenValid(r) {
		http.Redirect(w, r, "/login", http.StatusNotFound)
		return
	}
	t, err := template.ParseFiles("static/site/Add-Category.html")

	if err != nil {
		log.Printf("webhandler.AddCategoryHandler: %v\n", err)
	}
	t.Execute(w, nil)
}

func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if !isSessionTokenValid(r) {
		http.Redirect(w, r, "/login", http.StatusNotFound)
		return
	}
	categories, err := repository.Mongo.GetCategoriesFromDB()
	if err != nil {
		log.Printf("webhandler.GetCategoriesHandler: %v\n", err)
	}
	t, err := template.ParseFiles("static/site/Categories.html")

	if err != nil {
		log.Printf("webhandler.GetCategoriesHandler: %v\n", err)
	}
	data := struct {
		Categories []model.Category
	}{
		Categories: categories,
	}
	t.Execute(w, data)
}

func GetLanguagesHandler(w http.ResponseWriter, r *http.Request) {
	if !isSessionTokenValid(r) {
		http.Redirect(w, r, "/login", http.StatusNotFound)
		return
	}
	languages, err := repository.Mongo.GetLanguagesFromDB()
	if err != nil {
		log.Printf("webhandler.GetLanguagesHandler: %v\n", err)
	}
	t, err := template.ParseFiles("static/site/Languages.html")

	if err != nil {
		log.Printf("webhandler.GetLanguagesHandler: %v\n", err)
	}
	data := struct {
		Languages []model.Language
	}{
		Languages: languages,
	}
	t.Execute(w, data)
}

func GetProvidersHandler(w http.ResponseWriter, r *http.Request) {
	if !isSessionTokenValid(r) {
		http.Redirect(w, r, "/login", http.StatusNotFound)
		return
	}
	providers, err := repository.Mongo.GetProvidersFromDB()
	if err != nil {
		log.Printf("webhandler.GetProvidersHandler: %v\n", err)
	}
	t, err := template.ParseFiles("static/site/Providers.html")

	if err != nil {
		log.Printf("webhandler.GetProvidersHandler: %v\n", err)
	}
	data := struct {
		Providers []model.Provider
	}{
		Providers: providers,
	}
	t.Execute(w, data)
}

func GetSourcesHandler(w http.ResponseWriter, r *http.Request) {
	if !isSessionTokenValid(r) {
		http.Redirect(w, r, "/login", http.StatusNotFound)
		return
	}

	sources, err := repository.Mongo.GetSourcesFromDB()
	log.Printf("Length of sources: %v\n", len(sources))
	for _, source := range sources {
		log.Printf("Source: %v\n", source)
	}

	if err != nil {
		log.Printf("webhandler.GetSourcesHandler: %v\n", err)
	}
	t, err := template.ParseFiles("static/site/Sources.html")

	if err != nil {
		log.Printf("webhandler.GetSourcesHandler: %v\n", err)
	}
	data := struct {
		Sources []model.Source
	}{
		Sources: sources,
	}

	jsonString, _ := json.Marshal(data.Sources[0])

	fmt.Println(string(jsonString))
	t.Execute(w, data)
}

func GetNewsHandler(w http.ResponseWriter, r *http.Request) {
	if !isSessionTokenValid(r) {
		http.Redirect(w, r, "/login", http.StatusNotFound)
		return
	}
	newsList, err0 := repository.Mongo.GetAllNews()
	if err0 != nil {
		log.Printf("webhandler.GetNewsHandler: %v\n", err0)
	}
	sources, err := repository.Mongo.GetSourcesFromDB()
	log.Printf("Length: %d", len(sources))
	for i := 0; i < len(sources); i++ {
		log.Printf("Source: %s", sources[i].Name)
	}

	if err != nil {
		log.Printf("webhandler.GetNewsHandler: %v\n", err)
	}
	t, err := template.ParseFiles("static/site/News.html")

	if err != nil {
		log.Printf("webhandler.GetNewsHandler: %v\n", err)
	}

	data := struct {
		News    []model.NewsClNews
		Sources []model.Source
	}{
		News:    newsList.News,
		Sources: sources,
	}

	t.Execute(w, data)
}

func EditSourceHandler(w http.ResponseWriter, r *http.Request) {
	if !isSessionTokenValid(r) {
		http.Redirect(w, r, "/login", http.StatusNotFound)
		return
	}

	sourceID := r.URL.Query()["id"][0]

	source, err := repository.Mongo.GetSourceFromDB(sourceID)
	if err != nil {
		log.Printf("webhandler.EditSourceHandler: %v\n", err)
	}

	t, err := template.ParseFiles("static/html/editSource.html")

	if err != nil {
		log.Printf("webhandler.EditSourceHandler: %v\n", err)
	}
	categories, err := repository.Mongo.GetCategoriesFromDB()
	if err != nil {
		log.Printf("webhandler.EditSourceHandler: %v\n", err)
	}
	providers, err := repository.Mongo.GetProvidersFromDB()
	if err != nil {
		log.Printf("webhandler.EditSourceHandler: %v\n", err)
	}
	languages, err := repository.Mongo.GetLanguagesFromDB()
	if err != nil {
		log.Printf("webhandler.EditSourceHandler: %v\n", err)
	}
	data := struct {
		Categories []model.Category
		Providers  []model.Provider
		Languages  []model.Language
		Source     model.Source
	}{
		Categories: categories,
		Providers:  providers,
		Languages:  languages,
		Source:     source,
	}

	t.Execute(w, data)
}

func AccountHandler(w http.ResponseWriter, r *http.Request) {
	if !isSessionTokenValid(r) {
		http.Redirect(w, r, "/login", http.StatusNotFound)
		return
	}
	t, err := template.ParseFiles("static/site/Account.html")

	if err != nil {
		log.Printf("webhandler.AccountHandler: %v\n", err)
	}
	t.Execute(w, nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionToken")
	if err != nil {
		log.Printf("webhandler.LogoutHandler: %v\n", err)
	}
	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/html/404.html")

	if err != nil {
		log.Printf("webhandler.NotFoundHandler: %v\n", err)
	}
	t.Execute(w, nil)
}
