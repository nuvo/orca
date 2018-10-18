## Credentials

### Kubernetes

Orca tries to get credentials in the following order:
If `KUBECONFIG` environment variable is set - orca will use the current context from that config file. Otherwise it will use `~/.kube/config`.