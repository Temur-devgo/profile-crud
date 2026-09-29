package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "sarkor/docs"
	"sarkor/internal/config"
	"sarkor/internal/handler"
	"sarkor/internal/middlware"
	"sarkor/internal/repository"
	"sarkor/internal/service"
	"sarkor/migrations"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf := config.GetConfigs()

	sqLiteConn, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		panic(fmt.Errorf("DB connection err: %w", err).Error())
	}
	defer func(sqLiteConn *sql.DB) {
		if err = sqLiteConn.Close(); err != nil {
			panic(fmt.Errorf("DB connection close err: %w", err).Error())
		}
	}(sqLiteConn)

	if err = migrations.Run(ctx, sqLiteConn); err != nil {
		panic(fmt.Errorf("migrations err: %w", err).Error())
	}

	repo := repository.NewUsersRepo(sqLiteConn)
	svc := service.NewUserSvc(repo)
	userHandler := handler.NewUserHandler(*svc)

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	baseGroup := router.Group("/")
	userHandler.AuthFreeRoutes(baseGroup)

	baseGroup.Use(middlware.ValidateToken)
	userHandler.AuthRoutes(baseGroup)

	if err = router.Run(fmt.Sprintf("%s:%v", conf.ServerHost, conf.ServerPort)); err != nil {
		fmt.Println(err.Error())
	}
}
