package http

import (
	"fmt"
	"go-core/util"
	"testing"
)

func TestPostWithJson(t *testing.T) {
	paras:= struct {
		Id int `json:"id"`
		UserName string `json:"username"`
	}{
		1234,
		"jinyidong",
	}

	parasAddr:=&paras

	body:=util.StructToJson(parasAddr)

	if body=="" {
		t.Error("反序列化失败！")
		return
	}
	bytes,err:=PostWithJson("http://www.tmall.com/test.do",body)

	if err!=nil {
		t.Error(err)
		return
	}

	fmt.Println(string(bytes))
}
