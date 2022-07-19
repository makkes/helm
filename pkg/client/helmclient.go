package client

import "bytes"

type HelmClient interface {
	GetChart(chartName, chartVersion string) (*bytes.Buffer, error)
	ListVersions(chartName string) ([]string, error)
	Login()
}

type HelmClientBase struct {
	CertFile              string
	KeyFile               string
	CaFile                string
	InsecureSkipVerifyTLS bool
}

type Builder interface {
	WithAuthenticationMethodBasicAuth(username, password string) Builder
	WithAuthenticationMethodClientCerts(certFile, keyFile, caFile string) Builder
	WithInsecureSkipVerifyTLS(insecureSkipVerifyTLS bool) Builder
	Build() (HelmClient, error)
}
