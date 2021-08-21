package api

import (
	"golangBackend/db"
	"golangBackend/form"
	"golangBackend/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApplyUserAPI(app *gin.RouterGroup, resource *db.Resource) {
	userEntity := repository.NewUserEntity(resource)
	authRoute := app.Group("")
	//authRoute.POST("/login", login(userEntity))
	authRoute.POST("/signup", signUp(userEntity))
	authRoute.GET("/getallusers", GetAllUsers(userEntity))

}

func signUp(userEntity repository.UserDetails) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		//fmt.Printf("%+s\n", ctx)
		userRequest := form.Users{}

		if err := ctx.Bind(&userRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		user, code, err := userEntity.Register(userRequest)
		//fmt.Printf("%+v\n", user.Username)
		response := map[string]interface{}{
			"user":  user,
			"error": err.Error(),
		}
		ctx.JSON(code, response)
	}

}

func GetAllUsers(userEntity repository.UserDetails) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		list, code, err := userEntity.GetAll()
		response := map[string]interface{}{
			"username": list,
			"error":    err.Error(),
		}
		ctx.JSON(code, response)
	}
}
