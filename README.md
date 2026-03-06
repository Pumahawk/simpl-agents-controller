```bash
function mise_remote() {
  local pod="$(kubectl get pod -oname -n default --sort-by={.metadata.creationTimestamp} | grep mise-deployer | tail -n 1)"
  kubectl exec -n default -it $pod -- bash -c 'env -C /app PATH="$PATH:/root/.local/bin" mise "$@"' -- "$@"
}
```
