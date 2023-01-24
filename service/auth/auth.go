package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"newscl/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var SECRET_KEY = []byte("seyitsecretkey")

type User struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type Token struct {
	ID    string `json:"id" bson:"_id"`
	Token string `json:"token" bson:"token"`
	Exp   string `json:"exp" bson:"exp"`
}

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		log.Println("Error in JWT token generation")
		return "", err
	}
	return tokenString, nil
}

func UserSignup(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	name := request.PostFormValue("name")
	email := request.PostFormValue("email")
	password := request.PostFormValue("password")
	passwordv := request.PostFormValue("passwordv")
	response.Header().Set("Content-Type", "application/json")
	if password != passwordv {
		response.Write([]byte(`{"response":"Passwords do not match!"}`))
		return
	}
	var chkUser User
	var user User
	user.Email = email
	user.Password = password
	user.Name = name

	user.Password = getHash([]byte(user.Password))
	collection := repository.Mongo.Client.Database("newscl").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&chkUser)
	if err == nil {
		response.Write([]byte(`{"response":"User with same email already exists!"}`))
		return
	}
	log.Printf("New User\n Name: %s\n Email: %s\n Password: %s\n", user.Name, user.Email, user.Password)
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)
}

func UserLogin(response http.ResponseWriter, request *http.Request) {
	//Get the cookie
	cookie, err := request.Cookie("token")
	if err != nil {
		//log.Println("No cookie found")
	}
	if cookie != nil {
		//log.Println("Cookie found")
		token := cookie.Value
		collection := repository.Mongo.Client.Database("newscl").Collection("tokens")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var dbToken Token
		err := collection.FindOne(ctx, bson.M{"token": token}).Decode(&dbToken)
		if err != nil {
			//log.Println("Token not found")
		} else if dbToken.Token == token {
			//log.Println("Token found")
			http.Redirect(response, request, "/index", http.StatusSeeOther)
			return
		}
	}
	request.ParseForm()
	email := request.PostForm.Get("emaill")
	password := request.PostForm.Get("passwordl")
	checkbox := request.PostForm.Get("checkbox")
	response.Header().Set("Content-Type", "application/json")
	var user User
	var dbUser User
	user.Email = email
	user.Password = password
	log.Printf("User Login\n Email: %s\n Password: *secret*", user.Email)
	collection := repository.Mongo.Client.Database("newscl").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	errUsr := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbUser)

	if errUsr != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	userPass := []byte(user.Password)
	dbPass := []byte(dbUser.Password)

	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)

	if passErr != nil {
		log.Println(passErr)
		response.Write([]byte(`{"response":"Wrong Password!"}`))
		return
	}
	jwtToken, err := GenerateJWT()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	if checkbox == "on" {
		cookie := http.Cookie{Name: "token", Value: jwtToken, Expires: time.Now().Add(24 * time.Hour)}
		http.SetCookie(response, &cookie)
		collection1 := repository.Mongo.Client.Database("newscl").Collection("tokens")
		ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel1()
		result, _ := collection1.InsertOne(ctx1, bson.M{"token": jwtToken, "uid": dbUser.ID, "exp": time.Now().Add(24 * time.Hour).String()})
		log.Printf("Result: %v\n", result)
	} else {
		cookie := http.Cookie{Name: "token", Value: jwtToken}
		http.SetCookie(response, &cookie)
	}

	log.Println("User Logged In")
	http.Redirect(response, request, "/index", http.StatusSeeOther)
}
