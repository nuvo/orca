# Orca

Orca is a simplifier. It focuses on the world around Kubernetes and CI\CD, but it is also handy in daily work. It takes complex tasks and makes them easy to accomplish.

## Install

### Prerequisits

1. git
2. [dep](https://github.com/golang/dep)
3. [helm](https://helm.sh/) (required for runtime)

### Install
```
wget -qO- https://github.com/maorfr/orca/releases/download/<TAG>/orca.tar.gz | sudo tar xvz -C /usr/local/bin

```

### Use it in your CI\CD processes
```
docker pull maorfr/orca
```
Docker hub repository: https://hub.docker.com/r/maorfr/orca

### Build from source

Orca uses dep as a dependency management tool.

```
go get -u github.com/golang/dep/cmd/dep
dep ensure
go build -o orca cmd/orca.go
```

## Commands

The following commands are available:
```
deploy chart            Deploy a Helm chart from chart repository
push chart              Push Helm chart to chart repository

get env                 Get list of Helm releases in an environment (Kubernetes namespace)
deploy env              Deploy a list of Helm charts to an environment (Kubernetes namespace)
delete env              Delete an environment (Kubernetes namespace) along with all Helm releases in it

create resource         Create or update a resource via REST API
get resource            Get a resource via REST API
delete resource         Delete a resource via REST API

determine buildtype     Determine build type based on path filters
```

## Why should you use Orca?

* If you want to create environments dynamically (as part of a Pull Request for example) -

```
# Get the "stable" environment (this could be production)
orca get env --name $SRC_RELEASE_NS --kube-context $SRC_KUBE_CONTEXT > charts.yaml
# Deploy the same configuration to a new namespace (you can override specific versions to have a great test environment)
orca deploy env --name $COMMIT_HASH -c charts.yaml --kube-context $DST_KUBE_CONTEXT --override $CHART_NAME=$CHART_TEST_VERSION
```

* If you want to achieve different CI processes depending on changed paths -

```
# Get the previous pipeline`s commit
orca get resource --url $PIPELINES_API_URL --headers "PRIVATE-TOKEN:$USER_TOKEN" --key sha --value $CI_COMMIT_SHA --offset 1 -p sha > previous_pipeline_commit
# Determine the build type based on path filters
orca determine buildtype --default-type code --path-filter ^config.*$=config --prev-commit $(cat previous_pipeline_commit) > .buildtype
```

* If you want to delete a resource through REST API by ID, but only have the name of the resource -
```
# Get the resource ID
ID=$(orca get resource --url $API_URL --headers "PRIVATE-TOKEN:$USER_TOKEN" --key name --value $NAME -p id)
# Delete the resouce
orca delete resource --url $API_URL/$ID --headers "PRIVATE-TOKEN:$USER_TOKEN"
```

See the Examples section for more examples!


## Examples

### Build type determination
This function provides the ability to execute different tasks on different branches and different changed paths. It is essentially a path filter implementation for tools which do not support path filters.

In this example, if files changed only in the `src` directory, the build type will be set to `code`, if files changed only in the `config` directory, the build type will be set to `config`. If files changed in both (or anywhere else), it will be set to `code`:
```
orca determine buildtype \
    --default-type code \
    --path-filter ^src.*$=code,^config.*$=config \
    --prev-commit <previousCommitHash>
```

In this example, if the current reference is different from the mainline and a release branch, the build type will be set to `default`:
```
orca determine buildtype \
    --default-type default \
    --curr-ref develop \
    --main-ref master \
    --rel-ref ^.*/rel-.*$
```

The two examples can be combined.


### Get resource
This function gets a resource from REST API.

In this example, it gets the previous commit hash (offset of 1 from the current commit hash), to be able to compare it to the current commit hash (which will help with the previous example):

```
orca get resource \
    --url <pipelinesURL> \
    --headers "<header>:<value>" \
    --key sha \
    --value <currentCommit> \
    --offset 1 \
    --print-key sha
```

### Get env
This functions gets all Helm installed releases from an environment (Kubernetes namespace).

In this example, only orca managed releases will be displayed (a managed release is considered one with release name in the form of namespace-chartName):

```
orca get env \
    --kube-context <kubeContext> \
    --name <namespace>
```

You can add the `--only-managed=false` to show all releases in a namespace.


### Deploy chart
This function deploys a Helm chart from a chart repository, using values files which are packed along with the chart.

In this example, the specified chart repository will be added, the chart will be fetched from it, and deployed using `prod-values.yaml` (packed within the chart) to the specified kubernetes context and namespace:

```
orca deploy chart \
    --name <chartName> \
    --version <chartVersion> \
    --release-name <releaseName> \
    --kube-context <kubeContext> \
    --namespace <namespace> \
    -f prod-values.yaml \
    --repo myrepo=<repoURL>
```

### Deploy env
This function deploys a list of Helm charts from a chart repository. This function supports runtime dependencies between charts.

If this is `charts.yaml`:
```
charts:
- name: cassandra
  version: 0.4.0
- name: mariadb
  version: 0.5.4
- name: serviceA
  version: 0.1.7
  depends_on:
  - cassandra
  - mariadb
- name: serviceB
  version: 0.2.3
  depends_on:
  - serviceA
```
Then the below line will deploy the charts in the following order (using topological sort algorithm):
1. cassandra, mariadb
2. serviceA
3. serviceB

```
orca deploy env \
    --name <namespace> \
    -c charts.yaml \
    --kube-context <kubeContext> \
    -f prod-values.yaml
```

## Environment variables support

Orca commands support the usage of environment variables instead of most of the flags. For example:
The `get env` command can be executed as mentioned in the example:
```
orca get env \
    --kube-context <kubeContext> \
    --name <namespace>
```

You can also set the appropriate envrionment variables (ORCA_FLAG, _ instead of -):

```
export ORCA_KUBE_CONTEXT=<kubeContext>
export ORCA_NAME=<namespace>

orca get env
```

## Credentials


### Kubernetes

Orca tries to get credentials in the following order:
If `KUBECONFIG` environment variable is set - orca will use the current context from that config file. Otherwise it will use `~/.kube/config`.

## All commands

### Deploy chart
```
Deploy a Helm chart from chart repository

Usage:
  orca deploy chart [flags]

Flags:
      --helm-tls-store string   path to TLS certs and keys. Overrides $HELM_TLS_STORE
      --inject                  enable injection during helm upgrade. Overrides $ORCA_INJECT (requires helm inject plugin: https://github.com/maorfr/helm-inject)
      --kube-context string     name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT
      --name string             name of chart to deploy. Overrides $ORCA_NAME
  -n, --namespace string        kubernetes namespace to deploy to. Overrides $ORCA_NAMESPACE
      --release-name string     release name. Overrides $ORCA_RELEASE_NAME
      --repo string             chart repository (name=url). Overrides $ORCA_REPO
  -s, --set strings             set additional parameters
      --tls                     enable TLS for request. Overrides $ORCA_TLS
  -f, --values strings          values file to use (packaged within the chart)
      --version string          version of chart to deploy. Overrides $ORCA_VERSION
```

### Push chart
```
Push Helm chart to chart repository (requires helm push plugin: https://github.com/chartmuseum/helm-push)

Usage:
  orca push chart [flags]

Flags:
      --append string   string to append to version. Overrides $ORCA_APPEND
      --lint            should perform lint before push. Overrides $ORCA_LINT
      --path string     path to chart. Overrides $ORCA_PATH
      --repo string     chart repository (name=url). Overrides $ORCA_REPO
```

### Get env
```
Get list of Helm releases in an environment (Kubernetes namespace)

Usage:
  orca get env [flags]

Flags:
      --helm-tls-store string   path to TLS certs and keys. Overrides $HELM_TLS_STORE
      --kube-context string     name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT
  -n, --name string             name of environment (namespace) to get. Overrides $ORCA_NAME
      --only-managed            list only releases managed by orca. Overrides $ORCA_ONLY_MANAGED (default true)
  -o, --output string           output format (yaml, md). Overrides $ORCA_OUTPUT
      --tls                     enable TLS for request. Overrides $ORCA_TLS
```

### Deploy env
```
Deploy a list of Helm charts to an environment (Kubernetes namespace)

Usage:
  orca deploy env [flags]

Aliases:
  env, environment

Flags:
  -c, --charts-file string                   path to file with list of Helm charts to install. Overrides $ORCA_CHARTS_FILE
  -x, --deploy-only-override-if-env-exists   if environment exists - deploy only override(s) (support for features spanning multiple services). Overrides $ORCA_DEPLOY_ONLY_OVERRIDE_IF_ENV_EXISTS
      --helm-tls-store string                path to TLS certs and keys. Overrides $HELM_TLS_STORE
      --inject                               enable injection during helm upgrade. Overrides $ORCA_INJECT (requires helm inject plugin: https://github.com/maorfr/helm-inject)
      --kube-context string                  name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT
  -n, --name string                          name of environment (namespace) to deploy to. Overrides $ORCA_NAME
      --override strings                     chart to override with different version (can specify multiple): chart=version
      --repo string                          chart repository (name=url). Overrides $ORCA_REPO
  -s, --set strings                          set additional parameters
      --tls                                  enable TLS for request. Overrides $ORCA_TLS
  -f, --values strings                       values file to use (packaged within the chart)
```

### Delete env
```
Delete an environment (Kubernetes namespace) along with all Helm releases in it

Usage:
  orca delete env [flags]

Flags:
      --force                   force environment deletion. Overrides $ORCA_FORCE
      --helm-tls-store string   path to TLS certs and keys. Overrides $HELM_TLS_STORE
      --kube-context string     name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT
  -n, --name string             name of environment (namespace) to delete. Overrides $ORCA_NAME
      --tls                     enable TLS for request. Overrides $ORCA_TLS
```

### Create resource
```
Create or update a resource via REST API

Usage:
  orca create resource [flags]

Flags:
      --headers strings   headers of the request (supports multiple)
      --method string     method to use in the request. Overrides $ORCA_METHOD (default "POST")
      --update            should method be PUT instead of POST. Overrides $ORCA_UPDATE
      --url string        url to send the request to. Overrides $ORCA_URL
```

### Get resource
```
Get a resource via REST API

Usage:
  orca get resource [flags]

Flags:
  -e, --error-indicator string   string indicating an error in the request. Overrides $ORCA_ERROR_INDICATOR (default "E")
      --headers strings          headers of the request (supports multiple)
      --key string               find the desired object according to this key. Overrides $ORCA_KEY
      --offset int               offset of the desired object from the reference key. Overrides $ORCA_OFFSET
  -p, --print-key string         key to print. If not specified - prints the response. Overrides $ORCA_PRINT_KEY
      --url string               url to send the request to. Overrides $ORCA_URL
      --value string             find the desired object according to to key`s value. Overrides $ORCA_VALUE
```

### Delete resource
```
Delete a resource via REST API

Usage:
  orca delete resource [flags]

Flags:
      --headers strings   headers of the request (supports multiple)
      --url string        url to send the request to. Overrides $ORCA_URL
```

### Determine buildtype
```
Determine build type based on path filters

Usage:
  orca determine buildtype [flags]

Flags:
      --allow-multiple-types       allow multiple build types. Overrides $ORCA_ALLOW_MULTIPLE_TYPES
      --curr-ref string            current reference name. Overrides $ORCA_CURR_REF
      --default-type string        default build type. Overrides $ORCA_DEFAULT_TYPE (default "default")
      --main-ref string            name of the reference which is the main line. Overrides $ORCA_MAIN_REF
      --path-filter strings        path filter (supports multiple) in the path=buildtype form (supports regex)
      --prev-commit string         previous commit for paths comparison. Overrides $ORCA_PREV_COMMIT
      --prev-commit-error string   identify an error with the previous commit by this string. Overrides $ORCA_PREV_COMMIT_ERROR (default "E")
      --rel-ref string             release reference name (or regex). Overrides $ORCA_REL_REF
```
