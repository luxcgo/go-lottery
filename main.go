package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/luxcgo/go-gallery/controllers"
	"github.com/luxcgo/go-gallery/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "root"
	password = "secret"
	dbname   = "luxcgo_gallery"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "肥肠抱歉,你要找的页面不见了")
}

// A helper function that panics on any error
func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	cfg := LoadConfig()
	dbCfg := cfg.Database
	services, err := models.NewServices(
		models.WithGorm(dbCfg.ConnectionInfo()),
		// Only log when not in prod
		models.WithNewLogger(!cfg.IsProd()),
		models.WithGallery(),
	)
	if err != nil {
		panic(err)
	}

	services.AutoMigrate()
	// services.DestructiveReset()

	r := mux.NewRouter()
	staticC := controllers.NewStatic()
	galleriesC := controllers.NewGalleries(services.Gallery, r)
	// galleriesC.Crawl()

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.FAQ).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(notFound)

	r.HandleFunc("/galleries", galleriesC.Create).Methods("POST")
	r.HandleFunc("/galleries", galleriesC.Index).Methods("GET").
		Name(controllers.IndexGalleries)

	http.ListenAndServe(":"+fmt.Sprint(cfg.Port), r)
}
