package repository

import (
	"context"
	"golangBackend/db"
	"golangBackend/form"
	"golangBackend/models"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ChallengeEntity ChallengeDetails

type challengeEntity struct {
	resource *db.Resource
	repo     *mongo.Collection
}

type ChallengeDetails interface {
	//History(challengeForm form.Challenge) (*models.Users, int, error)
	Challenging(challengeForm form.Challenge) (*models.Challenge, int, error)
	GetChallengeByID(username string, challenger string) (*models.Challenge, int, error)
}

func NewChallengeEntity(resource *db.Resource) ChallengeDetails {
	chalRepo := resource.DB.Collection("challengeDetails")
	ChallengeEntity = &challengeEntity{resource: resource, repo: chalRepo}

	return ChallengeEntity
}

func (entity *challengeEntity) GetChallengeByID(username string, challenger string) (*models.Challenge, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var user models.Challenge
	err := entity.repo.FindOne(ctx, bson.M{"id": username + "and" + challenger}).Decode(&user)
	user.Type = 1
	if err != nil {
		err2 := entity.repo.FindOne(ctx, bson.M{"id": challenger + "and" + username}).Decode(&user)
		user.Type = 2
		if err2 != nil {
			logrus.Print(err)
			return nil, 400, err
		}
	}
	//log.Fatal(user)

	return &user, http.StatusOK, nil
}

func (cEntity *challengeEntity) Challenging(challengeForm form.Challenge) (*models.Challenge, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	//Total score
	var user1_Totalwin = 0
	var user1_Totallose = 0
	var user2_Totalwin = 0
	var user2_Totallose = 0

	//This match score
	//var user1_win = 0
	//var user1_lose = 0
	//var user2_win = 0
	//var user2_lose = 0
	var users models.Challenge
	//log.Fatal(challengeForm.User1.Total_lose)
	found, _, _ := cEntity.GetChallengeByID(challengeForm.Username, challengeForm.Challenger)
	//log.Fatal(found)
	if found != nil {

		//users.Users1 = found.Users1
		//users.Users2 = found.Users2
		users.Id = found.Id
		users.User = found.User
		users.Challenger = found.Challenger
		users.User_win = found.User_win
		users.User_lose = found.User_lose
		users.Challenger_win = found.Challenger_win
		users.Challenger_lose = found.Challenger_lose
		users.Type = found.Type

		if found.Type == 1 {

			users.Action_user1 = challengeForm.Action_user
			users.Action_challenger = challengeForm.Action_challenger

		} else if found.Type == 2 {

			users.Action_user1 = challengeForm.Action_challenger
			users.Action_challenger = challengeForm.Action_user
		}
	} else {
		users.User = challengeForm.Username
		users.Challenger = challengeForm.Challenger
		users.User_win = 0
		users.User_lose = 0
		users.Challenger_win = 0
		users.Challenger_lose = 0
		users.Type = 0
		//users.Users1 = models.Users(challengeForm.User1)
		//users.Users2 = models.Users(challengeForm.User2)
		users.Action_user1 = challengeForm.Action_user
		users.Action_challenger = challengeForm.Action_challenger
	}
	user1_Totalwin = users.User_win
	user1_Totallose = users.User_lose
	user2_Totalwin = users.Challenger_win
	user2_Totallose = users.Challenger_lose

	//log.Fatal(users)

	// Rock
	if users.Action_user1 == "Rock" && users.Action_challenger == "Paper" {
		users.Winner = 2
		user1_Totallose++
		user2_Totalwin++

	} else if users.Action_user1 == "Rock" && users.Action_challenger == "Scissors" {
		users.Winner = 1
		user1_Totalwin++
		user2_Totallose++
	} else if users.Action_user1 == "Rock" && users.Action_challenger == "Rock" {
	}

	// end Rock

	// Paper
	if users.Action_user1 == "Paper" && users.Action_challenger == "Rock" {
		users.Winner = 1
		user1_Totalwin++
		user2_Totallose++
	} else if users.Action_user1 == "Paper" && users.Action_challenger == "Scissors" {
		users.Winner = 2
		user1_Totallose++
		user2_Totalwin++

	} else if users.Action_user1 == "Paper" && users.Action_challenger == "Paper" {

	}

	// end Paper

	// Scissors
	if users.Action_user1 == "Scissors" && users.Action_challenger == "Rock" {
		users.Winner = 2
		user1_Totallose++
		user2_Totalwin++
	} else if users.Action_user1 == "Scissors" && users.Action_challenger == "Paper" {
		users.Winner = 1
		user1_Totalwin++
		user2_Totallose++
	} else if users.Action_user1 == "Scissors" && users.Action_challenger == "Scissors" {

	}

	// end Scissors
	users.User_win = user1_Totalwin
	users.User_lose = user1_Totallose
	users.Challenger_win = user2_Totalwin
	users.Challenger_lose = user2_Totallose

	//log.Fatal(users)

	//fmt.Printf("%+s\n", user1_Totalwin)
	//fmt.Printf("User1: %+s\n", user1_Totallose)
	//fmt.Printf("User2: %+s\n", user2_Totalwin)
	//fmt.Printf("User2: %+s\n", user2_Totallose)

	challenger := models.Challenge{
		Id: challengeForm.Username + "and" + challengeForm.Challenger,
		//	Users1:            users.Users1,
		//	Users2:            users.Users2,
		User:              users.User,
		Challenger:        users.Challenger,
		User_win:          user1_Totalwin,
		User_lose:         user1_Totallose,
		Challenger_win:    user2_Totalwin,
		Challenger_lose:   user2_Totallose,
		Action_user1:      challengeForm.Action_user,
		Action_challenger: challengeForm.Action_challenger,
		Type:              users.Type,
		Winner:            users.Winner,
	}

	/*history := models.History{
		Date:              time.Time{},
		Win:               "",
		Lose:              "",
		Action_user:       challengeForm.Action_user,
		Action_challenger: challengeForm.Action_challenger,
	}*/

	//fmt.Printf("%+s\n", found.Id)
	if found != nil {
		_, err3 := cEntity.repo.UpdateOne(ctx, bson.M{"id": users.Id},
			bson.M{"$set": bson.M{"user_win": users.User_win, "user_lose": users.User_lose, "challenger_win": users.Challenger_win, "challenger_lose": users.Challenger_lose}})
		//fmt.Printf("User2: %+s\n", test)
		if err3 != nil {
			logrus.Print(err3)
			return nil, 400, err3
		}
	} else {
		_, err3 := cEntity.repo.InsertOne(ctx, challenger)
		if err3 != nil {
			logrus.Print(err3)
			return nil, 400, err3
		}
	}
	//_, err := UserEntity.repo.UpdateOne(ctx, bson.M{"username": challengeForm.User1.Username}, bson.M{"$set": bson.M{"total_win": user1_Totalwin, "total_lose": user1_Totallose}})
	//_, err2 := cEntity.repo.UpdateOne(ctx, bson.M{"username": challengeForm.User2.Username}, bson.M{"$set": bson.M{"total_win": user2_Totalwin, "total_lose": user2_Totallose}})

	/*if err != nil {
		logrus.Print(err)
		return nil, 400, err
	}

	if err2 != nil {
		logrus.Print(err2)
		return nil, 400, err2
	}*/

	return &challenger, http.StatusOK, nil
}
