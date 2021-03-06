### _module 02_

## 1 Go & Git

```shell
go build 1-main.go
sudo ./1-main

git add 1-main.go
git commit -m "xxx"
git push
```

### _module 03_

## 2 Docker Build & Push
```shell
docker build . -t lslsp2/httpserver:0.8 -f 2-Dockerfile
docker login
docker push lslsp2/httpserver:0.8
docker run -p 80:80 lslsp2/httpserver:0.8
```

### _module 08 (1/2)_

## 3.Create configmap

```shell
kubectl apply configmap lslsp-config --from-env-file=lslsp-config.env
```or
kubectl apply -f 3-lslsp-cofig.yaml
```

## 4.Create httpserver deployment (with 5 pods)

```shell
kubectl apply -f 4-deploy.yaml
```

### _module 08 (2/2)_

## 5.Create lslsp-httpserver service

```shell
kubectl apply -f 5-service.yaml
```

## 6.Create lslsp-tls secret

```shell
kubectl apply -f 6-secret.yaml
```

## 7.install nginx-ingress ( with jessie_meng's yaml :-> )
```shell
kubectl apply -f 7-ng-ingress-deploy-cp-from-cncamp.yaml
```

## 8.Create lslsp-httpserver-ingress ingress

```shell
kubectl apply -f 8-ingress.yaml
kubectl get svc -n ingress-nginx
curl -H "Host: lslsp.com" https://<ingress-nginx-controller CLUSTER-IP> -v -k
```

### _module 10_

## 9. Go & Git ( override step 1 with Metrics.go )

```shell
go mod init
go mod tidy
go build 1-main.go
sudo ./1-main

#(1/5)
git add 1-main.go
#(2/5)
git add metrics/metrics.go
git commit -m "xxx"
git push
```

## 10. Docker build with go.mod & Push ( override step 2 )
```shell
docker build . -t lslsp2/httpserver:0.10 -f 2-Dockerfile
docker login
docker push lslsp2/httpserver:0.10
docker run -p 80:80 lslsp2/httpserver:0.10
```

## 11.Update httpserver deployment ( override step 4-5 )

```shell
kubectl apply -f 4-deploy.yaml
kubectl apply -f 5-service.yaml
```

## 12.Install loki-stack ( Loki Grafana Promethus ) 

```sh
helm repo add grafana https://grafana.github.io/helm-charts
helm upgrade --install loki grafana/loki-stack --set grafana.enabled=true,prometheus.enabled=true,prometheus.alertmanager.persistentVolume.enabled=false,prometheus.server.persistentVolume.enabled=false
```

```sh
kubectl get svc
NAME                            TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                        AGE
jenkins                         NodePort    10.10.45.86     <none>        80:31650/TCP,50000:31623/TCP   31h
kubernetes                      ClusterIP   10.10.0.1       <none>        443/TCP                        5d2h
loki                            ClusterIP   10.10.143.255   <none>        3100/TCP                       33h
loki-grafana                    NodePort    10.10.117.223   <none>        80:32495/TCP                   33h
loki-headless                   ClusterIP   None            <none>        3100/TCP                       33h
loki-kube-state-metrics         ClusterIP   10.10.150.60    <none>        8080/TCP                       33h
loki-prometheus-alertmanager    ClusterIP   10.10.121.165   <none>        80/TCP                         33h
loki-prometheus-node-exporter   ClusterIP   None            <none>        9100/TCP                       33h
loki-prometheus-pushgateway     ClusterIP   10.10.173.203   <none>        9091/TCP                       33h
loki-prometheus-server          NodePort    10.10.79.112    <none>        80:31286/TCP                   33h
lslsp-httpserver                NodePort    10.10.167.152   <none>        80:30000/TCP                   5d1h
```

```sh
#(3/5)
kubectl edit svc loki-prometheus-server
# type: NodePort
kubectl edit svc loki-grafana
# type: NodePort
```
![image](screenshot/httpserver.png)


## 13. Query Metrics in Prometheus

```sh
#(4/5)
histogram_quantile(0.55, sum by(le) (rate(httpserver_execution_latency_seconds_bucket[5m])))
-> Execute
```
![image](screenshot/Prometheus.png)


## 14. Inport dashboard in Grafana

```sh
#(5/5)
Create -> Import -> Upload JSON file -> 14-httpserver-latency.json -> open
-> Import
```
![image](screenshot/Grafana.png)


### _module 12_

## 15. Go & Git (Transform httpserver into 3 microservices: service0 -> service1 -> service2 )

```shell
cd service0 (or service1, service2)
go mod init
go mod tidy
go build
sudo ./service0 (or service1, service2)

git add service0
git add service1
git add service2
git commit -m "xxx"
git push
```

## 16. Docker Build & Push
```shell
//build
cd service0
docker build . -t lslsp2/service0:0.12 -f ../dockerfile/16.1-Dockerfile
cd service1
docker build . -t lslsp2/service1:0.12 -f ../dockerfile/16.2-Dockerfile
cd service2
docker build . -t lslsp2/service2:0.12 -f ../dockerfile/16.3-Dockerfile

//push
docker login
docker push lslsp2/service0:0.12
docker push lslsp2/service1:0.12
docker push lslsp2/service2:0.12

//unit test
docker run -p 8080:80 lslsp2/service0:0.12
docker run -p 8081:80 lslsp2/service1:0.12
docker run -p 8082:80 lslsp2/service2:0.12

//service0's log:
2022/05/15 12:22:07 HTTP get failed with error: %!(EXTRA string=error, *url.Error=Get "http://service1": dial tcp: lookup service1 on 8.8.8.8:53: no such host

//remove
docker rm -f <container-id>
```

## 17. Install Istio
```shell
curl -L https://istio.io/downloadIstio | sh -
cd istio-1.13.3
export PATH=$PWD/bin:$PATH
istioctl install --set profile=demo -y
```

## 18. Deploy Secure Istio Ingress Gateway (1/3)
```shell
//secret
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=lslsp Inc./CN=*.lslsp.io' -keyout lslsp.io.key -out lslsp.io.crt

//deploy with sidecar
kubectl create ns istio-httpserver
kubectl label ns istio-httpserver istio-injection=enabled
kubectl apply -f 18-deployment.yaml -n istio-httpserver
```
![image](screenshot/istio-https.png)


## 19. Istio L7 Routing (2/3)
```shell
kubectl apply -f 19-l7route.yaml -n istio-httpserver
```
![image](screenshot/l7route.png)

## 20. Jaeger Tracing (3/3)
```shell
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.13/samples/addons/jaeger.yaml

//istioctl dashboard jaeger -> ubuntu server has no GUI -> "http://localhost:16686" doesn't work ->  edit tracing svc mode : "NodePort".  or

istioctl dashboard --address=0.0.0.0 jaeger
```
![image](screenshot/jaeger-trace.png)

