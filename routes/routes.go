package routes

import (
	"context"
	"fmt"

	"github.com/PawelKowalski99/gogapps/manager"
)

func InitRoutes(m *manager.Manager) error {
	for _, f := range []func(m *manager.Manager)error{
		pictures,
	} {
		err := 	f(m)
		if err != nil {
			m.L.Errorf("could not init routes: %v", err)
			return fmt.Errorf("could not init routes: %v", err)
		}
	}

	return nil
}

func pictures(m *manager.Manager) error {

	m.R.Get("/pictures", m.Provider.GetPictures(context.Background()))

	return nil
}
