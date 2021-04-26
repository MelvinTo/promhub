package main

import (
	"net/http"

	"github.com/melvinto/promhub/config"
	"github.com/melvinto/promhub/prober"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/sirupsen/logrus"
)

func main() {
	config := config.GlobalConfig()

	probers := []prober.Prober{}

	for _, pc := range config.ProberConfigs {
		t := pc["type"].(string)
		name := pc["name"].(string)
		log.Infof("initializing %v prober %v", t, name)
		switch t {
		case "bandwagon":
			bh := prober.BandwagonHost{}
			bh.LoadFromConfig(pc)
			probers = append(probers, bh)
		}
	}

	h := promhttp.Handler()
	
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		for _, p := range probers {
			p.Run()
		}
		h.ServeHTTP(w, r)
	})


	http.ListenAndServe(":2112", nil)
}
