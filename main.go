package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
)

var configFile = flag.String("c", "etc/sinai-api.yaml", "the config file")
var config = &Config{}

type Config struct {
	JmsToken    string `yaml:"jms_token"`
	JmsAddr     string `yaml:"jms_addr"`
	Interval    int    `yaml:"interval"`
	DialTimeout int    `yaml:"dial_timeout"`
	HttpPort    int    `yaml:"http_port"`
}

func main() {
	flag.Parse()

	f, err := os.ReadFile(*configFile)
	if err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(f, config); err != nil {
		panic(err)
	}

	jmsClient := NewJmsClient(config.JmsAddr, config.JmsToken)

	connections, err = jmsClient.GatewayList()
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(time.Duration(config.Interval) * time.Second)
	defer ticker.Stop()

	// start one go routine to monitor gateway can reachable
	go func() {
		for {
			select {
			case <-ticker.C:
				for i := range connections {
					status := checkConnection(connections[i])
					updateConnectionStatus(&status)
				}
			}
		}
	}()

	prometheus.MustRegister(NewConnectionCollector())

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/gateways", getGatewayList)

	fmt.Printf("jms_domain_exporter will start at :%d", config.HttpPort)
	if err = http.ListenAndServe(fmt.Sprintf(":%d", config.HttpPort), nil); err != nil {
		fmt.Println("start http server fail:", err)
	}
}
func getGatewayList(w http.ResponseWriter, r *http.Request) {
	if len(connections) > 0 {
		ret, _ := json.Marshal(connections)
		_, _ = w.Write(ret)
	}
}
