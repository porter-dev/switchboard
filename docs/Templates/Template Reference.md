# Template Reference
- `version`:
	- Type: `String`
	- Description: the template version being used. Should correspond with a Porter API version, like `v1`.
- `context`:
	- Type: [[Template Reference#Context|Context]]
	- Description: defines the context set for the template which can be referenced in queries.
- `data`:
	- Type: \[\][[Template Reference#Query|Query]]
	- Description: defines the set of data queries for the Renderer. 

## Context
-  `preset`: 
	- Type: \[\][[Template Reference#Key|Key]]
	- Description: defines the keys which can be set before the template is rendered.
- `static`
	- Type: \[\][[Template Reference#Key|Key]]
	- Description: defines the keys which are static for the template. This is like setting a global constant.
- `dynamic`
	- Type: \[\][[Template Reference#Key|Key]]
	- Description: defines the keys which are dynamic for the template. These keys can be set by the Renderer.  

Example:

```yaml
context: 
  preset:
  - name: project
  - name: cluster
  static:
  - name: namespace
    default: default
  dynamic:
  - name: tag_filter
    default: applications
```

### Key
- `name`:
	- Type: `String`
	- Description: the name of the key
- `default`:
	- Type: `String`
	- Description: the default value for the key. Setting the default value makes this key non-required. 
	
## Query
- `name`:
	- Type: `String`
	- Description: the name of the query. This can be whatever you'd like. 
- `driver`:
	- Type: `String`
	- Description: references a driver, like `helm`, `k8s`, or `terraform`. See the list of supported drivers. 
- `resource`:
	- Type: `String`
	- Description: references a resource for that driver, like `charts`
- `path`:
	- Type: `[]String`
	- Description: references a set of context fields to be passed as path parameters. Context fields will be added to the path in the following order:
		- `preset` fields will be added in the order they are listed, **before** the resource is uniquely identified (i.e. before the `/<driver>/<resource>` part of the path)
		- `static` and `dynamic` fields will be added in the order they are listed, **after** the resource is uniquely identified
- `query`:
	- Type: `[]string`
	- Description: references a set of context fields to be passed as query parameters.
- `transform`:
	- Type: [[Template Reference#Transform|Transform]]
	- Description: the transformation to apply to this query.

Example:

```yaml
data: 
- name: get_helm_charts
  driver: helm
  resource: charts
  path: 
  - project 
  query: 
  - tag_filter
- name: get_namespaces
  driver: kubernetes
  resource: namespaces
  path: 
  - project 
  - cluster
  transform:
    engine: jq 
    query: get_namespaces 
    statement: | 
    .items[] | { name: .metadata.name }
```

### Transform
- `engine`:
	- Type: `string`
	- Description: references a transform engine, like `jq`. See the list of transform engines. 
- `query`: 
	- Type: `string`
	- Description: references the `name` field of a query as defined in the [[Template Reference#Query|Query]] section. 
- `statement`:
	- Type: `string`
	- Description: the transform instructions.

# Full Example
```yaml
version: v1
context: 
  preset:
  - name: project
  - name: cluster
  static:
  - name: namespace
    default: default
  dynamic:
  - name: tag_filter
    default: applications
data: 
- name: get_helm_charts
  driver: helm
  resource: charts
  path: 
  - project 
  query: 
  - tag_filter
- name: get_namespaces
  driver: kubernetes
  resource: namespaces
  path: 
  - project 
  - cluster
  transform:
    engine: jq 
    query: get_namespaces 
    statement: | 
    .items[] | { name: .metadata.name }
```