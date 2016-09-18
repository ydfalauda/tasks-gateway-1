package main

import (
	"errors"

	"github.com/twinj/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DB struct {
	Session         *mgo.Session
	DBName          string
	Collection      string
	TokenCollection string
	TaskCollection  string
}

func NewDB(address string, dbname string, collection string, token string, task string) (*DB, error) {
	session, err := mgo.Dial(address)
	if err != nil {
		return nil, err
	}
	return &DB{
		Session:         session,
		DBName:          dbname,
		Collection:      collection,
		TokenCollection: token,
		TaskCollection:  task,
	}, nil
}

func (db *DB) GetCollection() *mgo.Collection {
	database := db.Session.DB(db.DBName)
	return database.C(db.Collection)
}

func (db *DB) GetTokenCollection() *mgo.Collection {
	database := db.Session.DB(db.DBName)
	return database.C(db.TokenCollection)
}

func (db *DB) GetTaskCollection() *mgo.Collection {
	database := db.Session.DB(db.DBName)
	return database.C(db.TaskCollection)
}

func (db *DB) Signup(name string, username string, password string) (*User, error) {

	collection := db.GetCollection()

	number, err := collection.Find(bson.M{
		"username": username,
	}).Count()
	if err != nil {
		return nil, err
	}
	if number > 0 {
		return nil, errors.New("Username already registered...")
	}

	user := User{
		ID:       bson.NewObjectId(),
		Username: username,
		Name:     name,
		Password: password,
	}
	err = collection.Insert(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *DB) Login(username string, password string) (*User, error) {
	collection := db.GetCollection()
	number, err := collection.Find(bson.M{
		"username": username,
		"password": password,
	}).Count()
	if err != nil {
		return nil, err
	}
	if number <= 0 {
		return nil, errors.New("Username/Password invalid")
	}

	user := User{}
	collection.Find(bson.M{
		"username": username,
		"password": password,
	}).One(&user)

	return &user, nil
}

func (db *DB) GetToken(user *User) (*AccessToken, error) {
	collection := db.GetTokenCollection()
	number, err := collection.Find(bson.M{
		"_id": user.ID,
	}).Count()
	if err != nil {
		return nil, err
	}
	var token *AccessToken
	if number <= 0 {
		token = &AccessToken{
			UserID: user.ID,
			Token:  uuid.NewV4().String(),
		}
		collection.Insert(token)
	} else {
		token = &AccessToken{}
		collection.Find(bson.M{
			"_id": user.ID,
		}).One(&token)
	}
	return token, nil
}

func (db *DB) GetUserByToken(token string) (*User, error) {
	collection := db.GetTokenCollection()
	number, err := collection.Find(bson.M{
		"token": token,
	}).Count()
	if err != nil {
		return nil, err
	}

	if number <= 0 {
		return nil, errors.New("Not authorized")
	}
	tokenObj := AccessToken{}
	collection.Find(bson.M{
		"token": token,
	}).One(&tokenObj)

	userCollection := db.GetCollection()
	userObj := User{}
	userCollection.Find(bson.M{
		"_id": tokenObj.UserID,
	}).One(&userObj)
	return &userObj, nil
	// collection

}

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Username string        `json:"username"`
	Name     string        `json:"name"`
	Password string        `json:"password"`
}

type AccessToken struct {
	UserID bson.ObjectId `bson:"_id,omitempty"`
	Token  string        `json:"token"`
}
