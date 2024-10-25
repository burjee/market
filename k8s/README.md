# K8S

k8s configuration file

## Install Ingress NGINX Controller

[ingress nginx](https://github.com/kubernetes/ingress-nginx)

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.12.0-beta.0/deploy/static/provider/cloud/deploy.yaml
```

## Run

Before performing this operation, you need to create Docker images for `server_match`, `server_notify`, and `server_order`.

```bash
kubectl apply -f .
```

## System Logs

You can view the logs of the matching system, which display the average matches per second for the past one minute, five minutes, fifteen minutes, and overall.

```bash
kubectl logs --tail=1 -f {your-market-match-pod-name}
```

## Delete

```bash
kubectl delete -f .
```
