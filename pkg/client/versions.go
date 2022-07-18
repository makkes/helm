package client

import (
	"fmt"
	"sort"

	"github.com/Masterminds/semver/v3"
	"helm.sh/helm/v3/pkg/repo"
)

func SortSemver(in []string) ([]string, error) {
	versions := make([]*semver.Version, len(in))
	var err error

	for idx, v := range in {
		versions[idx], err = semver.NewVersion(v)
		if err != nil {
			return nil, fmt.Errorf("failed parsing semantic chart version '%s': %w", v, err)
		}
	}

	sort.Sort(semver.Collection(versions))

	res := make([]string, len(versions))
	for idx, v := range versions {
		res[idx] = v.String()
	}

	return res, nil
}

func SortChartVersions(v repo.ChartVersions) ([]string, error) {
	versions := make([]*semver.Version, len(v))
	var err error

	for idx, v := range v {
		versions[idx], err = semver.NewVersion(v.Version)
		if err != nil {
			return nil, fmt.Errorf("failed parsing semantic chart version '%s': %w", v.Version, err)
		}
	}

	sort.Sort(semver.Collection(versions))

	res := make([]string, len(versions))
	for idx, v := range versions {
		res[idx] = v.String()
	}

	return res, nil
}
