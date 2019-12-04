<img src="/docs/img/logo.png" width="250px" alt="orca logo">


[![Release](https://img.shields.io/github/release/nuvo/orca.svg)](https://github.com/nuvo/orca/releases)
[![Travis branch](https://img.shields.io/travis/nuvo/orca/master.svg)](https://travis-ci.org/nuvo/orca)
[![Docker Pulls](https://img.shields.io/docker/pulls/nuvo/orca.svg)](https://hub.docker.com/r/nuvo/orca/)
[![Go Report Card](https://goreportcard.com/badge/github.com/nuvo/orca)](https://goreportcard.com/report/github.com/nuvo/orca)
[![license](https://img.shields.io/github/license/nuvo/orca.svg)](https://github.com/nuvo/orca/blob/master/LICENSE)

# Orca

Orca is an advanced CI\CD tool which focuses on the world around Kubernetes, Helm and CI\CD, and it is also handy in daily work.
Orca is a simplifier - It takes complex tasks and makes them easy to accomplish.
Is is important to note that Orca is not intended to replace Helm, but rather to empower it and enable advanced usage with simplicity.

## Install

### Prerequisites

1. git
2. [dep](https://github.com/golang/dep)
3. [Helm](https://helm.sh/) (required for `env` and `chart` subcommands)
4. [ChartMuseum](https://github.com/helm/charts/tree/master/stable/chartmuseum) or any other chart repository implementation (required for `deploy` commands)

### From a release

Download the latest release from the [Releases page](https://github.com/nuvo/orca/releases) or use it in your CI\CD process with a [Docker image](https://hub.docker.com/r/nuvo/orca)

### From source

```
mkdir -p $GOPATH/src/github.com/nuvo && cd $_
git clone https://github.com/nuvo/orca.git && cd orca
make
```

* Please note that the `master` branch and the `latest` docker image are unstable. For a stable and tested version, use a tagged release or a tagged docker image.

## Why should you use Orca?

What Orca does best is manage environments. An Environment is a Kubernetes namespace with a set of Helm charts installed on it.
There are a few use cases you will probably find useful right off the bat.

### Create a dynamic environment (as part of a Pull Request for example)

#### Get the "stable" environment and deploy the same configuration to a new environment

This will deploy the "stable" configuration (production?) to a destination namespace.

```
orca get env --name $SRC_NS --kube-context $SRC_KUBE_CONTEXT > charts.yaml
orca deploy env --name $DST_NS -c charts.yaml \
    --kube-context $DST_KUBE_CONTEXT \
    --repo myrepo=$REPO_URL
```

Additional flags:

* Use the `-p` (parallel) flag to specify parallelism of chart deployments.
* Use the `-f` flag to specify different values files to use during deployment.
* Use the `-s` flag to set additional parameters.

#### Get the "stable" environment and deploy the same configuration to a new environment, with override(s) and environment refresh

Useful for creating test environments for a single service or for multiple services.
Handy for testing a single feature spanning across one or more services.
This will deploy the "stable" configuration to a destination namespace, except for the specified override(s), which will be deployed with version `CHART_VERSION`. In addition, it will set an annotation stating that `CHART_NAME` is protected and can only be overridden by itself.
If the reference environment has changed between environment deployments, the new environment will be updated with these changes.   

The following commands will be a part of all CI\CD processes in all services:

```
orca get env --name $SRC_NS --kube-context $SRC_KUBE_CONTEXT > charts.yaml
orca deploy env --name $DST_NS -c charts.yaml \
    --kube-context $DST_KUBE_CONTEXT \
    --repo myrepo=$REPO_URL \
    --override $CHART_NAME=$CHART_VERSION \
    --protected-chart $CHART_NAME
```

* When the first service's process starts, it creates the environment and deploys the configuration from the "stable" environment (exactly the same as the previous example).
* When the Nth service's process starts, the environment already exists, and a previous chart that was deployed is `protected`. The current service will also be marked as `protected`, and will update the environment, without changing the previous protected service(s).
* After deploying from (for example) 3 different repositories, the new environment will have the latest "stable" configuration, except for the 3 services which are currently under test, which will be deployed with their respective `CHART_VERSION`s (protected by `--protected-chart`)
* You can add the `--protected-chart` flag even if this service is completely isolated (for consistency).
* Orca also handles a potential race condition between 2 or more services by "locking" the environment during deployment (using a `busy` annotation on the namespace).

#### Get the "stable" environment and deploy the same configuration to a new environment, with override(s) and without environment refresh

Useful for creating test environments for a single service or for multiple services.
Handy for testing a single feature spanning across one or more services, when you want to prevent updates from the reference environment after the initial creation of the new environment.
This will deploy the "stable" configuration to a destination namespace, except for the specified override(s), which will be deployed with version `CHART_VERSION`.
If the reference environment has changed between environment deployments, the new environment will NOT be updated with these changes.

The following commands will be a part of all CI\CD processes in all services:

```
orca get env --name $SRC_NS --kube-context $SRC_KUBE_CONTEXT > charts.yaml
orca deploy env --name $DST_NS -c charts.yaml \
    --kube-context $DST_KUBE_CONTEXT \
    --repo myrepo=$REPO_URL \
    --override $CHART_NAME=$CHART_VERSION \
    -x
```

* When the first service's process starts, it creates the environment and deploys the configuration from the "stable" environment (exactly the same as the previous example).
* When the Nth service's process starts, the environment already exists, so only the specified override(s) will be deployed.
* Assuming each process deploys a different chart (or charts), there is no need to protect them using the `--protected-chart`.


### Create and update static environments

#### Manage multiple versions of your product without constantly maintaining the CI\CD process for all services

If you are supporting more then one version of your product, you can use Orca as the CI\CD tool to deploy and update environments with different configurations with ease.
Assuming you are required to create a new environment of your product, create a new Git repository with a single `charts.yaml` file, which you can update and deploy as you need.

Your CI\CD process may be as slim as:

```
orca deploy env --name $NS -c charts.yaml \
    --kube-context $KUBE_CONTEXT \
    --repo myrepo=$REPO_URL
```

### Keep track of an environment's state

This is a bonus! If you need to document changes in your environments, you can use Orca to accomplish it. Trigger an event of your choice whenever an environment is updated and use Orca to get the current state:

```
orca get env --name $SRC_NS --kube-context $SRC_KUBE_CONTEXT -o md
```

This will print the list of currently installed charts in Markdown format.

### Prepare for disaster recovery

You can use Orca to prepare for a rainy day. Trigger an event of your choice whenever an environment is updated and use Orca to get the current state into a file (ideally keep it under source control):

```
orca get env --name $NS --kube-context $KUBE_CONTEXT -o yaml > charts.yaml
```

In case of emergency, you can deploy the same configuration using the `deploy env` command as explained above.

## Environment variables support

Orca commands support the usage of environment variables instead of most of the flags. For example:
The `get env` command can be executed as mentioned in the example:
```
orca get env \
    --kube-context <kubeContext> \
    --name <namespace>
```

You can also set the appropriate environment variables (ORCA_FLAG, _ instead of -):

```
export ORCA_KUBE_CONTEXT=<kubeContext>
export ORCA_NAME=<namespace>

orca get env
```

## Docs

### Commands

Since Orca is a tool designed for CI\CD, it has additional commands and options to help with common actions:
```
deploy artifact         Deploy an artifact to Artifactory
deploy chart            Deploy a Helm chart from chart repository
push chart              Push Helm chart to chart repository
get env                 Get list of Helm releases in an environment (Kubernetes namespace)
deploy env              Deploy a list of Helm charts to an environment (Kubernetes namespace) from chart repository
delete env              Delete an environment (Kubernetes namespace) along with all Helm releases in it
diff env                Show differences in Helm releases between environments (Kubernetes namespace)
lock env                Lock an environment (Kubernetes namespace)
unlock env              Unlock an environment (Kubernetes namespace)
validate env            Validate an environment (Kubernetes namespace)
create resource         Create or update a resource via REST API
get resource            Get a resource via REST API
delete resource         Delete a resource via REST API
determine buildtype     Determine build type based on path filters
```

For a more detailed description of all commands, see the [Commands](/docs/commands) section

## Examples

Be sure to check out the [Examples](/docs/examples) section!

## References

* [The Nuvo Group CI\CD journey](https://medium.com/nuvo-group-tech/the-nuvo-group-ci-cd-journey-132ab70bf452)
* [Dynamic Kubernetes Environments with Orca](https://medium.com/nuvo-group-tech/dynamic-kubernetes-environments-with-orca-a288afd36aa8)
