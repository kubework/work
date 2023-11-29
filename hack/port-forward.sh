#!/usr/bin/env bash
set -eu -o pipefail

killall kubectl || true

info() {
    echo '[INFO] ' "$@"
}

info "MinIO on http://localhost:9000"
kubectl -n work port-forward pod/minio 9000:9000 &

work_server=$(kubectl -n work get pod -l app=work-server -o name)
if [[ "$work_server" != "" ]]; then
  info "Work Server on http://localhost:2746"
  kubectl -n work port-forward svc/work-server 2746:2746 &
fi

postgres=$(kubectl -n work get pod -l app=postgres -o name)
if [[ "$postgres" != "" ]]; then
  info "Postgres on http://localhost:5432"
  kubectl -n work port-forward "$postgres" 5432:5432 &
fi

mysql=$(kubectl -n work get pod -l app=mysql -o name)
if [[ "$mysql" != "" ]]; then
  info "MySQL on http://localhost:3306"
  kubectl -n work port-forward "$mysql" 3306:3306 &
fi

wait