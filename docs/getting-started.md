# Getting Started

To see how Work works, you can run examples of simple workflows and workflows that use artifacts.
For the latter, you'll set up an artifact repository for storing the artifacts that are passed in
the workflows. Here are the requirements and steps to run the workflows.

## 0. Requirements
* Kubernetes 1.9 or later
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* Have a [kubeconfig](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/) file (default location is `~/.kube/config`)

## 1. Download the Work CLI

### Mac

Available via `brew`
```sh
brew install kubework/tap/work
```
And via `curl`
```sh
# Download the binary
curl -sLO https://github.com/kubework/work/releases/download/v2.4.3/work-darwin-amd64

# Make binary executable
chmod +x work-darwin-amd64

# Move binary to path
mv ./work-darwin-amd64 /usr/local/bin/work

# Test installation
work version
```

### Linux

Available via `curl`
```sh
# Download the binary
curl -sLO https://github.com/kubework/work/releases/download/v2.4.3/work-linux-amd64

# Make binary executable
chmod +x work-linux-amd64

# Move binary to path
mv ./work-linux-amd64 /usr/local/bin/work

# Test installation
work version
```

### Binaries

You can download the latest and previous Work binaries from our [releases page](https://github.com/kubework/work/releases/).

## 2. Install the Controller and UI
```sh
kubectl create namespace work
kubectl apply -n work -f https://raw.githubusercontent.com/kubework/work/v2.4.3/manifests/install.yaml
```
NOTE: On GKE, you may need to grant your account the ability to create new `clusterrole`s
```sh
kubectl create clusterrolebinding YOURNAME-cluster-admin-binding --clusterrole=cluster-admin --user=YOUREMAIL@gmail.com
```

## 3. Configure the service account to run Workflows

### Roles, RoleBindings, and ServiceAccounts

In order for Work to support features such as artifacts, outputs, access to secrets, etc. it needs to communicate with Kubernetes resources
using the Kubernetes API. To communicate with the Kubernetes API, Work uses a `ServiceAccount` to authenticate itself to the Kubernetes API.
You can specify which `Role` (i.e. which permissions) the `ServiceAccount` that Work uses by binding a `Role` to a `ServiceAccount` using a `RoleBinding`

Then, when submitting Workflows you can specify which `ServiceAccount` Work uses using:

```sh
work submit --serviceaccount <name>
```

When no `ServiceAccount` is provided, Work will use the `default` `ServiceAccount` from the namespace from which it is run, which will almost always have insufficient privileges by default.

For more information about granting Work the necessary permissions for your use case see [Workflow RBAC](workflow-rbac.md).

### Granting admin privileges

For the purposes of this demo, we will grant the `default` `ServiceAccount` admin privileges (i.e., we will bind the `admin` `Role` to the `default` `ServiceAccount` of the current namespace):

```sh
kubectl create rolebinding default-admin --clusterrole=admin --serviceaccount=default:default
```

**Note that this will grant admin privilages to the `default` `ServiceAccount` in the namespace that the command is run from, so you will only be able to
run Workflows in the namespace where the `RoleBinding` was made.**

## 4. Run Sample Workflows
```sh
work submit --watch https://raw.githubusercontent.com/kubework/work/master/examples/hello-world.yaml
work submit --watch https://raw.githubusercontent.com/kubework/work/master/examples/coinflip.yaml
work submit --watch https://raw.githubusercontent.com/kubework/work/master/examples/loops-maps.yaml
work list
work get xxx-workflow-name-xxx
work logs xxx-pod-name-xxx #from get command above
```

Additional examples and more information about the CLI are available on the [Work Workflows by Example](../examples/README.md) page.

You can also create Workflows directly with `kubectl`. However, the Work CLI offers extra features
that `kubectl` does not, such as YAML validation, workflow visualization, parameter passing, retries
and resubmits, suspend and resume, and more.
```sh
kubectl create -f https://raw.githubusercontent.com/kubework/work/master/examples/hello-world.yaml
kubectl get wf
kubectl get wf hello-world-xxx
kubectl get po --selector=workflows.kubework.io/workflow=hello-world-xxx --show-all
kubectl logs hello-world-yyy -c main
```


## 5. Install an Artifact Repository

Work supports S3 (AWS, GCS, Minio) and Artifactory as artifact repositories. Instructions on how to configure artifact repositories are available on the [Configuring your Artifact Repository](configure-artifact-repository.md) page.

This tutorial uses Minio for the sake of portability.

Install Minio:
```sh
helm install work-artifacts stable/minio \
  --set service.type=LoadBalancer \
  --set defaultBucket.enabled=true \
  --set defaultBucket.name=my-bucket \
  --set persistence.enabled=false \
  --set fullnameOverride=work-artifacts
```

Login to the Minio UI using a web browser (port 9000) after exposing obtaining the external IP using `kubectl`.
```sh
kubectl -n work get service work-artifacts -o wide
```
On Minikube:
```sh
minikube -n work service --url work-artifacts
```

NOTE: When minio is installed via Helm, it uses the following hard-wired default credentials,
which you will use to login to the UI:
* AccessKey: `AKIAIOSFODNN7EXAMPLE`
* SecretKey: `wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY`

Create a bucket named `my-bucket` from the Minio UI.

## 6. Reconfigure the workflow controller to use the Minio artifact repository

Edit the `workflow-controller` `ConfigMap` to reference the service name (`work-artifacts`) and
secret (`work-artifacts`) created by the Helm install:

Edit the `workflow-controller` `ConfigMap`:
```sh
kubectl edit cm -n work workflow-controller-configmap
```
Add the following:
```yaml
data:
  config: |
    artifactRepository:
      s3:
        bucket: my-bucket
        endpoint: work-artifacts:9000
        insecure: true
        # accessKeySecret and secretKeySecret are secret selectors.
        # It references the k8s secret named 'work-artifacts'
        # which was created during the minio helm install. The keys,
        # 'accesskey' and 'secretkey', inside that secret are where the
        # actual minio credentials are stored.
        accessKeySecret:
          name: work-artifacts
          key: accesskey
        secretKeySecret:
          name: work-artifacts
          key: secretkey
```

NOTE: the Minio secret is retrieved from the namespace you use to run Workflows. If Minio is
installed in a different namespace then you will need to create a copy of its secret in the
namespace you use for Workflows.

## 7. Run a workflow which uses artifacts
```sh
work submit https://raw.githubusercontent.com/kubework/work/master/examples/artifact-passing.yaml
```

## 8. Access the Work UI

> v2.5 and after

```
kubectl -n work port-forward deployment/work-server 2746:2746
```

Then visit: http://127.0.0.1:2746

> v2.4 and before

By default, the Work UI service is not exposed with an external IP. To access the UI, use one of the
following:

### Method 1: kubectl port-forward
```
kubectl -n work port-forward deployment/work-ui 8001:8001
```
Then visit: http://127.0.0.1:8001

### Method 2: kubectl proxy
```
kubectl proxy
```
Then visit: http://127.0.0.1:8001/api/v1/namespaces/work/services/work-ui/proxy/

NOTE: artifact download and webconsole is not supported using this method

### Method 3: Expose a LoadBalancer
Update the work-ui service to be of type `LoadBalancer`.
```
kubectl patch svc work-ui -n work -p '{"spec": {"type": "LoadBalancer"}}'
```
Then wait for the external IP to be made available:
```
kubectl get svc work-ui -n work
NAME      TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)        AGE
work-ui   LoadBalancer   10.19.255.205   35.197.49.167   80:30999/TCP   1m
```

NOTE: On Minikube, you won't get an external IP after updating the service -- it will always show
`pending`. Run the following command to determine the Work UI URL:
```
minikube service -n work --url work-ui
```
