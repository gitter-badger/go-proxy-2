package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var config Config

type Config struct {
	ListenAddr string `json:"listen"`
	RemoteAddr string `json:"remote"`
	Proxy	   string `json:"proxy"`
	Password   string `json:"password"`
}


// 保存配置到配置文件
func (config *Config) SaveConfig(configPath string) {
	configJson, _ := json.MarshalIndent(config, "", "	")
	err := ioutil.WriteFile(configPath, configJson, 0644)
	if err != nil {
		fmt.Errorf("保存配置到文件 %s 出错: %s", configPath, err)
	}
	log.Printf("保存配置到文件 %s 成功\n", configPath)
}

func (config *Config) ReadConfig(configPath string) {
	// 如果配置文件存在，就读取配置文件中的配置 assign 到 config
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		log.Printf("从文件 %s 中读取配置\n", configPath)
		file, err := os.Open(configPath)
		if err != nil {
			log.Fatalf("打开配置文件 %s 出错:%s", configPath, err)
		}
		defer file.Close()

		err = json.NewDecoder(file).Decode(config)
		if err != nil {
			log.Fatalf("格式不合法的 JSON 配置文件:\n%s", file)
		}
	}
}
