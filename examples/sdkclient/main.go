package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"helm.sh/helm/v3/pkg/client"
	"helm.sh/helm/v3/pkg/client/http"
	"helm.sh/helm/v3/pkg/client/oci"
)

func buildHelmClient(repoURL *url.URL) (client.HelmClient, error) {
	switch repoURL.Scheme {
	case "oci":
		return oci.NewOCIHelmClientBuilder(strings.TrimPrefix(repoURL.String(), "oci://")).
			Build()
	case "https":
		return http.NewHTTPHelmClientBuilder(repoURL.String()).
			Build()
	default:
		return nil, fmt.Errorf("unsupported URL scheme %s. Only 'oci' and 'https' are supported", repoURL.Scheme)
	}

}

func main() {
	if len(os.Args) < 3 {
		panic(fmt.Sprintf("usage: %s REPO_URL CHART_NAME", os.Args[0]))
	}

	repoURL, err := url.Parse(os.Args[1])
	if err != nil {
		panic(err)
	}
	chartName := os.Args[2]

	fmt.Printf("fetching latest chart version of chart %s from %s\n", repoURL, chartName)

	client, err := buildHelmClient(repoURL)
	if err != nil {
		panic(err)
	}

	versions, err := client.ListVersions(chartName)
	if err != nil {
		panic(err)
	}
	version := versions[len(versions)-1]

	chartBuf, err := client.GetChart(chartName, version)
	if err != nil {
		panic(err)
	}

	chartFileName := fmt.Sprintf("%s-%s.tar.gz", chartName, version)
	if err := ioutil.WriteFile(chartFileName, chartBuf.Bytes(), 0644); err != nil {
		panic(err)
	}

	fmt.Printf("wrote chart to %s\n", chartFileName)
}
