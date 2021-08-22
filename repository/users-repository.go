package repository

import (
	"context"
	"errors"
	"golangBackend/db"
	"golangBackend/form"
	"golangBackend/models"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var UserEntity UserDetails

type usersEntity struct {
	resource *db.Resource
	repo     *mongo.Collection
}

type UserDetails interface {
	Register(userForm form.Users) (*models.Users, int, error)
	GetUserByID(username string) (*models.Users, int, error)
	OnlineUsers() ([]models.Users, int, error)
	GetAll() ([]models.Users, int, error)
}

func NewUserEntity(resource *db.Resource) UserDetails {
	userRepo := resource.DB.Collection("userDetails")
	UserEntity = &usersEntity{resource: resource, repo: userRepo}

	return UserEntity
}

func (entity *usersEntity) OnlineUsers() ([]models.Users, int, error) {
	usersList := []models.Users{}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	cursor, err := entity.repo.Find(ctx, bson.M{"online": true})
	if err != nil {
		logrus.Print(err)
		return []models.Users{}, 400, err
	}
	for cursor.Next(ctx) {
		var user models.Users
		err = cursor.Decode(&user)
		if err != nil {
			logrus.Print(err)
		}
		usersList = append(usersList, user)
	}
	return usersList, http.StatusOK, nil

}

func (entity *usersEntity) GetUserByID(username string) (*models.Users, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var user models.Users
	err := entity.repo.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	//log.Fatal("username:" + username)
	if err != nil {
		logrus.Print(err)
		return nil, 400, err
	}

	return &user, http.StatusOK, nil
}

func (entity *usersEntity) Register(userForm form.Users) (*models.Users, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	//fmt.Printf("%+s\n", userForm.Username)
	user := models.Users{
		Username: userForm.Username,
		Online:   false,
	}
	//log.Fatal(userForm.Username)
	found, _, _ := entity.GetUserByID(user.Username)
	//fmt.Printf("Username for getID: %+s\n", user.Username)
	if found != nil {
		return nil, http.StatusBadRequest, errors.New("Username is already used")
	}
	_, err := entity.repo.InsertOne(ctx, user)

	if err != nil {
		logrus.Print(err)
		return nil, 400, err
	}

	return &user, http.StatusOK, nil
}

func (entity *usersEntity) GetAll() ([]models.Users, int, error) {
	usersList := []models.Users{}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	cursor, err := entity.repo.Find(ctx, bson.M{})

	if err != nil {
		logrus.Print(err)
		return []models.Users{}, 400, err
	}

	for cursor.Next(ctx) {
		var user models.Users
		err = cursor.Decode(&user)
		if err != nil {
			logrus.Print(err)
		}
		usersList = append(usersList, user)
	}
	return usersList, http.StatusOK, nil
}
