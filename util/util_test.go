package util

import (
	"fmt"
	"testing"
	"time"
)

func TestStructToJson(t *testing.T) {
	user:=struct {
		Name string
		Age int
		Birthday time.Time
	}{
		"jinyidong",
		12,
		time.Date(2019,10,23,0,0,0,0,time.Local),
	}

	userAddr:=&user

	userJson:=StructToJson(userAddr)

	if userJson=="" {
		t.Error("序列化失败！")
		return
	}
	fmt.Println(userJson)
}
