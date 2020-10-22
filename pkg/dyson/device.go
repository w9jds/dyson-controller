package dyson

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Device is the spec of a registered Dyson Device on your account
type Device struct {
	client mqtt.Client

	Serial              string
	Name                string
	Version             string
	LocalCredentials    string
	ProductType         string
	ConnectionType      string
	AutoUpdate          bool
	NewVersionAvailable bool
}

func (dyson *Device) decipherCredentials() string {
	iv := []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	key := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20}

	cipherText, err := b64.StdEncoding.DecodeString(dyson.LocalCredentials)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		panic(err)
	}

	stream := cipher.NewCBCDecrypter(block, iv)
	stream.CryptBlocks(cipherText, cipherText)

	var credentials map[string]string
	err = json.Unmarshal(bytes.Trim(cipherText, "\x08"), &credentials)
	if err != nil {
		panic(err)
	}

	return credentials["apPasswordHash"]
}

func (dyson *Device) getDeviceTopic(command string) string {
	return fmt.Sprintf("%s/%s/%s", dyson.ProductType, dyson.Serial, command)
}

// Connect to the device's MQTT service
func (dyson *Device) Connect(ipAddress string) {
	password := dyson.decipherCredentials()
	broker := fmt.Sprintf("tcp://%s:1883", ipAddress)

	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetProtocolVersion(3)
	opts.SetClientID("dyson_" + string(rand.Intn(1000000000)))
	opts.SetUsername(dyson.Serial)
	opts.SetPassword(password)

	opts.SetDefaultPublishHandler(
		func(client mqtt.Client, msg mqtt.Message) {
			payload := msg.Payload()
			log.Println(string(payload))
		},
	)

	opts.SetOnConnectHandler(
		func(mclient mqtt.Client) {
			mclient.Subscribe(dyson.getDeviceTopic("status/current"), 0, nil)
		},
	)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	dyson.client = client
}
