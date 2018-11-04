package config

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"testing"
)

func TestZK_newZK(t *testing.T) {
	var zookeeper = newZK(BasicConfig.Servers(), func(event zk.Event) {})

	fmt.Println(zookeeper.conn.State().String())
}

func TestZK_GetChildren(t *testing.T) {
	var zookeeper = newZK(BasicConfig.Servers(), func(event zk.Event) {

	})

	var res, err = zookeeper.GetChildren("/")

	if err != nil {
		fmt.Errorf(err.Error())
	}

	for _, v := range res {
		fmt.Println(v)
	}

}

func BenchmarkZK_GetChildren(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var zookeeper = newZK(BasicConfig.Servers(), func(event zk.Event) {

		})

		var _, err = zookeeper.GetChildren("/")

		if err != nil {
			b.Fatalf(err.Error())
		}
	}
}

func TestZK_Set(t *testing.T) {
	var zookeeper = newZK(BasicConfig.Servers(), func(event zk.Event) {

	})

	//err := zk.Create("/service2/1231",[]byte("hello"),0)
	//
	//if err != nil {
	//	fmt.Println("create",err)
	//}

	err := zookeeper.Set("/service2/product-search-admin", []byte("{\"Endpoints\":[{\"Namespace\":\"dev\",\"Ip\":\"172.20.1.240\",\"Port\":8080,\"Weight\":1}],\"Name\":\"product-search-admin\",\"Version\":\"af368d4e5406d096542ae4fc3594faf4\"}"), -1)

	if err != nil {
		fmt.Println(err)
	}

	//s,err:=zk.Get("/service2/1231")
	//
	//if err!= nil{
	//	fmt.Println("get",err)
	//}
	//fmt.Println(s)
	//
	//zk.Delete("/service2/1231",0)
	//
	//zk.Delete("/config/123",0)
}

func TestZK_Create(t *testing.T) {
	var zookeeper = newZK(BasicConfig.Servers(), func(event zk.Event) {

	})

	err := zookeeper.Create("/service2/1231", []byte("hello"), 0)

	if err != nil {
		fmt.Println("create", err)
	}

	//err = zk.Set("/config/123",[]byte("333"),1)
	//
	//if err!= nil{
	//	fmt.Println(err)
	//}

	s, err := zookeeper.Get("/service2/1231")

	if err != nil {
		fmt.Println("get", err)
	}
	fmt.Println(s)

	zookeeper.Delete("/service2/1231", -1)
	//
	//zk.Delete("/config/123",0)
}
