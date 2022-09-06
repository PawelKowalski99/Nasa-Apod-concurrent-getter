package main

import (
	config "github.com/PawelKowalski99/gogapps/config"
	"log"
	"net/http"

	"github.com/PawelKowalski99/gogapps/server"
	"github.com/sirupsen/logrus"
)

func main() {
	c, err := config.Init()
	if err != nil {
		log.Fatal(err.Error())
	}

	s, err := server.Init(c)
	if err != nil {
		logrus.Fatalf("could not create server: %v", err)
	}

	logrus.Infof("Listening on port: %s", s.Config.GetPort())
	err = http.ListenAndServe(":"+s.Config.GetPort(), s.R)
	if err != nil {
		log.Fatal(err.Error())
	}
}
