### module02 homework

## 1 httpserver golang source code

```shell
go build 1-main.go
./1-main
git add 1-main.go
git commit -m "xxx"
git push
```

### module03 homework

## 2 docker build and push to dockerhub
```shell
docker build -f Dockerfile
docker run
docker push
```

### module 08 part one

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

### module 08 part two

## 5.Create lslsp-httpserver service

```shell
kubectl create -f 5-service.yaml
```

## 6.Create lslsp-tls secret

```shell
kubectl create -f 6.secret.yaml
```

## 7.install nginx-ingress (with jessie_meng's yaml :-> )
```shell
kubectl create -f 7.ng-ingress-deploy-cp-from-cncamp.yaml
```

## 8.Create lslsp-httpserver-ingress ingress

```shell
kubectl create -f 8-ingress.yaml
kubectl get svc -n ingress-nginx
curl -H "Host: lslsp.com" http://<ingress-nginx-controller CLUSTER-IP> -v -k
```
