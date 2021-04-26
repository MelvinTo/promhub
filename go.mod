module github.com/melvinto/promhub

go 1.15

replace github.com/melvinto/promhub/prober => /prober

require (
	github.com/go-resty/resty/v2 v2.6.0
	github.com/prometheus/client_golang v1.10.0
	github.com/sirupsen/logrus v1.6.0
	gopkg.in/yaml.v2 v2.3.0
)
