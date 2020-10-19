package dyson

// Device is the spec of a registered Dyson Device on your account
type Device struct {
	Serial              string
	Name                string
	Version             string
	LocalCredentials    string
	ProductType         string
	ConnectionType      string
	AutoUpdate          bool
	NewVersionAvailable bool
}
