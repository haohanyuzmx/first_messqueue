package model

import (
	"fmt"
	"testing"
)

func TestDealOrder(t *testing.T) {
	//o:=Order{
	//	Id:        1,
	//	WhatTime:  "1.20",
	//	WhoId:     1,
	//	FilmId:    1,
	//	Hashandle: "",
	//}
	//DealOrder(o)
	f:=Film{
		Id:       1,
		Name:     "",
		WhatTime: "",
		Where:    "",
		Num:      0,
	}
	err:=DB.Where(&f).First(&f).Error
	f.Num--
	err=DB.Model(&f).Update("num",f.Num).Error
	fmt.Println(err)
}
