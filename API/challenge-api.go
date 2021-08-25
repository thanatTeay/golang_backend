package api

import (
	"fmt"
	"golangBackend/db"
	"golangBackend/form"
	"golangBackend/models"
	"golangBackend/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ApplyChallengeAPI(app *gin.RouterGroup, resource *db.Resource) {
	challengeEntity := repository.NewChallengeEntity(resource)
	userEn := repository.NewUserEntity(resource)
	hisEn := repository.NewRankingEntity(resource)
	authRouteChallenge := app.Group("/challenge")
	authRouteChallenge.POST("/fighting", Challenge(challengeEntity, userEn, hisEn))
	authRouteChallenge.POST("/status", StatusChallenger(challengeEntity, userEn, hisEn))
}

func StatusChallenger(challengeEntity repository.ChallengeDetails, userEn repository.UserDetails, hisTn repository.HistoryDetails) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var cRequest form.Challenge
		err := ctx.BindJSON(&cRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		challengeDetails, code, err := challengeEntity.GetChallengeByID(cRequest.Username, cRequest.Challenger)

		var user *models.Users
		var challenger *models.Users
		var lastM []models.History
		if challengeDetails != nil {
			user, code, err = userEn.GetUserByID(challengeDetails.User)
			challenger, code, err = userEn.GetUserByID(challengeDetails.Challenger)
			lastM, code, err = hisTn.Lastmatch(challengeDetails.Id)
		}

		if challengeDetails == nil || user == nil || challenger == nil {

			response := map[string]interface{}{
				"err": err.Error(),
			}
			ctx.JSON(code, response)
		} else {
			response := map[string]interface{}{
				"data":       challengeDetails,
				"user":       user,
				"challenger": challenger,
				"last_match": lastM,
				//"score": updateScore,
				//"error": "err",
			}
			ctx.JSON(code, response)
		}

	}
}

func Challenge(challengeEntity repository.ChallengeDetails, userEn repository.UserDetails, hisEn repository.HistoryDetails) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var cRequest form.Challenge
		err := ctx.BindJSON(&cRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		challengeDetails, code, err := challengeEntity.Challenging(cRequest)
		RankingRequest := form.Ranking{
			User:            challengeDetails.User,
			Challenger:      challengeDetails.Challenger,
			User_win:        challengeDetails.User_win,
			User_lose:       challengeDetails.User_lose,
			Challenger_win:  challengeDetails.Challenger_win,
			Challenger_lose: challengeDetails.Challenger_lose,
			Winner:          challengeDetails.Winner,
		}

		History := form.History{
			Id:                challengeDetails.Id,
			Date:              time.Now().String(),
			User:              challengeDetails.User,
			Challenger:        challengeDetails.Challenger,
			User_win:          challengeDetails.User_win,
			User_lose:         challengeDetails.User_lose,
			Challenger_win:    challengeDetails.Challenger_win,
			Challenger_lose:   challengeDetails.Challenger_lose,
			Action_user:       challengeDetails.Action_user1,
			Action_challenger: challengeDetails.Action_challenger,
		}

		history, code, err := hisEn.History(History)

		updateScore, code, err := userEn.UpdateScore(RankingRequest)
		fmt.Printf("%+v\n", updateScore)
		if challengeDetails == nil || updateScore == nil || history == nil {
			response := map[string]interface{}{
				"err": err.Error(),
			}
			ctx.JSON(code, response)
		} else {
			response := map[string]interface{}{
				"data": challengeDetails,
			}
			ctx.JSON(code, response)
		}

	}
}
