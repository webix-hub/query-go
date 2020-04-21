package main

import (
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/unrolled/render"

	"log"
)

var format = render.New()

func main() {
	// config
	Config.LoadFromFile("./config.yml")

	// db connection
	db, err := sqlx.Connect("mysql", Config.DataSourceName())
	if err != nil {
		log.Fatal(err)
	}

	// http router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)

	r.Post("/api/data/{table}", func(w http.ResponseWriter, r *http.Request) {
		table := chi.URLParam(r, "table")

		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			format.Text(w, 500, err.Error())
			return
		}

		err = getDataFromDB(w, db, table, body)
		if err != nil {
			format.Text(w, 500, err.Error())
			return
		}
	})

	r.Get("/api/data/{table}/{field}/suggest", func(w http.ResponseWriter, r *http.Request) {
		table := chi.URLParam(r, "table")
		field := chi.URLParam(r, "field")

		err := getSuggestFromDB(w, db, table, field)
		if err != nil {
			format.Text(w, 500, err.Error())
			return
		}
	})

	log.Printf("Start server at %s", Config.Server.Port)
	http.ListenAndServe(Config.Server.Port, r)
}
