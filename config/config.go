package config

import (
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
	"sync"
)

var Paths = "config,connection"

type Manager struct {
	configMap map[string]string
}

var manager *Manager
var once sync.Once
var zoo *ZK

func NewManager() *Manager {

	once.Do(func() {
		manager = &Manager{}
		manager.configMap = make(map[string]string)

		zoo = newZK(strings.Split(BasicConfig.Zookeeper, ","), func(event zk.Event) {
			switch event.State {
			case zk.StateExpired:
				{
					fmt.Println("StateExpired")
					setAll()
					fmt.Println("setAll /")
				}
			}
			switch event.Type {
			case zk.EventNodeDataChanged:
				{
					manager.update(event.Path)

					fmt.Printf("EventNodeDataChanged: %s\n", event.Path)
				}
			case zk.EventNodeDeleted:
				{
					manager.delete(event.Path)
					fmt.Printf("EventNodeDeleted: %s\n", event.Path)
				}
			case zk.EventNodeChildrenChanged:
				{
					initMap(event.Path)
					fmt.Printf("EventNodeChildrenChanged: %s\n", event.Path)
				}
			case zk.EventNodeCreated:
				{
					fmt.Printf("EventNodeCreated: %s\n", event.Path)
				}
			}
		})
		setAll()
	})

	return manager
}

func (m *Manager) Get(configPath string) string {
	if strings.HasPrefix(configPath, "/connection") || strings.HasPrefix(configPath, "/config") {
		v := m.configMap[configPath]
		if v == "" {
			v, _, err := zoo.GetW(configPath)

			if err == nil {
				m.configMap[configPath] = v
				return v
			}
			return ""
		}
	}
	return m.configMap[configPath]
}

func initMap(path string) {
	if path == "" {
		return
	}
	if strings.EqualFold(path, "/") && !(strings.Index(Paths, path) >= 0) {
		return
	}
	var zk = zoo
	children, _, err := zk.GetChildrenW(path)

	if err != nil {
		panic(err)
	}
	if children == nil || len(children) <= 0 {
		return
	}
	for _, sPath := range children {
		var nextPath = path

		if strings.EqualFold(path, "/") && !(strings.Index(Paths, sPath) >= 0) {
			continue
		}
		if !strings.HasSuffix(nextPath, "/") {
			nextPath += "/"
		}
		nextPath += sPath
		v, _, _ := zk.GetW(nextPath)
		manager.configMap[nextPath] = v
	}
}

func (m *Manager) Create(name string, v interface{}) error {
	var zk = zoo
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = zk.Create(name, data, -1)
	return err
}

func (m *Manager) Set(path string, v interface{}) error {
	var zk = zoo
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = zk.Set(path, data, -1)
	return err
}

func (m *Manager) update(path string) {
	if path == "" {
		return
	}
	if strings.EqualFold(path, "/") && !(strings.Index(Paths, path) >= 0) {
		return
	}
	var zk = zoo
	v, _, err := zk.GetW(path)
	if err != nil {
		panic(err)
	}
	manager.configMap[path] = v
}

func (m *Manager) delete(path string) {
	if path == "" {
		return
	}
	delete(m.configMap, path)
}

func setAll() {
	paths := strings.Split(Paths, ",")
	for _, v := range paths {
		initMap("/" + v)
	}
}

func (m *Manager) Dispose() {
	zoo.conn.Close()
}
