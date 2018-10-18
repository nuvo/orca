workflow "New workflow" {
  on = "push"
  resolves = ["Orca"]
}

action "Orca" {
  uses = "./.github/orca"
  env = {
    SRC_KUBE_CONTEXT="prod",
    SRC_NS="default",
    DST_KUBE_CONTEXT="dev",
    DST_NS="orca-1",
    CHART_NAME="example-chart",
    CHART_VERSION="0.1.0"
  }
}