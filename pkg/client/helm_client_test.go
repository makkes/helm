package client_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"helm.sh/helm/v3/pkg/client"
	"helm.sh/helm/v3/pkg/client/http"
	"helm.sh/helm/v3/pkg/client/oci"
)

func TestHelmClientsWithBasicAuth(t *testing.T) {

	ociUser := os.Getenv("OCI_USER")
	ociPassword := os.Getenv("OCI_PASSWORD")
	if ociUser == "" || ociPassword == "" {
		t.Fatalf("to run this test you need to set the OCI_USER and OCI_PASSWORD environment variables.")
	}

	httpUser := os.Getenv("HTTP_USER")
	httpPassword := os.Getenv("HTTP_PASSWORD")
	if ociUser == "" || ociPassword == "" {
		t.Fatalf("to run this test you need to set the HTTP_USER and HTTP_PASSWORD environment variables.")
	}

	tests := []struct {
		name     string
		repoURL  string
		bldrCtor func(repoURL string) client.Builder
		user     string
		password string
	}{
		{
			"successfully downloads a chart from an OCI registry",
			"ghcr.io/stefanprodan/charts",
			oci.NewOCIHelmClientBuilder,
			ociUser,
			ociPassword,
		},
		{
			"successfully downloads a chart from an HTTP repository",
			"http://localhost:8080",
			http.NewHTTPHelmClientBuilder,
			httpUser,
			httpPassword,
		},
	}

	chartName := "podinfo"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := tt.bldrCtor(tt.repoURL).
				WithAuthenticationMethodBasicAuth(tt.user, tt.password).
				Build()
			if err != nil {
				t.Fatalf("unexpected error building Helm client: %s", err)
			}

			versions, err := client.ListVersions(chartName)
			if err != nil {
				t.Fatalf("unexpected error listing versions: %s", err)
			}
			version := versions[len(versions)-1]

			chartBuf, err := client.GetChart(chartName, version)
			if err != nil {
				t.Fatalf("unexpected error getting chart: %s", err)
			}

			tmpDir := t.TempDir()
			chartFileName := filepath.Join(tmpDir, fmt.Sprintf("%s-%s.tar.gz", chartName, version))
			if err := ioutil.WriteFile(chartFileName, chartBuf.Bytes(), 0644); err != nil {
				t.Fatalf("unexpected error writing chart file: %s", err)
			}

		})
	}
}
