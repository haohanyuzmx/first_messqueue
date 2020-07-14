package main

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"messQueue"
	"strconv"
	"time"
)

var st = make(chan int,1)
type film struct {
	id       int    `gorm:"primary_key;auto_increment"json:"id"`
	name     string `json:"name"`
	whatTime string `json:"what_time"`
	where    string `json:"where"`
	num      int    `json:"num"`
}
type order struct {
	id       int    `gorm:"primary_key;auto_increment"json:"id"`
	whatTime string `json:"what_time"`
	whoid    int    `json:"whoid"`
	filmid   int    `json:"filmid"`
}

// 实现发送者

func (o *order) MsgContent() []byte {
	afilm, err := json.Marshal(o)
	if err != nil {
		log.Print(err)
		return nil
	}
	return afilm
}

// 实现接收者

func (t *order) Consumer(dataByte []byte) error {
	err := json.Unmarshal(dataByte, t)
	if err != nil {
		return err
	}
	var oldf film
	messQueue.DB.Create(t)
	messQueue.DB.Table("films").Where("id=?",t.filmid).First(&oldf)
	oldf.num --
	if oldf.num <= 0 {
		err = errors.New("票空了")
		return err
	}
	err=messQueue.DB.Table("films").Where("id=?", t.filmid).Update("num", oldf.num).Error
	return err
}


func main() {
	messQueue.MQ.RegisterReceiver(&order{})
	r:=gin.Default()
	r.POST("order",getorder)
	go start()
	r.Run()
}
func start()  {
	for  {
		messQueue.MQ.Start()
		time.Sleep(10*time.Minute)
	}
}
func getorder(ctx *gin.Context)  {
	fid:=ctx.PostForm("fileid")
	fidint,err:=strconv.Atoi(fid)
	if err != nil {
		log.Println(err)
		return
	}
	oid:=ctx.PostForm("whoid")
	oidint,err:=strconv.Atoi(oid)
	if err != nil {
		log.Println(err)
		return
	}
	timeTemplate1 := "2006-01-02 15:04:05"
	od:=&order{
		filmid: fidint,
		whoid: oidint,
		whatTime: time.Now().Format(timeTemplate1),
	}
	messQueue.MQ.RegisterProducer(od)
	ctx.JSON(200,gin.H{
		"mess":"订票成功",
	})
	st<-0
}