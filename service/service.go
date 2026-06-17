package service

import (
	"errors"
	"os"
	"strings"

	usecaseimpl "github.com/YagoSchramm/GoNotify/domain/usecase/impl"
	repoimpl "github.com/YagoSchramm/GoNotify/infrastructure/datastore/repository/impl"
	"github.com/YagoSchramm/GoNotify/infrastructure/foundation/db"
	approuter "github.com/YagoSchramm/GoNotify/infrastructure/router"
	modules "github.com/YagoSchramm/GoNotify/infrastructure/router/module"
	"github.com/gorilla/mux"
)

func Build() (*mux.Router, func(), error) {
	content, err := os.ReadFile(".env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, func() {}, err
	}

	if err == nil {
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			key, value, found := strings.Cut(line, "=")
			if !found {
				continue
			}

			key = strings.TrimSpace(key)
			value = strings.TrimSpace(value)
			value = strings.Trim(value, `"'`)

			if key == "" {
				continue
			}

			if os.Getenv(key) == "" {
				_ = os.Setenv(key, value)
			}
		}
	}

	dsn := os.Getenv("DATABASE_URL")
	secret := os.Getenv("JWT_SECRET")

	if dsn == "" {
		return nil, func() {}, errors.New("DATABASE_URL is not set")
	}
	if secret == "" {
		return nil, func() {}, errors.New("JWT_SECRET is not set")
	}

	dbConn, err := db.NewPostgresConnection(dsn)
	if err != nil {
		return nil, func() {}, err
	}

	authRepository := repoimpl.NewAuthRepository(dbConn)
	triggerRepository := repoimpl.NewTriggerRepository(dbConn)
	notificationRepository := repoimpl.NewNotificationRepository(dbConn)
	cleanup := func() {
		_ = dbConn.Close()
	}

	authUseCase := usecaseimpl.NewAuthRepository(authRepository, secret)
	notificationUseCase := usecaseimpl.NewNotificationUseCase(notificationRepository, triggerRepository)
	triggerUseCase := usecaseimpl.NewTriggerUseCase(triggerRepository)

	authModule := modules.NewAuthModule(authUseCase, secret)
	notificationModule := modules.NewNotificationModule(notificationUseCase)
	triggerModule := modules.NewTriggerModule(triggerUseCase)

	router := mux.NewRouter()
	approuter.Mount(
		router,
		authModule.Middlewares(),
		authModule,
		notificationModule,
		triggerModule,
	)

	return router, cleanup, nil
}
