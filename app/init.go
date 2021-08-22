package app

import (
	api "golangBackend/API"
	"golangBackend/db"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Routes struct {
}

func (app Routes) StartMain() {
	router := gin.Default()
	router.Use(gin.Logger())
	usersGroup := router.Group("api")
	resource, err := db.InitResourc()
	if err != nil {
		logrus.Error(err)
	}
	api.ApplyUserAPI(usersGroup, resource)

	router.Run(":8080")
}
