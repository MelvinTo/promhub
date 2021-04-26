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

	log.Info(config)

	for _, pc := range config.ProberConfigs {
		t := pc["type"].(string)
		switch t {
		case "bandwagon":
			bh := prober.BandwagonHost{}
			bh.LoadFromConfig(pc)
			go bh.ScheduleRun()
		}
	}
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
