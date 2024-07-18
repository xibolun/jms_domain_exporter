package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type ConnectionStatus struct {
	IP       string
	Port     int
	Name     string
	IsUp     bool
	LastSeen time.Time
}

var (
	connectionStatuses = make(map[string]*ConnectionStatus)
	connections        = make([]*ConnectionStatus, 0)
	mutex              sync.RWMutex
)

type ConnectionCollector struct {
	upDesc *prometheus.Desc
}

func NewConnectionCollector() *ConnectionCollector {
	return &ConnectionCollector{
		upDesc: prometheus.NewDesc(
			"connection_up",
			"Indicates if the connection is up (1) or down (0)",
			[]string{"ip", "port", "name"},
			nil,
		),
	}
}

func (c *ConnectionCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.upDesc
}

func (c *ConnectionCollector) Collect(ch chan<- prometheus.Metric) {
	mutex.RLock()
	defer mutex.RUnlock()

	for _, status := range connectionStatuses {
		isUp := 0.0
		if status.IsUp {
			isUp = 1.0
		}

		ip, port, name := status.IP, status.Port, status.Name

		ch <- prometheus.MustNewConstMetric(
			c.upDesc,
			prometheus.GaugeValue,
			isUp,
			ip,
			strconv.Itoa(port),
			name,
		)
	}
}

func checkConnection(status *ConnectionStatus) ConnectionStatus {
	// 连接到指定IP和端口
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", status.IP, status.Port), time.Duration(config.DialTimeout)*time.Second)
	if err != nil {
		status.IsUp = false
	} else {
		status.IsUp = true
		conn.Close()
	}

	// 更新最后一次检查的时间
	status.LastSeen = time.Now()

	return *status
}

func updateConnectionStatus(status *ConnectionStatus) {
	mutex.Lock()
	defer mutex.Unlock()

	key := getConnectionKey(status.IP, status.Port, status.Name)
	connectionStatuses[key] = status
}

func getConnectionKey(ip string, port int, name string) string {
	return fmt.Sprintf("%s:%d:%s", ip, port, name)
}
