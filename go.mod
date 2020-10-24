module get.porter.sh/mixin/kubernetes

go 1.13

require (
	get.porter.sh/porter v0.22.1-beta.1
	github.com/Masterminds/semver v1.5.0
	github.com/ghodss/yaml v1.0.0
	github.com/gobuffalo/packr/v2 v2.8.0
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.5.1
	github.com/xeipuuv/gojsonschema v1.2.0
	golang.org/x/mod v0.2.0
	gopkg.in/yaml.v2 v2.2.4
)

replace github.com/hashicorp/go-plugin => github.com/carolynvs/go-plugin v1.0.1-acceptstdin
