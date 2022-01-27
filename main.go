package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/luxcgo/go-gallery/controllers"
	"github.com/luxcgo/go-gallery/models"
	"github.com/robfig/cron/v3"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "root"
	password = "secret"
	dbname   = "luxcgo_gallery"
)

var (
	confDir string
)

func init() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	flag.StringVar(&confDir, "conf_dir", path, "conf文件存放路径")
	flag.Parse()
}

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
	galleriesC.Run()
	c := cron.New()
	c.AddJob("0 22 * * *", galleriesC)
	c.AddJob("0 12 * * *", galleriesC)
	c.Start()

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.FAQ).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(notFound)

	r.HandleFunc("/galleries", galleriesC.Create).Methods("POST")
	r.HandleFunc("/galleries", galleriesC.Index).Methods("GET").
		Name(controllers.IndexGalleries)

	http.ListenAndServe(":"+fmt.Sprint(cfg.Port), r)
}
