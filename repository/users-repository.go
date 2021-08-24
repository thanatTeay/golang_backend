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
	//History(challengeForm form.Challenge) (*models.Users, int, error)
	//Challenging(challengeForm form.Challenge) (*models.Challenge, int, error)
	UpdateStatus(userForm form.Users) (*models.Users, int, error)
	UpdateScore(rankForm form.Ranking) (*models.Users, int, error)
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

func (entity *usersEntity) UpdateScore(rankForm form.Ranking) (*models.Users, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var user_total_win = 0
	var user_total_lose = 0
	var challenger_total_win = 0
	var challenger_total_lose = 0

	found1, _, _ := entity.GetUserByID(rankForm.User)
	found2, _, _ := entity.GetUserByID(rankForm.Challenger)
	//fmt.Printf("Username for getID: %+s\n", user.Username)
	if found1 == nil {
		return nil, http.StatusBadRequest, errors.New("User not found")
	}
	if found2 == nil {
		return nil, http.StatusBadRequest, errors.New("Challenger not found")
	}
	user_total_win = found1.Total_win
	challenger_total_lose = found2.Total_lose
	challenger_total_win = found2.Total_win
	user_total_lose = found1.Total_lose

	userScore := models.Users{
		//Users1: models.Users(challengeForm.User1),
		//Users2: models.Users(challengeForm.User2),
	}
	//log.Fatal(found)
	if rankForm.Winner == 1 {

		user_total_win++

		challenger_total_lose++

	} else {

		challenger_total_win++

		user_total_lose++
	}

	_, err := entity.repo.UpdateOne(ctx, bson.M{"username": rankForm.User}, bson.M{"$set": bson.M{"total_win": user_total_win, "total_lose": user_total_lose}})
	_, err2 := entity.repo.UpdateOne(ctx, bson.M{"username": rankForm.Challenger}, bson.M{"$set": bson.M{"total_win": challenger_total_win, "total_lose": challenger_total_lose}})
	if err != nil {
		logrus.Print(err)
		return nil, 400, err
	}
	if err2 != nil {
		logrus.Print(err2)
		return nil, 400, err
	}
	return &userScore, http.StatusOK, nil
}

func (entity *usersEntity) UpdateStatus(userForm form.Users) (*models.Users, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	userStatus := models.Users{
		Username:    userForm.Username,
		Status_user: userForm.Status_user,
	}
	_, err := entity.repo.UpdateOne(ctx, bson.M{"username": userStatus.Username}, bson.M{"$set": bson.M{"status_user": userStatus.Status_user}})
	if err != nil {
		logrus.Print(err)
		return nil, 400, err
	}

	return &userStatus, http.StatusOK, nil

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
		Username:    userForm.Username,
		Status_user: "neutral",
		Online:      false,
		Total_win:   0,
		Total_lose:  0,
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
	//usersList := []models.Users{}
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
