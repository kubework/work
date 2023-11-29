# REST API

## Work Server API

![alpha](assets/alpha.svg)

> v2.5 and after

Since version v2.5 Work Workflows ships with a server that provide more features and security than before.

The server can be configured with or without client auth (`server --auth-mode client`). When it is disabled, then clients must pass their Kubeconfig base 64 encoded in the HTTP `Authorization` header:

```
token=$(work auth token)
curl -H "Authorization: $token" http://localhost:2746/api/v1/workflows/work
```

To view the API:
 
1. Open [https://editor.swagger.io/](https://editor.swagger.io/)
2. Copy and paste ../api/work-server/swagger.json

## Classic API

![ga](assets/ga.svg)

> v2.4 and before

Work is implemented as a kubernetes controller and Workflow [Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/).
Work itself does not run an API server, and with all CRDs, it extends the Kubernetes API server by
introducing a new API Group/Version (kubework.io/v1alpha1) and Kind (Workflow). When CRDs are
registered in a cluster, access to those resources are made available by exposing new endpoints in
the kubernetes API server. For example, to list workflows in the default namespace, a client would
make an HTTP GET request to: `https://<k8s-api-server>/apis/kubework.io/v1alpha1/namespaces/default/workflows`

> NOTE: the optional work-ui does run a thin API layer to power the UI, but is not intended for
  programatic interaction.

A common scenario is to programatically submit and retrieve workflows. To do this, you would use the
existing Kubernetes REST client in the language of preference, which often libraries for performing
CRUD operation on custom resource objects.

## Examples

### Golang 

A kubernetes Workflow clientset library is auto-generated under [kubework/work/pkg/client](https://github.com/kubework/work/tree/master/pkg/client) and can be imported by golang
applications. See the [golang code example](../examples/example-golang/main.go) on how to make use of this client.

### Python
The python kubernetes client has libraries for interacting with custom objects. See: https://github.com/kubernetes-client/python/blob/master/kubernetes/docs/CustomObjectsApi.md


### Java
The Java kubernetes client has libraries for interacting with custom objects. See:
https://github.com/kubernetes-client/java/blob/master/kubernetes/docs/CustomObjectsApi.md

### Ruby
The Ruby kubernetes client has libraries for interacting with custom objects. See:
https://github.com/kubernetes-client/ruby/tree/master/kubernetes
See this [external Ruby example](https://github.com/fischerjulian/work_workflows_ruby_example) on how to make use of this client.

## OpenAPI

An OpenAPI Spec is generated under [kubework/work/api/openapi-spec](https://github.com/kubework/work/blob/master/api/openapi-spec/swagger.json). This spec may be
used to auto-generate concrete datastructures in other languages.
