####################################################################################################
# Builder image
# Initial stage which pulls prepares build dependencies and CLI tooling we need for our final image
# Also used as the image in CI jobs so needs all dependencies
####################################################################################################
FROM golang:1.13.4 as builder

RUN apt-get update && apt-get install -y \
    git \
    make \
    wget \
    gcc \
    zip && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /tmp

# Install docker
ENV DOCKER_CHANNEL stable
ENV DOCKER_VERSION 18.09.1
RUN wget -O docker.tgz "https://download.docker.com/linux/static/${DOCKER_CHANNEL}/x86_64/docker-${DOCKER_VERSION}.tgz" && \
    tar --extract --file docker.tgz --strip-components 1 --directory /usr/local/bin/ && \
    rm docker.tgz

# Install dep
ENV DEP_VERSION=0.5.0
RUN wget https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 -O /usr/local/bin/dep && \
    chmod +x /usr/local/bin/dep

####################################################################################################
# workexec-base
# Used as the base for both the release and development version of workexec
####################################################################################################
FROM debian:10.3-slim as workexec-base

# NOTE: keep the version synced with https://storage.googleapis.com/kubernetes-release/release/stable.txt
ENV KUBECTL_VERSION=1.15.1
ENV JQ_VERSION=1.6
RUN apt-get update && \
    apt-get install -y curl procps git tar mime-support && \
    rm -rf /var/lib/apt/lists/* && \
    curl -L -o /usr/local/bin/kubectl -LO https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl && \
    chmod +x /usr/local/bin/kubectl && \
    curl -L -o /usr/local/bin/jq -LO https://github.com/stedolan/jq/releases/download/jq-${JQ_VERSION}/jq-linux64 && \
    chmod +x /usr/local/bin/jq
COPY hack/ssh_known_hosts /etc/ssh/ssh_known_hosts
COPY --from=builder /usr/local/bin/docker /usr/local/bin/

####################################################################################################

FROM node:11.15.0 as work-ui

ADD ["ui", "."]

RUN yarn install --frozen-lockfile --ignore-optional --non-interactive
RUN yarn build

####################################################################################################
# Work Build stage which performs the actual build of Work binaries
####################################################################################################
FROM builder as work-build

# Perform the build
WORKDIR /go/src/github.com/kubework/work
COPY . .
# check we can use Git
RUN git rev-parse HEAD
COPY --from=work-ui node_modules ui/node_modules
RUN mkdir -p ui/dist
COPY --from=work-ui dist/app ui/dist/app
# stop make from trying to re-build this without yarn installed
RUN touch ui/dist/app
RUN make dist/work-linux-amd64 dist/workflow-controller-linux-amd64 dist/workexec-linux-amd64

####################################################################################################
# workexec
####################################################################################################
FROM workexec-base as workexec
COPY --from=work-build /go/src/github.com/kubework/work/dist/workexec-linux-amd64 /usr/local/bin/workexec
ENTRYPOINT [ "workexec" ]

####################################################################################################
# workflow-controller
####################################################################################################
FROM scratch as workflow-controller
# Add timezone data
COPY --from=work-build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=work-build /go/src/github.com/kubework/work/dist/workflow-controller-linux-amd64 /bin/workflow-controller
ENTRYPOINT [ "workflow-controller" ]

####################################################################################################
# workcli
####################################################################################################
FROM scratch as workcli
COPY --from=workexec-base /etc/ssh/ssh_known_hosts /etc/ssh/ssh_known_hosts
COPY --from=workexec-base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=work-build /go/src/github.com/kubework/work/dist/work-linux-amd64 /bin/work
ENTRYPOINT [ "work" ]
