package target

import (
	"net/http"
)

// DefaultRoundTripper is our default RoundTripper.
var DefaultRoundTripper http.RoundTripper = &http.Transport{
	Proxy:             http.ProxyFromEnvironment,
	DisableKeepAlives: true,
}

// BasicAuthTransport is the struct to add basic auth to a RoundTripper.
type BasicAuthTransport struct {
	Transport http.RoundTripper
	Username  string
	Password  string
}

// RoundTrip implements the RoundTrip for our RoundTripper with basic auth
// support.
func (bat BasicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(bat.Username, bat.Password)
	return bat.Transport.RoundTrip(req)
}

// TokenAuthTransporter is the struct to add token auth to a RoundTripper.
type TokenAuthTransporter struct {
	Transport http.RoundTripper
	Token     string
}

// RoundTrip implements the RoundTrip for our RoundTripper with token auth
// support.
func (tat TokenAuthTransporter) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+tat.Token)
	return tat.Transport.RoundTrip(req)
}
