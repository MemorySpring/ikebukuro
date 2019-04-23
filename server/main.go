package main

import (
	"Ikebukuro/server/connection"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"os"
)

type Config struct {
	Mysql string `json:"mysql"`
}

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("配置加载失败", err)
		return
	}
	connection.DataBaseInit(config.Mysql)
	ormer := orm.NewOrm()
	err = ormer.Using("default")
	if err != nil {
		fmt.Println(err)
		return
	}
	diskServer := &connection.DiskServer{
		Ormer: ormer,
		Port:  ":1318",
	}
	diskServer.Start()
}

func loadConfig() (Config, error) {
	file, err := os.Open("./server/conf/config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := Config{}
	err = decoder.Decode(&conf)
	return conf, err
}
