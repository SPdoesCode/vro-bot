package util

import (
	"encoding/json"
	"fmt"
	"os"
)

func writeData() {
	data, err := json.MarshalIndent(userdata, "", " ")
	if err != nil {
		fmt.Println("ERROR: Couldnt update save.json: ", err)
		os.Exit(0)
	}
	err = os.WriteFile("save.json", data, 0644)
	if err != nil {
		fmt.Println("ERROR: Couldnt save to data.json: ", err)
	}
}

func GetData() {

	data, err := os.ReadFile("save.json")
	if err != nil {
		fmt.Println("ERROR: Couldnt read the save.json: ", err)
		os.Exit(1)
	}
	err = json.Unmarshal(data, &userdata)
	if err != nil {
		fmt.Println("ERROR: Couldnt get the data from the save.json: ", err)
		os.Exit(1)
	}

}
