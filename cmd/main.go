package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/viktor-viktor/shard/internal"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// initializing server dependencies.
	url := fmt.Sprintf("mongodb://%v:%v@%v", os.Getenv("MONGO_USER"), os.Getenv("MONGO_PWD"), os.Getenv("MONGO_ADDR"))
	db, err := internal.NewMongoDBDAL(context.TODO(), url, os.Getenv("MONGO_DB_NAME"))
	if err != nil {
		panic(err)
	}
	p := internal.StartPools(db)

	// starting server
	router := gin.Default()
	router.POST("/events", internal.BuildEventsController(p))
	router.GET("/workers", internal.BuildWorkersController(db))
	router.GET("/workers/:id", internal.BuildSingleWorkerController(db))
	srv := &http.Server{
		Addr:    ":9753",
		Handler: router,
	}
	go func() {
		srv.ListenAndServe()
	}()

	// graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Println("Error when shutting down server. Error: ", err)
	}

	p.ShutdownPools()
}
