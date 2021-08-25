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
	"go.mongodb.org/mongo-driver/mongo/options"
)

var HistoryEntity HistoryDetails

type historyEntity struct {
	resource *db.Resource
	repo     *mongo.Collection
}

type HistoryDetails interface {
	History(historyForm form.History) (*models.History, int, error)
	Lastmatch(username string) ([]models.History, int, error)
}

func NewRankingEntity(resource *db.Resource) HistoryDetails {
	chalRepo := resource.DB.Collection("history")
	HistoryEntity = &historyEntity{resource: resource, repo: chalRepo}

	return HistoryEntity
}

func (entity *historyEntity) Lastmatch(username string) ([]models.History, int, error) {
	history := []models.History{}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"date", -1}})

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	cursor, err := entity.repo.Find(ctx, bson.M{}, findOptions.SetLimit(1))

	if err != nil {
		logrus.Print(err)
		return nil, 400, err
	}

	for cursor.Next(ctx) {
		var user models.History
		err = cursor.Decode(&user)
		if err != nil {
			logrus.Print(err)
		}
		history = append(history, user)
	}
	return history, http.StatusOK, nil
}

func (cEntity *historyEntity) History(historyForm form.History) (*models.History, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	history := models.History{
		Id:                historyForm.Id,
		Date:              historyForm.Date,
		User:              historyForm.User,
		Challenger:        historyForm.Challenger,
		User_win:          historyForm.User_win,
		User_lose:         historyForm.User_lose,
		Challenger_win:    historyForm.Challenger_win,
		Challenger_lose:   historyForm.Challenger_lose,
		Action_user:       historyForm.Action_user,
		Action_challenger: historyForm.Action_challenger,
	}

	_, err3 := cEntity.repo.InsertOne(ctx, history)
	if err3 != nil {
		logrus.Print(err3)
		return nil, 400, err3
	}

	return &history, http.StatusOK, nil
}
