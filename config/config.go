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
	node *Node
}

var manager *Manager
var once sync.Once
var zoo *ZK

func NewManager() *Manager {

	once.Do(func() {
		manager = &Manager{}

		zoo = newZK(strings.Split(BasicConfig.Zookeeper, ","), func(event zk.Event) {
			switch event.State {
			case zk.StateExpired:
				{
					fmt.Println("StateExpired")
					manager.setNodes("/")
					fmt.Println("setNodes /")
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
					newPath := event.Path
					manager.setNodes(newPath)
					fmt.Printf("EventNodeChildrenChanged: %s\n", newPath)
				}
			case zk.EventNodeCreated:
				{
					fmt.Printf("EventNodeCreated: %s\n", event.Path)
				}
			}
		})

		manager.node = initNodes()
	})

	return manager
}

func initNodes() *Node {

	var node = newNode()

	nextNode(node, "/")

	if node.Children == nil {
		panic("no child")
	}

	return node
}

func nextNode(node *Node, path string) {
	if node == nil {
		return
	}
	var zookeeper = zoo

	children, _, err := zookeeper.GetChildrenW(path)

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

		v, _, _ := zookeeper.GetW(nextPath)

		var child = &Node{
			Path:  sPath,
			Value: v,
		}

		node.appendChild(child)

		nextNode(child, nextPath)
	}
}

func (m *Manager) GetChildren(path string) []*Node {

	node, _ := getNode(m.node, path)

	if node == nil {
		return nil
	}

	return node.Children
}

func (m *Manager) Get(path string) string {

	node, _ := getNode(m.node, path)

	if node == nil {
		return ""
	}
	return node.Value
}

func (m *Manager) Create(name string, v interface{}) error {
	var zookeeper = zoo

	data, err := json.Marshal(v)

	if err != nil {
		return err
	}

	err = zookeeper.Create(name, data, -1)

	return err
}

func (m *Manager) Set(path string, v interface{}) error {
	var zookeeper = zoo

	data, err := json.Marshal(v)

	if err != nil {
		return err
	}

	err = zookeeper.Set(path, data, -1)

	return err
}

func (m *Manager) add(path string, v string) {

	node, paths := getNode(m.node, path)

	newNode := &Node{Path: paths[len(paths)-1], Value: v}

	node.appendChild(newNode)
}

func (m *Manager) update(path string) {
	node, _ := getNode(m.node, path)

	v, _, _ := zoo.GetW(path)

	node.Value = v
}

func (m *Manager) delete(path string) {

	node, _ := getNode(m.node, path)

	if node == nil {
		return
	}

	paths := strings.Split(path, "/")

	newPaths := paths[1 : len(paths)-1]

	if newPaths == nil {
		return
	}

	newPath := "/" + strings.Join(newPaths, "/")

	pNode, _ := getNode(m.node, newPath)

	if pNode == nil {
		return
	}

	removeIndex := 0

	for i := 0; i < len(pNode.Children); i++ {
		if node == pNode.Children[i] {
			removeIndex = i
			break
		}
	}

	pNode.Children = append(pNode.Children[:removeIndex], pNode.Children[removeIndex+1:]...)
}

func (m *Manager) setNodes(path string) {
	node, _ := getNode(m.node, path)

	nextNode(node, path)
}

func getNode(node *Node, path string) (*Node, []string) {
	if node == nil {
		panic("node is nil")
	}

	if path == "" {
		panic("path is empty")
	}

	paths := strings.Split(path, "/")

	if len(paths) <= 1 {
		panic("path len <=1")
	}

	var returnNode = node

	for i := 1; i < len(paths); i++ {
		returnNode = childSelector(returnNode, paths[i])

		if returnNode == nil {
			break
		}
	}

	return returnNode, paths
}

func childSelector(n *Node, p string) *Node {

	if p == "" {
		return nil
	}

	var childNode *Node
	for i, child := range n.Children {

		if child == nil {
			continue
		}

		if child.Path == "" {
			continue
		}

		if p == child.Path {
			childNode = n.Children[i]
			break
		}
	}

	return childNode
}

func (m *Manager) Dispose() {
	zoo.conn.Close()
}
