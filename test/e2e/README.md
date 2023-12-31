# E2E Testing

1. Run `make start pf`
2. Either (a) run your test in your IDE or (b) run `make test`.

Notes:

* Everything runs in the `work` namespace (including MinIO). 
* For speed, please only use `cowsay:v1`. 
* Test can take longer on CI. Adds 5s to timeout values.

## Debugging E2E Tests

### Logs

```
make logs
```

### Accessing Work Server UI

Open http://localhost:2746

### Accessing MinIO

Open http://localhost:9000 using admin/password.

### Expose Database

#### Postgres

Add to `/etc/hosts`:

```
127.0.0.1 postgres
```

#### MySQL

Add to `/etc/hosts`:

```
127.0.0.1 mysql
```

### Running Controller In Your IDE
 
If you want to run the controller in your IDE (e.g. to debug it), firstly scale down the controller:

```
kubectl -n work scale deploy/workflow-controller --replicas 0
```

The run `cmd/workflow-controller/main.go` using these arguments, which enable debug logging, and make sure you use locally build image:

```
--loglevel debug --executor-image kubework/workexec:your-branch --executor-image-pull-policy Never
```

### Running The Work Server In Your IDE

```
kubectl -n work scale deploy/work-server --replicas 0
```

Kill any port forwards on 2746.

The run `server/main.go` using these arguments, which enable debug logging, and make sure you use locally build image:

```
See dist/postgres.yaml
```

### Running The Work Server In Your IDE

```
kubectl -n work scale deploy/work-server --replicas 0
```

Kill any port forwards on 2746.

The run `server/main.go` using these arguments, which enable debug logging, and make sure you use locally build image:

```
See dist/postgres.yaml
```


### To Update The Executor

If you're making changes to the executor, run:

```
make executor-image DEV_IMAGE=true IMAGE_PREFIX=kubework/ IMAGE_TAG=dev 
```

### To Switch Between Postgres and MySQL

Edit `test/e2e/manifest/workflow-controller-config.yaml` and comment/un-comment correct section.

```
kubectl -n work apply test/e2e/manifest/workflow-controller-config.yaml
```

Then either for Postgres: 

```
kubectl -n work scale deploy/mysql --replicas 0
kubectl -n work scale deploy/postgres --replicas 1
```

Or for MySQL

```
kubectl -n work scale deploy/postgres --replicas 0
kubectl -n work scale deploy/mysql --replicas 1
```

To access the Postgres database as follows:

```
make postgres-cli
select * from work_workflows;
```

To access the MySQL database as follows:

```
make mysql-cli
select * from work_workflows;
```

