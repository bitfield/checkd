package main

import (
	"log"
	"net/http"

	"github.com/bitfield/pkg/checkd"
	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

func main() {
	log.Println("startup")
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to read config: %s\n", err)
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Println("config change detected")
		initChecks(v)
	})
	initChecks(v)
	checkd.Run()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initChecks(v *viper.Viper) {
	if err := checkd.Init(v); err != nil {
		log.Fatal(err)
	}

}
