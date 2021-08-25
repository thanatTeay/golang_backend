package api

import (
	"fmt"
	"golangBackend/db"
	"golangBackend/form"
	"golangBackend/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApplyUserAPI(app *gin.RouterGroup, resource *db.Resource) {
	userEntity := repository.NewUserEntity(resource)
	authRoute := app.Group("")

	authRoute.POST("/signup", signUp(userEntity))

	authRouteUser := app.Group("/users")
	authRouteUser.GET("/getall", GetAllUsers(userEntity))
	authRouteUser.GET("/getonline", GetOnlineUsers(userEntity))
	authRoute.POST("/updateStatus", UpdateStatusUser(userEntity))

	authRouteRanking := app.Group("/ranking")
	authRouteRanking.GET("", Ranking(userEntity))

}

func Ranking(userEntity repository.UserDetails) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		list, code, _ := userEntity.Ranking()
		response := map[string]interface{}{
			"username": list,
			//"error":    err.Error(),
		}
		ctx.JSON(code, response)
	}
}

func UpdateStatusUser(userEntity repository.UserDetails) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var userRequest form.Users
		err := ctx.BindJSON(&userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		user, code, err := userEntity.UpdateStatus(userRequest)

		if user == nil {
			response := map[string]interface{}{
				"err": err.Error(),
			}
			ctx.JSON(code, response)
		} else {
			response := map[string]interface{}{
				"username": user.Username,
				"status":   user.Status_user,
				"message":  "Successfully updated status",
			}
			ctx.JSON(code, response)
		}

	}
}

func signUp(userEntity repository.UserDetails) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		//fmt.Printf("%+s\n", ctx)
		var userRequest form.Users

		err := ctx.BindJSON(&userRequest)
		fmt.Printf("CheckUsernameRequest:  %+v\n", userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		user, code, err := userEntity.Register(userRequest)
		fmt.Printf("%+v\n", user)
		response := map[string]interface{}{
			"username": user,
			//"error":    err.Error(),
		}
		fmt.Printf("%+v\n", "test print before ctx.JSON in signUp")
		ctx.JSON(code, response)
	}

}

func GetAllUsers(userEntity repository.UserDetails) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		list, code, _ := userEntity.GetAll()
		response := map[string]interface{}{
			"username": list,
			//"error":    err.Error(),
		}
		ctx.JSON(code, response)
	}
}

func GetOnlineUsers(userEntity repository.UserDetails) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		list, code, _ := userEntity.OnlineUsers()
		response := map[string]interface{}{
			"username": list,
			//"error":    err.Error(),
		}
		ctx.JSON(code, response)
	}
}
