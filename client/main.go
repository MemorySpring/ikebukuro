package main

import (
	"Ikebukuro/client/connnection"
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Server string
}

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("配置文件加载失败")
		return
	}
	client := connnection.DiskClient{
		Addr: config.Server,
	}
	err = client.SignUp("baojizhong", "113118")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("注册完毕")
	err = client.Login("baojizhong", "113118")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("登录完毕")
}

func loadConfig() (Config, error) {
	file, _ := os.Open("./client/conf/config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := Config{}
	err := decoder.Decode(&conf)
	return conf, err
}
