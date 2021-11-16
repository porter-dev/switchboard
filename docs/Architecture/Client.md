# Client
The Switchboard client is in charge of managing and displaying [[Template Reference|Templates]]. Each client reads data from the Porter API and creates views for the developer. Both the CLI and the dashboard contain a client, although the CLI client is extremely limited in terms of the types of views it can support. The output of the client is a set of resources, each controlled by a separate driver, which are sent to the API. 

The internal components of the client, listed in order of data flow, are described below. 

## 1 - Context

The context sets the initial state of the engine. This section defines how resources are scoped, which fields need to be set, and which fields can be `dynamic`, `static`, or `preset`. [[Template Reference#Context|See the template reference]]. 

## 2 - Data Queries

Based on the context, we construct a set of data queries. [[Template Reference#Query|See the template reference]]. 

The Renderer computes protocols and generates the relevant API endpoints based on the template. For example, in the [[Template Reference#Full Example|example template]], we have two data sources to read from: `get_helm_charts` and `get_namespaces`. In order to read these data sources, we construct the following:

-   `get_helm_charts`
    -   `GET <https://dashboard.getporter.dev/api/projects/1/helm/charts?tag_filter=applications`
    -   `wss://dashboard.getporter.dev/api/projects/1/helm/charts?tag_filter=applications`
-   `get_namespaces`
    -   `GET https://dashboard.getporter.dev/api/projects/1/clusters/1/kubernetes/namespaces`
    -   `wss://dashboard.getporter.dev/api/projects/1/clusters/1/kubernetes/namespaces`

## 3 - Transforms

The data queries are transformed via a transform engine which modifies queries before sending them to the Renderer. `jq` will be the first transform engine supported. The transform engine can transform the raw response query to a more usable form for the rendering engine before the data arrives at the rendering engine. It applies to both streamable and initial data.

Because the transform engine can inject arbitrary data at this step, the frontend may not modify the transformation function which is sent to the backend. We add verification so that that API can check that the transform function has not been modified by the frontend. Thus, the frontend must pass the following parameters when passing a transform to the backend: 
- `X-Switch-Signature-256` -- a SHA 256 hash of the transform
- `transform` query param -- a base64-encoded transform 
- `template_session_id` -- the initial session id which is used as the salt

## 4 - Component Rendering

Component rendering refers to how frontend components are rendered. This renders three types of components: `mutator`, `display`, and `action`. Mutators modify the internal state, display components are simply for rendering internal state, and actions modify external state. This is basically a formalized version of what `form.yaml` currently is.

## 5 - Actions

Actions modify external state (Kubernetes objects, Helm releases, Terraform modules), and all supported actions are defined in the `actions` section. For example:

```yaml
actions:
- resource: terraform
  path:
  - project
  - module
  data:
    db_username: {{ db_username }}
    db_password: {{ db_password }}
```

Here you'll see the concept of **internal state referencing** -- these variables are present in the Renderer's internal state. All component's values can be referenced, and they can be undefined. 