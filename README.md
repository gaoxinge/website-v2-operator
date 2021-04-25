# website v2 operator

- [kubernetes/code-generator](https://github.com/kubernetes/code-generator)
- [kubernetes/sample-controller](https://github.com/kubernetes/sample-controller)
- [openshift-evangelists/crd-code-generation](https://github.com/openshift-evangelists/crd-code-generation)
- [How to generate client codes for Kubernetes Custom Resource Definitions (CRD)](https://itnext.io/how-to-generate-client-codes-for-kubernetes-custom-resource-definitions-crd-b4b9907769ba)

## install

### base on git

```shell
cd $GOPATH/src
mkdir k8s.io
cd k8s.io
git clone git@github.com:kubernetes/code-generator.git
git checkout tags/v0.20.5 -b v0.20.5
go install ./cmd/{defaulter-gen,client-gen,lister-gen,informer-gen,deepcopy-gen}
```

### base on module

```shell
# need import k8s.io/code-generator
go get k8s.io/code-generator@v0.20.5
go mod vendor
go install ./vendor/k8s.io/code-generator/cmd/{defaulter-gen,client-gen,lister-gen,informer-gen,deepcopy-gen}
```

## generate

### base on install

```shell
hack/generate-groups.sh all github.com/gaoxinge/website-v2-operator/pkg/client github.com/gaoxinge/website-v2-operator/pkg/apis extensions.example.com:v2
```

### base on module

```shell
# need import k8s.io/code-generator
go get k8s.io/code-generator@v0.20.5
go mod vendor
vendor/k8s.io/code-generator/generate-groups.sh all github.com/gaoxinge/website-v2-operator/pkg/client github.com/gaoxinge/website-v2-operator/pkg/apis extensions.example.com:v2
```

## docker

You can build and push image by yourself

```shell
docker build -f docker/Dockerfile -t gaoxinge/website-v2-controller .
docker push gaoxinge/website-v2-controller:latest
```

or directly use

- [docker hub](https://hub.docker.com/r/gaoxinge/website-v2-controller)

## yaml

```shell
kubectl create -f yaml/website-v2-crd.yaml
kubectl create -f yaml/website-v2-controller.yaml
kubectl create serviceaccount website-v2-controller
kubectl create clusterrolebinding website-v2-controller --clusterrole=cluster-admin --serviceaccount=default:website-v2-controller
kubectl create -f yaml/website-v2-example.yaml
kubectl delete website kubia
kubectl delete clusterrolebinding website-v2-controller
kubectl delete serviceaccount website-v2-controller
kubectl delete deployment website-v2-controller
kubectl delete crd websites.extensions.example.com
```

## TODO

- [ ] add update for website
- [ ] add watcher for deployment and service
- [ ] add unit test