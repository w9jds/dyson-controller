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

	email    string
	password string
	country  string
)

func setup() {
	email = os.Getenv("DYSON_EMAIL")
	password = os.Getenv("DYSON_PASS")
	country = os.Getenv("COUNTRY")

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
	_, err := dysonClient.GetDevices()
	if err != nil {
		log.Panic("Error occurred while pulling account Dyson Devices")
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
