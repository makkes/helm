package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"helm.sh/helm/v3/pkg/client"
	"helm.sh/helm/v3/pkg/repo"
)

type HTTPHelmClient struct {
	client.HelmClientBase
	repoURL    string
	idxEntries map[string]repo.ChartVersions
}

var _ client.HelmClient = &HTTPHelmClient{}

type HTTPHelmClientBuilder struct {
	c *HTTPHelmClient
}

var _ client.Builder = &HTTPHelmClientBuilder{}

func NewHTTPHelmClientBuilder(repoURL string) *HTTPHelmClientBuilder {
	return &HTTPHelmClientBuilder{
		c: &HTTPHelmClient{
			repoURL: repoURL,
		},
	}
}

func (b *HTTPHelmClientBuilder) WithAuthenticationMethodClientCerts(certFile string, keyFile string, caFile string) client.Builder {
	b.c.CertFile = certFile
	b.c.KeyFile = keyFile
	b.c.CaFile = caFile
	return b
}

func (b *HTTPHelmClientBuilder) WithInsecureSkipVerifyTLS(insecureSkipVerifyTLS bool) client.Builder {
	b.c.InsecureSkipVerifyTLS = insecureSkipVerifyTLS
	return b
}

func (b *HTTPHelmClientBuilder) Build() (client.HelmClient, error) {
	return b.c, nil
}

func (c HTTPHelmClient) GetChart(chartName string, chartVersion string) (*bytes.Buffer, error) {
	if c.idxEntries == nil {
		if err := c.loadIndexFile(); err != nil {
			return nil, fmt.Errorf("failed loading index file: %w", err)
		}
	}
	idxVersions, ok := c.idxEntries[chartName]
	if !ok {
		return nil, fmt.Errorf("no version of chart %s found in index file", chartName)
	}

	for _, v := range idxVersions {
		if v.Version == chartVersion {
			resp, err := http.Get(v.URLs[0])
			if err != nil {
				return nil, fmt.Errorf("failed downloading chart file: %w", err)
			}

			var b bytes.Buffer
			n, err := io.Copy(&b, resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed reading chart file from server: %w", err)
			}
			if resp.ContentLength != -1 && n != resp.ContentLength {
				return nil, fmt.Errorf("failed writing chart file data, %d bytes written, %d expected", n, resp.ContentLength)
			}
			return &b, nil
		}
	}

	return nil, fmt.Errorf("chart %s, version %s not found in index", chartName, chartVersion)
}

func (c *HTTPHelmClient) ListVersions(chartName string) ([]string, error) {
	if c.idxEntries == nil {
		if err := c.loadIndexFile(); err != nil {
			return nil, fmt.Errorf("failed loading index file: %w", err)
		}
	}
	idxVersions, ok := c.idxEntries[chartName]
	if !ok {
		return nil, fmt.Errorf("no version of chart %s found in index file", chartName)
	}

	return client.SortChartVersions(idxVersions)
}

func (c *HTTPHelmClient) loadIndexFile() error {
	file, err := os.CreateTemp(os.TempDir(), "")
	if err != nil {
		return fmt.Errorf("failed creating temporary local index file: %w", err)
	}
	defer os.Remove(file.Name())

	resp, err := http.Get(fmt.Sprintf("%s/index.yaml", c.repoURL))
	if err != nil {
		return fmt.Errorf("failed downloading index file: %w", err)
	}

	n, err := io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed writing index file: %w", err)
	}
	if resp.ContentLength != -1 && n != resp.ContentLength {
		return fmt.Errorf("failed writing index file, %d bytes written, %d expected", n, resp.ContentLength)
	}

	idxFile, err := repo.LoadIndexFile(file.Name())
	if err != nil {
		return fmt.Errorf("failed loading index file: %w", err)
	}

	c.idxEntries = idxFile.Entries

	return nil
}

func (c HTTPHelmClient) Login() {}
