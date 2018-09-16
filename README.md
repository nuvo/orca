# Orca

Orca is a simplifier. It focuses on the world around Kubernetes and CI\CD, but it is also handy in daily work. It takes complex tasks and makes them easy to accomplish.

## Install

### Prerequisits

1. git
2. [glide](https://github.com/Masterminds/glide)
3. [helm](https://helm.sh/) (required for runtime)

### Install from source

Orca uses glide as a dependency management tool, since some of the referenced packages are not available using [dep](https://github.com/golang/dep).

```
glide up -v
go build -o orca cmd/orca.go
```

## Commands

The following commands are available:
```
delete env              Delete an environment (Kubernetes namespace) along with all Helm releases in it
delete resource         Delete a resource from REST API
deploy chart            Deploy a Helm chart from chart repository
deploy env              Deploy a list of Helm charts to an environment (Kubernetes namespace)
determine buildtype     Determine build type based on path filters
get env                 Get list of Helm releases in an environment (Kubernetes namespace)
get resource            Get a resource from REST API
push chart              Push Helm chart to chart repository
```

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