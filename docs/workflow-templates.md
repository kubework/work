# Workflow Templates

![GA](assets/ga.svg)

> v2.4 and after

Workflow templates are reusable chunks of YAML that you can use within your workflows. This allows you to have a library of templates.

You can create some example templates as follows:

```
work template create https://raw.githubusercontent.com/kubework/work/master/examples/workflow-template/templates.yaml
```

The submit a workflow using one of those templates:

```
work submit https://raw.githubusercontent.com/kubework/work/master/examples/workflow-template/hello-world.yam
```
