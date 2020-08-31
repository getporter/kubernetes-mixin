module github.com/deislabs/porter-kubernetes

go 1.13

require (
	get.porter.sh/porter v0.28.1
	github.com/Masterminds/semver v1.5.0
	github.com/docker/cnab-to-oci v0.3.0-beta3 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/gobuffalo/envy v1.9.0 // indirect
	github.com/gobuffalo/packr/v2 v2.8.0
	github.com/karrick/godirwalk v1.16.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rogpeppe/go-internal v1.6.1 // indirect
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.6.1
	github.com/xeipuuv/gojsonschema v1.2.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208 // indirect
	golang.org/x/sys v0.0.0-20200831180312-196b9ba8737a // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/hashicorp/go-plugin => github.com/carolynvs/go-plugin v1.0.1-acceptstdin
