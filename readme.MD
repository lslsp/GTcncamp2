## 1.Create configmap

```shell
$ kubectl create configmap lslsp-config --from-env-file=lslsp-config.env
###or use 
### $ kubectl create -f lslsp-cofig.yaml
```
## 2.Create httpserver deployment

```shell
$ kubectl create -f deploy.yaml
```
