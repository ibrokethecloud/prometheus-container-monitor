package main

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ibrokethecloud/rancher-events/events"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rancher/go-rancher-metadata/metadata"

	log "github.com/Sirupsen/logrus"
)

const (
	metadataUrl = "http://rancher-metadata/latest"
)

var currentState = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "container_state",
	Help: "Current state of container 1=reachable, 0=unreachable",
}, []string{"container_name"})

func init() {
	// Metrics need to be registered //
	prometheus.MustRegister(currentState)
	log.SetFormatter(&log.JSONFormatter{})
}

// Setting up some global variables using our environment definitions //

var MODE = os.Getenv("MODE")
var ENV = os.Getenv("ENV")
var PORT = os.Getenv("PORT")

func main() {

	// Starting the handler for metrics //
	http.Handle("/metrics", promhttp.Handler())

	var wg sync.WaitGroup

	wg.Add(3)

	http_listen := ":" + PORT

	go func() {
		defer wg.Done()
		log.Info("About to start http server")
		log.Error(http.ListenAndServe(http_listen, nil))
	}()

	go func() {
		defer wg.Done()
		for { // Run an infinite loop for metric collections //
			m := metadata.NewClient(metadataUrl)
			retContainers, err := m.GetContainers()
			checkError(err)
			pingContainers(retContainers)
			time.Sleep(60 * time.Second)
		}
	}()

	//Stream Docker Events from Rancher API
	go func() {
		defer wg.Done()
		index := getServiceIndex()
		log.Infof("Service index returned is %d",index)
		for index == 1 {
			events.GetContainerEvents()
			time.Sleep(10 * time.Minute)
		}
	}()

	wg.Wait()

}
