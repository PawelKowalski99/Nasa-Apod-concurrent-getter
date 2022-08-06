package main

import (
	"github.com/gogoapps/manager"
	"github.com/gogoapps/routes"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"
)



func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("could not load .env file")
	}

	m, err  := manager.NewManager()
	if err != nil {
		logrus.Fatalf("could not create manager: %v", err)
	}

	err = routes.InitRoutes(m)

	http.ListenAndServe(manager.PORT, m.R)
}

