# Release Instructions

Allow 1h to do a release.

## Preparation

Cherry-pick your changes from master onto the release branch.

Mandatory: the release branch must be green [in CircleCI](https://app.circleci.com/github/kubework/work/pipelines).

It is a very good idea to clean up before you start:

    make clean
    kubectl delete ns work

## Release

To generate new manifests and perform basic checks:

    make prepare-release VERSION=v2.5.0-rc6

Next, build everything:

    make build

Publish the images and local Git changes:

    make publish-release

Create [the release](https://github.com/kubework/work/releases) in Github. You can get some text for this using [Github Toolkit](https://github.com/alexec/github-toolkit):

    ght relnote v2.5.0-rc5..v2.5.0-rc6

    
## Validation

K3D tip: you'll need to import the images:

    k3d import-images kubework/workcli:v2.5.0-rc6 kubework/workexec:v2.5.0-rc6 kubework/workflow-controller:v2.5.0-rc6

Install Work locally:

    kubectl create ns work
    kubectl apply -n work -f https://raw.githubusercontent.com/kubework/work/v2.5.0-rc6/manifests/quick-start-postgres.yaml
    make pf-bg 

Maybe run e2e tests?

    make test-e2e
