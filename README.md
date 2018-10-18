# Orca

Orca is an advanced CI\CD tool which focuses on the world around Kubernetes, Helm and CI\CD, and it is also handy in daily work.
Orca is a simplifier - It takes complex tasks and makes them easy to accomplish.
Is is important to note that Orca is not intended to replace Helm, but rather to empower it and enable advanced usage with simplicity.

## Install

### From a release

1. git
2. [dep](https://github.com/golang/dep)
3. [helm](https://helm.sh/) (required for runtime)

Download the latest release from the [Releases page](https://github.com/maorfr/orca/releases) or use it in your CI\CD process with a [Docker image](https://hub.docker.com/r/maorfr/orca)

### From source

```
git clone https://github.com/maorfr/orca.git
cd orca
make
```

## Why should you use Orca?

What Orca does best is manage environments. An Environment is a Kubernetes namespace with a set of Helm charts installed on it.
There are a few use cases you will probably find useful right off the bat.

### Create a dynamic environment (as part of a Pull Request for example)

#### Get the "stable" environment and deploy the same configuration to a new environment

This will deploy the "stable" configuration (production?) to a destination namespace.

```
orca get env --name $SRC_NS --kube-context $SRC_KUBE_CONTEXT > charts.yaml
orca deploy env --name $DST_NS -c charts.yaml --kube-context $DST_KUBE_CONTEXT
```

Additional flags:

Use the `-p` (parallel) flag to specify parallelism of chart deployments.

Use the `-f` flag to specify different values files to use during deployment.

Use the `-s` flag to set additional parameteres.

#### Get the "stable" environment and deploy the same configuration to a new environment, with override(s)

Useful for creating test environments for a single service.
This will deploy the "stable" configuration to a destination namespace, except for the specified override(s), which will be deployed with version `CHART_VERSION`.

```
orca get env --name $SRC_NS --kube-context $SRC_KUBE_CONTEXT > charts.yaml
orca deploy env --name $DST_NS -c charts.yaml --kube-context $DST_KUBE_CONTEXT \
    --override $CHART_NAME=$CHART_VERSION
```

#### Get the "stable" environment and deploy the same configuration to a new environment, with override(s) and existence check

Useful for creating test environments for multiple services. Handy for testing a single feature spanning across multiple services.
This will deploy the same configuration to a destination namespace, except for the specified override(s), which will be deployed with version CHART_VERSION. If the environment already exists, only the specified override(s) will be deployed (using the `-x` flag - deploy only override if environment exists).
The following commands will be a part of all CI\CD processes in all services:

```
orca get env --name $SRC_NS --kube-context $SRC_KUBE_CONTEXT > charts.yaml
orca deploy env --name $DST_NS -c charts.yaml --kube-context $DST_KUBE_CONTEXT \
    --override $CHART_NAME=$CHART_VERSION \
    -x
```

When the first service's process starts, it creates the environment and deploys the configuration from the "stable" environment (exactly the same as the previous example). When the Nth service's process starts, the environment already exists, and only the specified override(s) are deployed.
Orca also handles a potential race condition between 2 or more services by "locking" the environment during deployment (using a `busy` annotation on the namespace).

Using the `-x` flag, after deploying from (for example) 3 different repositories, the new environment will have the "stable" configuration, except for the 3 services which are currently under test, which will be deployed with their respective `CHART_VERSION`s.

You can add the `-x` flag even if this service is completely isolated (for consistency).

### Create and update static environments

#### Manage multiple versions of your product without constantly maintaining the CI\CD process for all services

If you are supporting more then one version of your product, you can use Orca as the CI\CD tool to deploy and update environments with different configurations with ease.
Assuming you are required to create a new environment of your product, create a new Git repository with a single `charts.yaml` file, which you can update and deploy as you need.

Your CI\CD process may be as slim as:

```
orca deploy env --name $NS -c charts.yaml --kube-context $KUBE_CONTEXT
```

### Keep track of an environment's state

This is a bonus! If you need to document changes in your environments, you can use Orca to accomplish it. Trigger an event of your choice whenever an environment is updated and use Orca to get the current state:

```
orca get env --name $SRC_NS --kube-context $SRC_KUBE_CONTEXT -o md
```

This will print the list of currently installed charts in Mardown format.

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

## Docs

### Commands

Since Orca is a tool designed for CI\CD, it has additional commands and options to help with common actions:
```
deploy chart            Deploy a Helm chart from chart repository
push chart              Push Helm chart to chart repository
get env                 Get list of Helm releases in an environment (Kubernetes namespace)
deploy env              Deploy a list of Helm charts to an environment (Kubernetes namespace)
delete env              Delete an environment (Kubernetes namespace) along with all Helm releases in it
lock env                Lock an environment (Kubernetes namespace)
unlock env              Unlock an environment (Kubernetes namespace)
create resource         Create or update a resource via REST API
get resource            Get a resource via REST API
delete resource         Delete a resource via REST API
determine buildtype     Determine build type based on path filters
```

For a more detailed description of all commands, see the [Commands](https://github.com/maorfr/orca/tree/master/docs/commands) section

## Examples

Be sure to check out the [Examples](https://github.com/maorfr/orca/tree/master/docs/examples) section!
