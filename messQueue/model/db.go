package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)
type Film struct {
	Id       int    `gorm:"primary_key;auto_increment"json:"id"`
	Name     string `json:"name"`
	WhatTime string `json:"what_time"`
	Where    string `json:"where"`
	Num      int    `json:"num"`
}
type Order struct {
	Id        int    `gorm:"primary_key;auto_increment"json:"id"`
	WhatTime  string `json:"what_time"`
	WhoId     int    `json:"who_id"`
	FilmId    int    `json:"film_id"`
	Hashandle string `json:"hashandle"`
}
var DB *gorm.DB
func init()  {
	my,err:=gorm.Open("mysql","root:@tcp(127.0.0.1:3306)/test?parseTime=true&charset=utf8&loc=Local")
	if err != nil {
		log.Println(err)
		return
	}
	DB=my
	if !DB.HasTable(&Film{}) {
		DB.CreateTable(&Film{})
	}
	if !DB.HasTable(&Order{}) {
		DB.CreateTable(&Order{})
	}
}

func DealOrder(order Order)  {
	var f Film
	DB.Table("films").Where("id=?",order.FilmId).First(&f)
	f.Num--
	if f.Num<0 {
		log.Println("票没了")
		order.Hashandle="false"
		DB.Create(&order)
		return
	}
	log.Println(f.Num)
	order.Hashandle="true"
	DB.Model(&f).Update("num",f.Num)
	DB.Create(&order)
}
