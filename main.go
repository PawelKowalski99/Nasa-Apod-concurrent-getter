package main

import (
	"log"
	"net/http"

	"github.com/PawelKowalski99/gogapps/manager"
	"github.com/PawelKowalski99/gogapps/routes"

	"github.com/sirupsen/logrus"
)

func main() {
	m, err := manager.NewManager()
	if err != nil {
		logrus.Fatalf("could not create manager: %v", err)
	}

	err = routes.InitRoutes(m)
	if err != nil {
		log.Fatal(err.Error())
	}

	m.L.Printf("Running on Port %s", m.Port)
	err = http.ListenAndServe(m.Port, m.R)
	if err != nil {
		log.Fatal(err.Error())
	}
}
