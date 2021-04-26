package prober

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

var (
	bandwagonAPIURL          = "https://api.64clouds.com/v1/getServiceInfo"
	bandwagonUsedDataMetrics = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ph_bandwagon_used_data_bytes",
		Help: "bandwidth used in bytes",
	}, []string{"hostname", "veid"})
	bandwagonTimeBeforeReset = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ph_bandwagon_time_before_reset_seconds",
		Help: "time before reset in seconds",
	}, []string{"hostname", "veid"})
	bandwagonPlanMonthlyData = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ph_bandwagon_plan_monthly_data_bytes",
		Help: "planned monthly data in bytes",
	}, []string{"hostname", "veid"})
)

type BandwagonHost struct {
	VEID        string `yaml:"veid"`
	APIKey      string `yaml:"api_key"`
	RunInterval int    `yaml:"interval"`
}

type BandwagonHostInfo struct {
	Hostname        string `json:"hostname"`
	DataCounter     int64  `json:"data_counter"`
	DataNextReset   int64  `json:"data_next_reset"`
	PlanMonthlyData int64  `json:"plan_monthly_data"`
}

func (b BandwagonHostInfo) String() string {
	gb := int(b.DataCounter / 1024 / 1024 / 1024)
	perc := b.DataCounter * 100 / b.PlanMonthlyData
	now := time.Now().Unix()
	diff := b.DataNextReset - now

	return fmt.Sprintf("%v: %vGB (%v%%) used, reset in %v days", b.Hostname, gb, perc, int(diff/3600/24))
}

func (b *BandwagonHost) Name() string {
	return "bandwagon"
}

func (b *BandwagonHost) LoadFromConfig(c map[string]interface{}) {
	bytes, err := yaml.Marshal(&c)
	if err != nil {
		log.Error("Failed to marshal config, err:", err)
		return
	}

	err = yaml.Unmarshal(bytes, b)
	if err != nil {
		log.Error("Failed to load config, err:", err)
		return
	}
}

func (b *BandwagonHost) ScheduleRun() {
	for {
		b.Run()
		time.Sleep(time.Duration(b.RunInterval) * time.Second)
	}
}

func (b BandwagonHost) Run() {

	// Create a Resty Client
	client := resty.New()

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"veid":    b.VEID,
			"api_key": b.APIKey,
		}).
		SetHeader("Accept", "application/json").
		SetResult(&BandwagonHostInfo{}).
		Get(bandwagonAPIURL)

	if err != nil {
		log.Error("Failed to call api, err:", err)
		return
	}

	result := resp.String()

	bhi := BandwagonHostInfo{}

	err = json.Unmarshal([]byte(result), &bhi)
	if err != nil {
		log.Error("Failed to convert json to BandwagonHostInfo, err", err)
		return
	}

	log.Info(bhi)

	// update prometheus metrics
	bandwagonUsedDataMetrics.WithLabelValues(bhi.Hostname, b.VEID).Set(float64(bhi.DataCounter))
	bandwagonPlanMonthlyData.WithLabelValues(bhi.Hostname, b.VEID).Set(float64(bhi.PlanMonthlyData))

	now := time.Now().Unix()
	diff := bhi.DataNextReset - now

	bandwagonTimeBeforeReset.WithLabelValues(bhi.Hostname, b.VEID).Set(float64(diff))
}
