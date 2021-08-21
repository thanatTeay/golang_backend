package main

import (
	api "golangBackend/API"
	"golangBackend/db"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Routes struct {
}

func main() {

	/*conn, err := connectDB()
	if err != nil {
		logrus.Error(err)
	}*/
	router := gin.Default()
	//router.Use(dbMiddleware(conn))

	usersGroup := router.Group("api/users")
	resource, err := db.InitResourc()
	if err != nil {
		logrus.Error(err)
	}
	api.ApplyUserAPI(usersGroup, resource)

	/*{
		usersGroup.POST("register", routes.UserResister)
		//usersGroup.POST("register", routes.UserResister)
	}*/

	router.Run(":8080")
}
