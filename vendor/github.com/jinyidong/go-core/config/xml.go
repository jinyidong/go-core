package config

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type XmlConfig struct {
	Zookeeper string `xml:"Zookeeper"`
}

var BasicConfig = &XmlConfig{}

func init() {
	configPath := "/config/zk.config" //注意在windows下，磁盘目录为GOPATH所对应的目录
	file, err := os.Open(configPath)
	if err != nil {
		log.Printf("error: %v", err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("error: %v", err)
	}
	err = xml.Unmarshal(data, BasicConfig)
	if err != nil {
		log.Printf("error: %v", err)
	}
}

func (c *XmlConfig) Servers() []string {
	ss := strings.Split(BasicConfig.Zookeeper, ",")
	return ss
}
