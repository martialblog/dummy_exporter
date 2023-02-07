package main

import (
	"flag"
	"log"
	"net/http"
)

// nolint: gochecknoglobals
var (
	// These get filled at build time with the proper vaules
	version = "development"
	commit  = "HEAD"
)

var (
	listenAddress = flag.String("web.listen-address", ":9123", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")
	configPath    = flag.String("config.file", "dummy.json", "Dummy exporter configuration file.")
)

func main() {
	flag.Parse()

	var config Config

	if err := config.LoadConfig(*configPath); err != nil {
		log.Fatal("Error loading config", "err", err)
	}

	dummy := DummyMetricHandler{
		Metrics: config.Metrics,
	}

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

	log.Printf("Version: %s Commit: %s", version, commit)
	log.Printf("Listening on address: %s", *listenAddress)

	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
