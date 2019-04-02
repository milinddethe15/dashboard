<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Developer guide for bootstrap](#developer-guide-for-bootstrap)
  - [Prerequisites](#prerequisites)
  - [Setting up the build environment](#setting-up-the-build-environment)
  - [Building kfctl](#building-kfctl)
        - [`make build-kfctl`](#make-build-kfctl)
        - [`make install`](#make-install)
        - [`make build-kftl-container`](#make-build-kftl-container)
        - [`make push-kftl-container`](#make-push-kftl-container)
        - [`make push-kftl-container-latest`](#make-push-kftl-container-latest)
        - [`make test-init`](#make-test-init)
  - [Building bootstrap](#building-bootstrap)
        - [`make build-bootstrap`](#make-build-bootstrap)
        - [`make build`](#make-build)
        - [`make push`](#make-push)
        - [`make push-latest`](#make-push-latest)
        - [`make static, make plugins`](#make-static-make-plugins)
  - [How to run bootstrapper with Click-to-deploy app locally](#how-to-run-bootstrapper-with-click-to-deploy-app-locally)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Developer guide for bootstrap

## Prerequisites

golang to 1.12

```sh
$ ☞  go version
go version go1.12 darwin/amd64
```

On mac osx you can run

```sh
brew upgrade golang
```

On linux you can install the [download](https://golang.org/dl/) and follow the [installation directions](https://golang.org/doc/install).

golang-1.12 uses the golang module framework. See [golang Modules](https://github.com/golang/go/wiki/Modules).
You should add the environment variable `GO111MODULE=on` to your shell init file

Normally go will build a dependency list in the go.mod file, installing 
a dependency explicitly is done by running `go get <dependency>`. 
See the [use case](https://github.com/golang/go/wiki/Modules#why-am-i-getting-an-error-cannot-find-module-providing-package-foo) in the golang Modules Wiki.
golang-1.12 no longer creates a vendor directory.

## Setting up the build environment

Download the kubeflow repo by running

```sh
git clone git@github.com:kubeflow/kubeflow.git
```

or

```sh
git clone https://github.com:kubeflow/kubeflow .
```

Create a symbolic link inside your GOPATH to the location you checked out the code

```sh
mkdir -p ${GOPATH}/src/github.com/kubeflow
ln -sf ${GIT_KUBEFLOW} ${GOPATH}/src/github.com/kubeflow/kubeflow
```

* GIT_KUBEFLOW should be the location where you checked out https://github.com/kubeflow/kubeflow. GOPATH is typically $HOME/go.


## Building kfctl

##### `make build-kfctl`
```sh
cd $GIT_KUBEFLOW/bootstrap
make build-kfctl
```

* This will create `bin/kfctl` with full debug information

* If you get an error about missing files in `/tmp/v2`, you are hitting [#2790](https://github.com/kubeflow/kubeflow/issues/2790) and need to delete `/tmp/v2` and rerun the build.

##### `make install`

```sh
make install #depends on build-kfctl
```

* Installs kfctl in /usr/local/bin


##### `make build-kftl-container`

```sh
make build-kfctl-container
```

* creates a docker container

##### `make push-kftl-container`

```sh
make push-kfctl-container
```

* pushes the docker container to $IMG_KFCTL

##### `make push-kftl-container-latest`

```sh
make push-kfctl-container-latest
```

* pushes the docker container to $IMG_KFCTL and tags it with kfctl:latest

##### `make run-kfctl-container`

```sh
make run-kfctl-container
```

* Runs the local docker container


##### `make test-init`

```sh
make test-init
```

* will run `kfctl init` for gcp, minikube and no platform

## Building bootstrap 

##### `make build-bootstrap`

```sh
cd $GIT_KUBEFLOW/bootstrap
make build-bootstrap
```

* Creates bin/bootstrapper with full debug information

##### `make build`

* Depends on `make build-local`. Creates a docker image gcr.io/$(GCLOUD_PROJECT)/bootstrapper:$(TAG)

##### `make push`

* Depends on `make build`. Pushes the docker image gcr.io/$(GCLOUD_PROJECT)/bootstrapper:$(TAG)

##### `make push-latest`

* Depends on `make push`. Tags the docker image gcr.io/$(GCLOUD_PROJECT)/bootstrapper:$(TAG) with latest.
Note: To use a different gcloud project than kubeflow-images-public.
```sh
export GCLOUD_PROJECT=mygcloudproject
make push
```

##### `make static, make plugins`
These targets are for kfctl and allows the goland debugger work by disabling plugins.
This is a problem in the go compiler which should be fixed in 1.12.
See the [kfctl/README.md](./cmd/kfctl) for additional information.

## How to run bootstrapper with Click-to-deploy app locally

Start the backend:

```
IMAGE=gcr.io/kubeflow-images-public/bootstrapper:latest  # change this

docker run -d -it --name bootstrapper \
    --mount type=bind,source=${HOME}/kf_app,target=/home/kubeflow -p 8080:8080 \
    ${IMAGE} /opt/kubeflow/bootstrapper \
    --install-istio --namespace=kubeflow  # change args if you want
```

Start the frontend:

```
cd ../components/gcp-click-to-deploy
npm start
```

## Releasing kfctl

See [release guide](https://github.com/kubeflow/kubeflow/blob/master/docs_dev/releasing.md)