# kubernetes Mixin for Porter

<img src="https://porter.sh/images/mixins/kubernetes.svg" align="right" width="150px"/>

This is a kubernetes mixin for [Porter](https://github.com/deislabs/porter). It executes the
appropriate helm command based on which action it is included within: `install`,
`upgrade`, or `delete`.

### Install or Upgrade

```shell
porter mixin install kubernetes
```

#### Install or Upgrade canary version

```shell
porter mixin install kubernetes --version canary --url https://cdn.porter.sh/mixins/kubernetes
```

#### Install or Upgrade from feed-url

```shell
porter mixin install kubernetes --feed-url https://cdn.porter.sh/mixins/atom.xml
```

#### Manually Install or Upgrade with a specific version from github

```shell
porter mixin install kubernetes --version $VERSION --url https://github.com/deislabs/porter-kubernetes/releases/download
```

### Mixin Configuration

#### Kubernetes client version

```yaml
- kubernetes:
    clientVersion: v1.15.5
```

### Mixin Actions Syntax

#### Install Action

```yaml
install:
  - kubernetes:
      description: "Install Hello World App"
      manifests:
        - /cnab/app/manifests/hello
      wait: true

```

#### Upgrade Action

```yaml
upgrade:
  - kubernetes:
      description: "Upgrade Hello World App"
      manifests:
        - /cnab/app/manifests/hello
      wait: true

```

#### Uninstall Action

```yaml
uninstall:
  - kubernetes:
      description: "Uninstall Hello World App"
      manifests:
        - /cnab/app/manifests/hello
      wait: true

```

#### Outputs

The mixin supports extracting resource metadata from Kubernetes as outputs.

```yaml
outputs:
    - name: NAME
      resourceType: RESOURCE_TYPE
      resourceName: RESOURCE_TYPE_NAME
      namespace: NAMESPACE
      jsonPath: JSON_PATH_DEFINITION
```
