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
	Challenging(challengeForm form.Challenge) (*models.Challenge, int, error)
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

func (entity *usersEntity) Challenging(challengeForm form.Challenge) (*models.Challenge, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	//fmt.Printf("%+s\n", userForm.Username)
	foundUser1, _, _ := entity.GetUserByID(challengeForm.Username)
	if foundUser1 == nil {
		return nil, http.StatusBadRequest, errors.New("Not found Username")
	}
	foundUser2, _, _ := entity.GetUserByID(challengeForm.Challenger)
	if foundUser2 == nil {
		return nil, http.StatusBadRequest, errors.New("Not found Username")
	}
	var user1_win = 0
	var user1_lose = 0
	var user2_win = 0
	var user2_lose = 0

	// Rock
	if challengeForm.Action_user == "Rock" && challengeForm.Action_challenger == "Paper" {
		user1_lose = foundUser1.Total_lose
		user1_lose++

	} else if challengeForm.Action_user == "Rock" && challengeForm.Action_challenger == "Scissors" {
		user1_win = foundUser1.Total_lose
		user1_win++
	} else if challengeForm.Action_user == "Rock" && challengeForm.Action_challenger == "Rock" {

	} else if challengeForm.Action_challenger == "Rock" && challengeForm.Action_user == "Paper" {
		user2_lose = foundUser2.Total_lose
		user2_lose++
	} else if challengeForm.Action_challenger == "Rock" && challengeForm.Action_user == "Scissors" {
		user2_win = foundUser2.Total_win
		user2_win++
	}
	// end Rock

	// Paper
	if challengeForm.Action_user == "Paper" && challengeForm.Action_challenger == "Rock" {
		user1_win = foundUser1.Total_lose
		user1_win++
	} else if challengeForm.Action_user == "Paper" && challengeForm.Action_challenger == "Scissors" {
		user1_lose = foundUser1.Total_lose
		user1_lose++

	} else if challengeForm.Action_user == "Paper" && challengeForm.Action_challenger == "Paper" {

	} else if challengeForm.Action_challenger == "Paper" && challengeForm.Action_user == "Rock" {
		user2_win = foundUser2.Total_win
		user2_win++
	} else if challengeForm.Action_challenger == "Paper" && challengeForm.Action_user == "Scissors" {
		user2_lose = foundUser2.Total_lose
		user2_lose++

	}

	// end Paper

	// Scissors
	if challengeForm.Action_user == "Scissors" && challengeForm.Action_challenger == "Rock" {
		user1_lose = foundUser1.Total_lose
		user1_lose++
	} else if challengeForm.Action_user == "Scissors" && challengeForm.Action_challenger == "Paper" {
		user1_win = foundUser1.Total_lose
		user1_win++
	} else if challengeForm.Action_user == "Scissors" && challengeForm.Action_challenger == "Scissors" {

	} else if challengeForm.Action_challenger == "Scissors" && challengeForm.Action_user == "Rock" {
		user2_lose = foundUser2.Total_lose
		user2_lose++

	} else if challengeForm.Action_challenger == "Scissors" && challengeForm.Action_user == "Paper" {
		user2_win = foundUser2.Total_win
		user2_win++
	}

	// end Scissors

	user1 := models.Users{
		Username:    foundUser1.Username,
		Online:      foundUser1.Online,
		Status_user: "challenging",
		Total_win:   user1_win,
		Total_lose:  user1_lose,
	}

	user2 := models.Users{
		Username:    foundUser2.Username,
		Online:      foundUser2.Online,
		Status_user: "being challenged",
		Total_win:   user2_win,
		Total_lose:  user2_lose,
	}

	challenger := models.Challenge{
		Id:     challengeForm.Username + "and" + challengeForm.Challenger,
		Users1: user1,
		Users2: user2,

		Action_user1: challengeForm.Action_user,
		Action_user2: challengeForm.Action_challenger,
	}

	/*history := models.History{
		Date:              time.Time{},
		Win:               "",
		Lose:              "",
		Action_user:       challengeForm.Action_user,
		Action_challenger: challengeForm.Action_challenger,
	}*/

	_, err := entity.repo.UpdateOne(ctx, bson.M{"username": foundUser1.Username}, bson.M{"$set": bson.M{"status_user": user1.Status_user, "total_win": user1.Total_win, "total_lose": user1.Total_lose}})
	_, err2 := entity.repo.UpdateOne(ctx, bson.M{"username": foundUser2.Username}, bson.M{"$set": bson.M{"status_user": user2.Status_user, "total_win": user2.Total_win, "total_lose": user2.Total_lose}})
	_, err3 := entity.repo.InsertOne(ctx, challenger)

	if err != nil {
		logrus.Print(err)
		return nil, 400, err
	}

	if err2 != nil {
		logrus.Print(err2)
		return nil, 400, err2
	}
	if err3 != nil {
		logrus.Print(err3)
		return nil, 400, err3
	}

	return &challenger, http.StatusOK, nil
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
