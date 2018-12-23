## All commands


### Deploy artifact
```
Deploy an artifact to Artifactory

Usage:
  orca deploy artifact [flags]

Flags:
      --file string       path to file to deploy. Overrides $ORCA_FILE
      --token string      token to use for deployment. Overrides $ORCA_TOKEN
      --url string        url to deploy to. Overrides $ORCA_URL
```

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
      --timeout int             time in seconds to wait for any individual Kubernetes operation (like Jobs for hooks). Overrides $ORCA_TIMEOUT (default 300)
      --tls                     enable TLS for request. Overrides $ORCA_TLS
      --validate                perform environment validation after deployment. Overrides $ORCA_VALIDATE
  -f, --values strings          values file to use (packaged within the chart)
      --version string          version of chart to deploy. Overrides $ORCA_VERSION
```

`helm-tls-store` - path to directory containing `<kube-context>.cert.pem` and `<kube-context>.key.pem` files

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
      --kube-context string   name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT
  -n, --name string           name of environment (namespace) to get. Overrides $ORCA_NAME
  -o, --output string         output format (yaml, md). Overrides $ORCA_OUTPUT
```

### Deploy env
```
Deploy a list of Helm charts to an environment (Kubernetes namespace)

Usage:
  orca deploy env [flags]

Aliases:
  env, environment

Flags:
      --annotations strings                  additional environment (namespace) annotations (can specify multiple): annotation=value
  -c, --charts-file string                   path to file with list of Helm charts to install. Overrides $ORCA_CHARTS_FILE
  -x, --deploy-only-override-if-env-exists   if environment exists - deploy only override(s) (avoid environment update). Overrides $ORCA_DEPLOY_ONLY_OVERRIDE_IF_ENV_EXISTS
      --helm-tls-store string                path to TLS certs and keys. Overrides $HELM_TLS_STORE
      --inject                               enable injection during helm upgrade. Overrides $ORCA_INJECT (requires helm inject plugin: https://github.com/maorfr/helm-inject)
      --kube-context string                  name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT
      --labels strings                       environment (namespace) labels (can specify multiple): label=value
  -n, --name string                          name of environment (namespace) to deploy to. Overrides $ORCA_NAME
      --override strings                     chart to override with different version (can specify multiple): chart=version
  -p, --parallel int                         number of releases to act on in parallel. set this flag to 0 for full parallelism. Overrides $ORCA_PARALLEL (default 1)
      --protected-chart strings              chart name to protect from being overridden (can specify multiple)
      --repo string                          chart repository (name=url). Overrides $ORCA_REPO
  -s, --set strings                          set additional parameters
      --timeout int                          time in seconds to wait for any individual Kubernetes operation (like Jobs for hooks). Overrides $ORCA_TIMEOUT (default 300)
      --tls                                  enable TLS for request. Overrides $ORCA_TLS
      --validate                             perform environment validation after deployment. Overrides $ORCA_VALIDATE
  -f, --values strings                       values file to use (packaged within the chart)
```

`helm-tls-store` - path to directory containing `<kube-context>.cert.pem` and `<kube-context>.key.pem` files

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
  -p, --parallel int            number of releases to act on in parallel. set this flag to 0 for full parallelism. Overrides $ORCA_PARALLEL (default 1)
      --timeout int             time in seconds to wait for any individual Kubernetes operation (like Jobs for hooks). Overrides $ORCA_TIMEOUT (default 300)
      --tls                     enable TLS for request. Overrides $ORCA_TLS
```

`helm-tls-store` - path to directory containing `<kube-context>.cert.pem` and `<kube-context>.key.pem` files

### Diff env
```
Show differences in Helm releases between environments (Kubernetes namespace)

Usage:
  orca diff env [flags]

Flags:
  -h, --help                        help for env
      --kube-context-left string    name of the left kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT_LEFT
      --kube-context-right string   name of the right kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT_RIGHT
      --name-left string            name of left environment to compare. Overrides $ORCA_NAME_LEFT
      --name-right string           name of right environment to compare. Overrides $ORCA_NAME_RIGHT
```

### Lock env
```
Lock an environment (Kubernetes namespace)

Usage:
  orca lock env [flags]

Flags:
      --kube-context string   name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT
  -n, --name string           name of environment (namespace) to delete. Overrides $ORCA_NAME
```

### Unlock env
```
Unlock an environment (Kubernetes namespace)

Usage:
  orca unlock env [flags]

Flags:
      --kube-context string   name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT
  -n, --name string           name of environment (namespace) to delete. Overrides $ORCA_NAME
```

### Validate env
```
Validate an environment (Kubernetes namespace)

Usage:
  orca validate env [flags]

Flags:
      --kube-context string   name of the kubeconfig context to use. Overrides $ORCA_KUBE_CONTEXT
  -n, --name string           name of environment (namespace) to delete. Overrides $ORCA_NAME
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
