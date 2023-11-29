## Configure the service account to run Workflows

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
