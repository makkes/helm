package oci

import (
	"bytes"
	"fmt"
	"io"

	"helm.sh/helm/v3/pkg/client"
	"helm.sh/helm/v3/pkg/registry"
)

type OCIHelmClient struct {
	client.HelmClientBase
	registryURL    string
	registryClient *registry.Client
	loginOpts      []registry.LoginOption
}

var _ client.HelmClient = &OCIHelmClient{}

type OCIHelmClientBuilder struct {
	c *OCIHelmClient
}

var _ client.Builder = &OCIHelmClientBuilder{}

func NewOCIHelmClientBuilder(registryURL string) client.Builder {
	return &OCIHelmClientBuilder{
		c: &OCIHelmClient{
			registryURL: registryURL,
		},
	}
}

func (b *OCIHelmClientBuilder) WithAuthenticationMethodClientCerts(certFile string, keyFile string, caFile string) client.Builder {
	b.c.CertFile = certFile
	b.c.KeyFile = keyFile
	b.c.CaFile = caFile
	return b
}

func (b *OCIHelmClientBuilder) WithInsecureSkipVerifyTLS(insecureSkipVerifyTLS bool) client.Builder {
	b.c.InsecureSkipVerifyTLS = insecureSkipVerifyTLS
	return b
}

func (b *OCIHelmClientBuilder) WithAuthenticationMethodBasicAuth(username, password string) client.Builder {
	b.c.loginOpts = []registry.LoginOption{registry.LoginOptBasicAuth(username, password)}
	return b
}

func (b *OCIHelmClientBuilder) Build() (client.HelmClient, error) {
	regClient, err := registry.NewClient(
		registry.ClientOptWriter(io.Discard),
		registry.ClientOptCredentialsFile("/tmp/oci-creds"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed creating client: %w", err)
	}
	b.c.registryClient = regClient

	if err := regClient.Login(b.c.registryURL, b.c.loginOpts...); err != nil {
		return nil, fmt.Errorf("failed logging into registry: %w", err)
	}

	return b.c, nil
}

func (c OCIHelmClient) GetChart(chartName string, chartVersion string) (*bytes.Buffer, error) {
	res, err := c.registryClient.Pull(fmt.Sprintf("%s/%s:%s", c.registryURL, chartName, chartVersion))
	if err != nil {
		return nil, fmt.Errorf("failed pulling chart: %w", err)
	}
	if res.Chart == nil {
		return nil, fmt.Errorf("pull resulted in empty chart data")
	}
	return bytes.NewBuffer(res.Chart.Data), nil
}

func (c OCIHelmClient) ListVersions(chartName string) ([]string, error) {
	tags, err := c.registryClient.Tags(fmt.Sprintf("%s/%s", c.registryURL, chartName))
	if err != nil {
		return nil, fmt.Errorf("failed getting tags: %w", err)
	}
	return client.SortSemver(tags)
}

func (c OCIHelmClient) Login() {}
