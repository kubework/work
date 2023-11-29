[![slack](https://img.shields.io/badge/slack-kubework-brightgreen.svg?logo=slack)](https://kubework.github.io/community/join-slack)

# Kubework - Get stuff done with Kubernetes

![Work Image](docs/assets/work.png)

## Quickstart
```bash
kubectl create namespace work
kubectl apply -n work -f https://raw.githubusercontent.com/kubework/work/stable/manifests/install.yaml
```

## What is Kubework?

Kubework is a collection of tools for getting work done with Kubernetes.
* [Work Workflows](https://github.com/kubework/work) - Container-native Workflow Engine
* [Work CD](https://github.com/kubework/work-cd) - Declarative GitOps Continuous Delivery
* [Work Events](https://github.com/kubework/work-events) - Event-based Dependency Manager
* [Work Rollouts](https://github.com/kubework/work-rollouts) - Deployment CR with support for Canary and Blue Green deployment strategies

Also kubework-labs is a separate GitHub org that we setup for community contributions related to the Kubework ecosystem. Repos in kubework-labs are administered by the owners of each project. Please reach out to us on the Work slack channel if you have a project that you would like to add to the org to make it easier to others in the Work community to find, use, and contribute back.
* https://github.com/kubework-labs

## What is Work Workflows?
Work Workflows is an open source container-native workflow engine for orchestrating parallel jobs on Kubernetes. Work Workflows is implemented as a Kubernetes CRD (Custom Resource Definition).

* Define workflows where each step in the workflow is a container.
* Model multi-step workflows as a sequence of tasks or capture the dependencies between tasks using a graph (DAG).
* Easily run compute intensive jobs for machine learning or data processing in a fraction of the time using Work Workflows on Kubernetes.
* Run CI/CD pipelines natively on Kubernetes without configuring complex software development products.

## Why Work Workflows?
* Designed from the ground up for containers without the overhead and limitations of legacy VM and server-based environments.
* Cloud agnostic and can run on any Kubernetes cluster.
* Easily orchestrate highly parallel jobs on Kubernetes.
* Work Workflows puts a cloud-scale supercomputer at your fingertips!

## Documentation
* [Get started here](docs/getting-started.md)
* [How to write Work Workflow specs](examples/README.md)
* [How to configure your artifact repository](docs/configure-artifact-repository.md)

## Features
* DAG or Steps based declaration of workflows
* Artifact support (S3, Artifactory, HTTP, Git, raw)
* Step level input & outputs (artifacts/parameters)
* Loops
* Parameterization
* Conditionals
* Timeouts (step & workflow level)
* Retry (step & workflow level)
* Resubmit (memoized)
* Suspend & Resume
* Cancellation
* K8s resource orchestration
* Exit Hooks (notifications, cleanup)
* Garbage collection of completed workflow
* Scheduling (affinity/tolerations/node selectors)
* Volumes (ephemeral/existing)
* Parallelism limits
* Daemoned steps
* DinD (docker-in-docker)
* Script steps

## Who uses Work?
As the Work Community grows, we'd like to keep track of our users. Please send a PR with your organization name.

Currently **officially** using Work:

1. [Adevinta](https://www.adevinta.com/)
1. [Admiralty](https://admiralty.io/)
1. [Adobe](https://www.adobe.com/)
1. [Alibaba Cloud](https://www.alibabacloud.com/about)
1. [Ant Financial](https://www.antfin.com/)
1. [BasisAI](https://basis-ai.com/)
1. [BioBox Analytics](https://biobox.io)
1. [BlackRock](https://www.blackrock.com/)
1. [Canva](https://www.canva.com/)
1. [Capital One](https://www.capitalone.com/tech/)
1. [CarTrack](https://www.cartrack.com/)
1. [CCRi](https://www.ccri.com/)
1. [Codec](https://www.codec.ai/)
1. [Commodus Tech](https://www.commodus.tech)
1. [CoreFiling](https://www.corefiling.com/)
1. [Cratejoy](https://www.cratejoy.com/)
1. [CyberAgent](https://www.cyberagent.co.jp/en/)
1. [Cyrus Biotechnology](https://cyrusbio.com/)
1. [Datadog](https://www.datadoghq.com/)
1. [DataStax](https://www.datastax.com/)
1. [EBSCO Information Services](https://www.ebsco.com/)
1. [Equinor](https://www.equinor.com/)
1. [Fairwinds](https://fairwinds.com/)
1. [Gardener](https://gardener.cloud/)
1. [GitHub](https://github.com/)
1. [Gladly](https://gladly.com/)
1. [Google](https://www.google.com/intl/en/about/our-company/)
1. [Greenhouse](https://greenhouse.io)
1. [HOVER](https://hover.to)
1. [IBM](https://ibm.com)
1. [InsideBoard](https://www.insideboard.com)
1. [Interline Technologies](https://www.interline.io/blog/scaling-openstreetmap-data-workflows/)
1. [Intuit](https://www.intuit.com/)
1. [Karius](https://www.kariusdx.com/)
1. [KintoHub](https://www.kintohub.com/)
1. [Localytics](https://www.localytics.com/)
1. [Maersk](https://www.maersk.com/solutions/digital-solutions)
1. [Max Kelsen](https://maxkelsen.com/)
1. [Mirantis](https://mirantis.com/)
1. [NVIDIA](https://www.nvidia.com/)
1. [OVH](https://www.ovh.com/)
1. [Peak AI](https://www.peak.ai/)
1. [Preferred Networks](https://www.preferred-networks.jp/en/)
1. [Quantibio](http://quantibio.com/us/en/)
1. [Ramboll Shair](https://ramboll-shair.com/)
1. [Red Hat](https://www.redhat.com/en)
1. [SAP Fieldglass](https://www.fieldglass.com/)
1. [SAP Hybris](https://cx.sap.com/)
1. [Sidecar Technologies](https://hello.getsidecar.com/)
1. [Styra](https://www.styra.com/)
1. [Threekit](https://www.threekit.com/)
1. [Tiger Analytics](https://www.tigeranalytics.com/)
1. [Wavefront](https://www.wavefront.com/)

## Community Blogs and Presentations
* [Work Ansible role: Provisioning Work Workflows on OpenShift](https://medium.com/@marekermk/provisioning-work-on-openshift-with-ansible-and-kustomize-340a1fda8b50)
* [Work Workflows vs Apache Airflow](http://bit.ly/30YNIvT)
* [CI/CD with Work on Kubernetes](https://medium.com/@bouwe.ceunen/ci-cd-with-work-on-kubernetes-28c1a99616a9)
* [Running Work Workflows Across Multiple Kubernetes Clusters](https://admiralty.io/blog/running-work-workflows-across-multiple-kubernetes-clusters/)
* [Open Source Model Management Roundup: Polyaxon, Work, and Seldon](https://www.anaconda.com/blog/developer-blog/open-source-model-management-roundup-polyaxon-work-and-seldon/)
* [Producing 200 OpenStreetMap extracts in 35 minutes using a scalable data workflow](https://www.interline.io/blog/scaling-openstreetmap-data-workflows/)
* [Work integration review](http://dev.matt.hillsdon.net/2018/03/24/work-integration-review.html)
* TGI Kubernetes with Joe Beda: [Work workflow system](https://www.youtube.com/watch?v=M_rxPPLG8pU&start=859)
* [Community meeting minutes and recordings](https://docs.google.com/document/d/16aWGQ1Te5IRptFuAIFtg3rONRQqHC1Z3X9rdDHYhYfE)

## Project Resources
* Work GitHub:  https://github.com/kubework
* Work website: https://kubework.github.io/
* Work Slack:   [click here to join](https://kubework.github.io/community/join-slack)
