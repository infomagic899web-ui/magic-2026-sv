package utils

import (
	"log"
)

func GetVariables() (string, string) {
	onlineURI := GetEnv("MONGO_ONLINE_URI")
	offlineURI := GetEnv("MONGO_OFFLINE_URI")

	if onlineURI == "" || offlineURI == "" {
		log.Fatal("MONGO_ONLINE_URI and MONGO_OFFLINE_URI environment variables must be set")
	}

	if onlineURI != "" {
		log.Print("Online URI set ")
	}

	if offlineURI != "" {
		log.Print("Offline URI set ")
	}

	return onlineURI, offlineURI
}
