package config

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"sync"
	"time"
)

type ZK struct {
	conn    *zk.Conn
	servers []string
	ch      <-chan zk.Event
	watcher func(event zk.Event)
}

func newZK(servers []string, watcher func(zk.Event)) *ZK {
	var zookeeper = &ZK{}
	zookeeper.servers = servers
	zookeeper.watcher = watcher

	var wg = &sync.WaitGroup{}
	wg.Add(1) //Add用来设置等待线程数量

	go func() {
		var ech <-chan zk.Event
		var err error
		var c *zk.Conn
		if zookeeper.conn == nil {
			c, ech, err = zk.Connect(servers, time.Second) //*10)
			zookeeper.conn = c
		}

		if err != nil {
			panic(err)
		}

		go func() {
			for {
				select {
				case ch := <-ech:
					{
						switch ch.State {
						case zk.StateConnecting:
							{
								fmt.Println("StateConnecting")
							}
						case zk.StateConnected: //若链接后出现断连现象 重连时会报异常 因为此时同步信号量已为0
							{
								wg.Done()
								fmt.Println("StateConnected")
							}
						case zk.StateExpired:
							{
								fmt.Println("StateExpired")
							}
						case zk.StateDisconnected:
							{
								wg.Add(1)
								fmt.Println("StateDisconnected")
							}
						}

						watcher(ch)

					}
				}
			}
		}()

	}()

	wg.Wait()
	return zookeeper
}

func (zookeeper *ZK) Create(path string, data []byte, version int32) error {

	err := zookeeper.do(func(conn *zk.Conn) error {
		_, err := conn.Create(path, data, 0, zk.WorldACL(zk.PermAll))
		return err
	})

	return err
}

func (zookeeper *ZK) Set(path string, data []byte, version int32) error {

	err := zookeeper.do(func(conn *zk.Conn) error {
		_, err := conn.Set(path, data, version)
		return err
	})

	return err
}

func (zookeeper *ZK) Delete(path string, version int32) error {

	err := zookeeper.do(func(conn *zk.Conn) error {

		err := conn.Delete(path, version)
		return err
	})

	return err
}

func (zookeeper *ZK) Get(path string) (string, error) {
	var res []byte
	err := zookeeper.do(func(conn *zk.Conn) error {
		r, _, err := conn.Get(path)
		res = r
		return err
	})

	return string(res), err
}

func (zookeeper *ZK) GetW(path string) (string, <-chan zk.Event, error) {
	var res []byte
	var c <-chan zk.Event
	err := zookeeper.do(func(conn *zk.Conn) error {
		var err error
		res, _, c, err = conn.GetW(path)
		return err
	})

	return string(res), c, err
}

func (zookeeper *ZK) GetChildren(path string) ([]string, error) {

	var res []string

	err := zookeeper.do(func(conn *zk.Conn) error {
		r, _, err := conn.Children(path)
		res = r
		return err
	})

	return res, err
}

func (zookeeper *ZK) GetChildrenW(path string) ([]string, <-chan zk.Event, error) {

	var res []string
	var c <-chan zk.Event
	err := zookeeper.do(func(conn *zk.Conn) error {
		var err error
		res, _, c, err = conn.ChildrenW(path)
		return err
	})

	return res, c, err
}

func (zookeeper *ZK) do(fn func(conn *zk.Conn) error) error {
	conn := zookeeper.conn
	if conn != nil {
		err := fn(conn)
		return err
	}

	oldZK := zookeeper

	DefaultZK := newZK(zookeeper.servers, zookeeper.watcher)

	oldZK.conn.Close()

	err := fn(DefaultZK.conn)

	return err
}
