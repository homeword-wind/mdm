package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

const Ext = "toml"

func init() {
	// Set Profile Path
	work, _ := os.Getwd()
	viper.AddConfigPath(work + "/conf")

	viper.SetConfigType(Ext)
}

func main() {
	fmt.Println("Thanks for using the mdm")

	/* ---------- Main ---------- */

	// Set File Name
	viper.SetConfigName("main")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("viper.ReadInConfig(): %v", err)
	}
	// 获取全部文件内容
	fmt.Println("all settings: ", viper.AllSettings())
	fmt.Println("--------------")
	// 根据内容类型，解析出不同类型
	fmt.Println(viper.GetString("database.server"))
	fmt.Println(viper.GetIntSlice("database.ports"))
	fmt.Println("--------------")
	fmt.Println(viper.GetString("servers.alpha.ip"))
}
