package main

import (
	"flag"
	"log"
	"net/http"
)

var build = "development"

var (
	listenAddress = flag.String("web.listen-address", ":9123", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")
)

func main() {

	// Handle metrics endpoint
	var dummy DummyMetricHandler
	http.Handle(*metricsPath, dummy)

	// Let the User know where the metrics endpoint is
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
			<head><title>Dummy Exporter</title></head>
			<body>
			<h1>Dummy Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Printf("Version: %s", build)
	log.Printf("Listening on address: %s", *listenAddress)

	log.Fatal(http.ListenAndServe(*listenAddress, nil))

}
