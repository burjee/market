# SIMPLE MARKET

A project for practicing microservices and Kubernetes, featuring a simple trade matching system. With 0.5 CPU and 128MB of memory, it can match 1000+ orders per second.

Users create orders, which are sent to the matching system via Kafka for processing. The notification system broadcasts the order list using WebSocket. Orders and other data are stored only in Redis.

## Folder

``` bash
├── k8s           # K8s Configuration File
├── market_test   # Market Test
├── server_match  # Order Match
├── server_notify # Order List
├── server_order  # Api
└── web           # Frontend page
```

## Run

#### build docker image

- server_match
``` bash
cd server_match
docker build . -t market/match
```

- server_notify
``` bash
cd server_notify
docker build . -t market/notify
```

- server_order
``` bash
cd server_order
docker build . -t market/order
```

#### Install Ingress NGINX Controller

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.12.0-beta.0/deploy/static/provider/cloud/deploy.yaml
```

#### Run

```bash
cd k8s
kubectl apply -f .
``` 

#### Market Test

This will continuously generate a large volume of random orders in a short period.

```bash
cd market_test
npm install -g artillery@latest
npx artillery run test.yaml
``` 

#### Match Metrics

You can view the logs of the matching system, which display the average matches per second for the past one minute, five minutes, fifteen minutes, and overall.

```bash
kubectl logs --tail=1 -f {your-market-match-pod-name}
```

#### Delete
```bash
cd k8s
kubectl delete -f .
``` 
