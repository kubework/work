# Work Server

![alpha](assets/alpha.svg)

> v2.5 and after

The Work Server is a server that exposes an API and UI for workflows. You'll need to use this if you want to [offload large workflows](offloading-large-workflows.md) or the [workflow archive](workflow-archive.md).

You can run this in either "hosted" or "local" mode.

It replaces the Work UI.

## Hosted Mode

Use this mode if:

* You want a drop-in replacement for the Work UI.
* If you need to prevent users from directly accessing the database.

Hosted mode is provided as part of the standard [manifests](../manifests), [specifically in work-server-deployment.yaml](../manifests/base/work-server/work-server-deployment.yaml) .

## Local Mode

Use this mode if:

* You want something that does not require complex set-up. 
* You do not need to run a database.

To run locally:

```
work server
```

This will start a server on port 2746 which you can view at [http://localhost:2746](http://localhost:2746]).

## Options

### Auth Mode

You can choose which kube config the server uses: 

* "server" - in hosted mode, use the kube config of service account, in local mode, use your local kube config.  
* "client" - requires client to provide their Kubernetes bearer token and use that.
* "hybrid" - use the client token if provided, fallback to the server token if note.

By default, the server will start with auth mode of "client" (highest security).

### Managed Namespace

See [managed namespace](managed-namespace.md).
