package util

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

func GetConfig() {

	_, err := toml.DecodeFile("config.toml", &Config)
	if err != nil {
		fmt.Println("ERROR: Couldnt load config.toml: ", err)
		os.Exit(1)
	}
	fmt.Println("Parsed toml")
	fmt.Println("Bot prefix is: ", Config.Prefix)

	files, err := os.ReadDir(Config.ImageDir)
	if err != nil {
		fmt.Println("ERROR: Couldnt read contents of ", Config.ImageDir, ": ", err)
		os.Exit(1)
	}
	for _, file := range files {
		if !file.IsDir() {
			Images = append(Config.Images, file.Name())
		}
	}
}
