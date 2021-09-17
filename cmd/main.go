package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/salihkemaloglu/btc-billionaire/pkg/handler"
)

func main() {

	// create new cache
	c := cache.New(24*30*time.Hour, 24*30*time.Hour)

	r := gin.Default()
	h := handler.NewHandler(c)
	handler.DefaultSet(c)

	health := r.Group("/wallet")
	{
		health.POST("history", h.GetTransactions)
		health.POST("", h.InsertTransaction)
	}

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
