module get.porter.sh/mixin/kubernetes

go 1.16

replace github.com/hashicorp/go-plugin => github.com/carolynvs/go-plugin v1.0.1-acceptstdin

require (
	get.porter.sh/porter v0.38.7
	github.com/Masterminds/semver v1.5.0
	github.com/ghodss/yaml v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.7.0
	github.com/xeipuuv/gojsonschema v1.2.0
	gopkg.in/yaml.v2 v2.4.0
)
