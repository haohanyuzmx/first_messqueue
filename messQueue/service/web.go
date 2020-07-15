package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"messQueue/model"
	"messQueue/mq"
)

func WhatOrder(ctx *gin.Context)  {
	var o model.Order
	err:=ctx.BindJSON(&o)
	fmt.Println(o)
	if err != nil {
		ctx.JSON(200,gin.H{
			"code":1,
			"mess":"输入有误",
		})
		return
	}
	mq.Joinmg(o)
	ctx.JSON(200,gin.H{
		"code":200,
		"mess":"处理完成",
	})
}