package settings

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type ConfigInfo struct {
	Driver							string	`json:"driver"`
	DB 								string	`json:"db"`
	PortRun							string	`json:"portRun"`
	TokenAuthPostman 				string	`json:"tokenAuthPostman"`
	TokenAuthBetweenServices 		string	`json:"tokenAuthBetweenServices"`
}


var Config ConfigInfo 


func ConfigSetup(FileName string) {
	// Чтение с файла Settings JSON
	byteValue, err := ioutil.ReadFile(FileName)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	err = json.Unmarshal(byteValue, &Config)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
}
