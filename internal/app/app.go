package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/stefanop21/goProject/internal/api"
	"github.com/stefanop21/goProject/internal/store"
	"github.com/stefanop21/goProject/migrations"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// our stores will go here
	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPosgresTokenStore(pgDB)

	// our handlers will go here
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHandler,
		TokenHandler:   tokenHandler,
		DB:             pgDB,
	}

	// nil for valid error types
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
