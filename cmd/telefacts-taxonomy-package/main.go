package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"ecksbee.com/telefacts-taxonomy-package/internal/cache"
	"ecksbee.com/telefacts-taxonomy-package/internal/web"
	"ecksbee.com/telefacts-taxonomy-package/pkg/install"
	"ecksbee.com/telefacts-taxonomy-package/pkg/throttle"
	"ecksbee.com/telefacts/pkg/hydratables"
	"ecksbee.com/telefacts/pkg/serializables"
)

var (
	zipVar    string
	volumeVar string
)

func main() {
	flag.StringVar(&zipVar, "zip", "", "taxonomy package zip file")
	flag.StringVar(&volumeVar, "volume", "", "taxonomy package zip file")
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	fmt.Printf("%s zip\n", zipVar)
	if zipVar == "" && volumeVar != "" {
		fmt.Println("-zip is empty")
		return
	}
	if volumeVar == "" && zipVar != "" {
		fmt.Println("-volume is empty")
		return
	}
	if zipVar != "" && volumeVar != "" {
		throttle.StartSECThrottle()
		_, err := install.Run(zipVar, volumeVar, throttle.Throttle)
		if err != nil {
			panic(err)
		}
		return
	}
	var ctx = context.Background()
	srv := setupServer()
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	listenForShutdown(ctx, wait, srv)
}

func setupServer() *http.Server {
	appCache := cache.NewCache()
	dir, err := os.Getwd()
	if err != nil {
		dir = path.Join(".")
	}
	wd := os.Getenv("WD")
	if wd == "" {
		wd = dir
	}
	serializables.WorkingDirectoryPath = path.Join(wd, "wd")
	gts := os.Getenv("GTS")
	if gts == "" {
		gts = dir
	}
	serializables.GlobalTaxonomySetPath = path.Join(gts, "gts")
	hydratables.InjectCache(appCache)

	r := web.NewRouter()

	fmt.Println("telefacts-taxonomy-package-manager<-0.0.0.0:8080")
	return &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
}

func listenForShutdown(ctx context.Context, grace time.Duration, srv *http.Server) {
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Shutting down")
	ctx, cancel := context.WithTimeout(ctx, grace)
	defer cancel()
	srv.Shutdown(ctx)
	os.Exit(0)
}
