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
kubectl create configmap lslsp-config --from-env-file=lslsp-config.env
```or
kubectl create -f 3-lslsp-cofig.yaml
```

## 4.Create httpserver deployment (with 5 pods)

```shell
kubectl create -f 4-deploy.yaml
```

### _module 08 (2/2)_

## 5.Create lslsp-httpserver service

```shell
kubectl create -f 5-service.yaml
```

## 6.Create lslsp-tls secret

```shell
kubectl create -f 6-secret.yaml
```

## 7.install nginx-ingress (with jessie_meng's yaml :-> )
```shell
kubectl create -f 7-ng-ingress-deploy-cp-from-cncamp.yaml
```

## 8.Create lslsp-httpserver-ingress ingress

```shell
kubectl create -f 8-ingress.yaml
kubectl get svc -n ingress-nginx
curl -H "Host: lslsp.com" https://<ingress-nginx-controller CLUSTER-IP> -v -k
```

### _module 10_

## 9. re-Go & Git with Metrics.go

```shell
go mod init
go mod tidy
go build 1-main.go
sudo ./1-main

git add 1-main.go
git add metrics/metrics.go
git commit -m "xxx"
git push
```
## 10. Docker Build & Push v0.10
```shell
docker build . -t lslsp2/httpserver:0.10 -f 2-Dockerfile
docker login
docker push lslsp2/httpserver:0.10
docker run -p 80:80 lslsp2/httpserver:0.10
```


