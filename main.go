package main

import (
	"crypto/tls"
	"flag"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/crerwin/dcosauth/pkg/dcosauth"
	"github.com/matt-deboer/go-marathon"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	listenAddress = flag.String(
		"web.listen-address", ":9088",
		"Address to listen on for web interface and telemetry.")

	metricsPath = flag.String(
		"web.telemetry-path", "/metrics",
		"Path under which to expose metrics.")

	marathonUri = flag.String(
		"marathon.uri", "http://marathon.mesos:8080",
		"URI of Marathon")

	dcosToken = flag.String(
		"dcos.token", "none",
		"Bearer token for DC/OS authentication")

	dcosServiceAccoundUID = flag.String(
		"dcos.serviceaccountuid", "none",
		"UID for DC/OS Service Account")

	dcosServiceAccountPrivateKey = flag.String(
		"dcos.serviceaccountprivatekey", "none",
		"Private Key for DC/OS Service Account")
)

func dcosSetup(uri *url.URL) marathon.Config {
	config := marathon.NewDefaultConfig()
	// add /marathon to end of admin router url
	config.URL = "https://" + uri.String()
	if *dcosToken != "none" {
		config.DCOSToken = *dcosToken
	}
	authURL := strings.TrimSuffix(uri.String(), "/marathon")

	if *dcosServiceAccoundUID != "none" && *dcosServiceAccountPrivateKey != "none" {
		dcosauther := dcosauth.New(authURL, *dcosServiceAccoundUID, *dcosServiceAccountPrivateKey)
		token, _ := dcosauther.Token()
		config.DCOSToken = token
	}
	return config
}

func marathonConnect(uri *url.URL) error {
	config := marathon.NewDefaultConfig()
	if *dcosToken != "none" || *dcosServiceAccoundUID != "non" || *dcosServiceAccountPrivateKey != "none" {
		config = dcosSetup(uri)
	} else {
		config.URL = uri.String()
	}

	if uri.User != nil {
		if passwd, ok := uri.User.Password(); ok {
			config.HTTPBasicPassword = passwd
			config.HTTPBasicAuthUser = uri.User.Username()
		}
	}
	config.HTTPClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).Dial,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	log.Debugln("Connecting to Marathon")
	client, err := marathon.NewClient(config)
	if err != nil {
		return err
	}

	info, err := client.Info()
	if err != nil {
		return err
	}

	log.Debugf("Connected to Marathon! Name=%s, Version=%s\n", info.Name, info.Version)
	return nil
}

func main() {
	flag.Parse()
	uri, err := url.Parse(*marathonUri)
	if err != nil {
		log.Fatal(err)
	}

	retryTimeout := time.Duration(10 * time.Second)
	for {
		err := marathonConnect(uri)
		if err == nil {
			break
		}
		if err.Error() == "all the Marathon hosts are presently down" && (*dcosToken != "none" || *dcosServiceAccoundUID != "none") {
			// if we're targeting the DC/OS admin router the marathon client will complain
			// about marathon hosts being down, but that's okay
			break
		}

		log.Infof("Problem connecting to Marathon: %v", err)
		log.Infof("Couldn't connect to Marathon! at %v Trying again in %v", uri, retryTimeout)
		time.Sleep(retryTimeout)
	}

	exporter := NewExporter(&scraper{uri}, defaultNamespace)
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Marathon Exporter</title></head>
             <body>
             <h1>Marathon Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	log.Info("Starting Server: ", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
