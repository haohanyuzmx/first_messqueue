package main

import (
	"github.com/gin-gonic/gin"
	"messQueue/mq"
	"messQueue/service"
)

func main()  {
	go mq.Dealmq()
	r:=gin.Default()
	r.POST("/order",service.WhatOrder)
	r.Run()
}