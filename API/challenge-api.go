package api

import (
	"fmt"
	"golangBackend/db"
	"golangBackend/form"
	"golangBackend/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApplyChallengeAPI(app *gin.RouterGroup, resource *db.Resource) {
	challengeEntity := repository.NewChallengeEntity(resource)
	userEn := repository.NewUserEntity(resource)
	authRouteChallenge := app.Group("/challenge")
	authRouteChallenge.POST("/fighting", Challenge(challengeEntity, userEn))
}

func Challenge(challengeEntity repository.ChallengeDetails, userEn repository.UserDetails) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var cRequest form.Challenge
		err := ctx.BindJSON(&cRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		/*users, _, err := userEn.GetUserByID(cRequest.Username)
		if err != nil {
			log.Fatal(err)
		}*/

		/*challenger, _, err := userEn.GetUserByID(cRequest.Challenger)
		//cRequest.User1 = append(cRequest.User1,form.Users(*users))
		cRequest.User1 = form.Users(*users)
		cRequest.User2 = form.Users(*challenger)
		//fmt.Printf("%+v\n", cRequest)*/

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
		//fmt.Printf("%+v\n", challengeDetails)
		//cRequest.User1 = form.Users(challengeDetails.Users1)
		//cRequest.User2 = form.Users(challengeDetails.Users2)
		updateScore, code, err := userEn.UpdateScore(RankingRequest)
		//history, code, err := userEntity.History(cRequest)
		fmt.Printf("%+v\n", updateScore)
		response := map[string]interface{}{
			"data": challengeDetails,
			//"score": updateScore,
			//"error": "err",
		}
		ctx.JSON(code, response)

	}
}