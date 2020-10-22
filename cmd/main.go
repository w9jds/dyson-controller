package main

import (
	"crypto/tls"
	"dyson-controller/pkg/dyson"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	dysonClient *dyson.Client

	email     string
	password  string
	country   string
	ipAddress string
)

func setup() {
	email = os.Getenv("DYSON_EMAIL")
	password = os.Getenv("DYSON_PASS")
	country = os.Getenv("COUNTRY")
	ipAddress = os.Getenv("IP_ADDRESS")

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	dysonClient = &dyson.Client{
		Client: httpClient,
	}
}

func main() {
	setup()

	dysonClient.Login(email, password, country)
	devices, err := dysonClient.GetDevices()
	if err != nil {
		log.Panic("Error occurred while pulling account Dyson Devices")
	}

	devices[0].Connect(ipAddress)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
