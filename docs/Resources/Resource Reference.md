# Resource Reference
- `version`:
	- Type: `String`
	- Description: the resource version being used. Should correspond with a Porter API version, like `v1`.
- `resources`:
	- Type: \[\][[Resource Reference#Resource|Resource]]
	- Description: describes a set of grouped resources.

## Resource
- `name`:
	- Type: `String`
	- Description: the name of the resource. This can be whatever you like, but will typically match the name of the underlying resource created by the driver. 
- `driver`:
	- Type: `String`
	- Description: references a driver, like `helm`, `k8s`, or `terraform`. See the list of supported drivers. 
- `source`:
	- Type: [[Resource Reference#Source|Source]]
	- Description: the source configuration for the driver. This is driver-specific, but is usually the path to a registry, Github repo, etc.
- `target`:
	- Type: [[Resource Reference#Target|Target]]
	- Description: the target for the resource. This can any supported Porter target, such as a project or a cluster. 
- `config`:
	- Type: `Object`
	- Description: arbitrary configuration used by the driver.

### Source
- `auth`
	- Type: [[Resource Reference#SourceAuth|SourceAuth]]

#### SourceAuth
- `kind`:
	- Type: `String`
	- Description: the type of source auth, like `basic`, `bearer`, `github`, etc

- `id`:
	- Type: `Integer|String`
	- Description: the ID of the source authentication from the Porter API, or the string value of the source authentication. 

### Target
- `project`:
	- Type: `Integer`
	- Description: the ID of the targeted Porter project.

# Examples

## Helm Chart

```yaml
version: v1
resources:
- name: web-example
  driver: helm
  source:
    kind: helm_registry
	config:
	  url: https://porter-charts.getporter.dev
	  name: web
	  version: v0.11.0
  target:
    project: 1234
    cluster: 4321
	namespace: default
  config:
    port: 8080D
```

## RDS + Helm Chart
This showcases a set of resources where an RDS database is created first, and output values from the RDS database are passed to the child resource, in this case a Helm chart that creates a web application. 

```yaml
version: v1
resources:
- name: rds
  driver: terraform
  source:
    kind: github
	config:
	  url: https://github.com/porter-dev/rds
	  path: ./ha
	  ref: v0.11.0
  target:
    project: 1234
  config:
    vpc: vpc-id-1234
	storage: 10GB
  resources:
  - name: web-example
    driver: helm
    source:
      kind: helm_registry
	  config:
	    url: https://porter-charts.getporter.dev
	    name: web
	    version: v0.11.0
    target:
      project: 1234
      cluster: 4321
	  namespace: default
    config:
      port: 8080
	  environment:
	  - name: DB_HOST
	    value: "{ .rds.db_host }"
	  - name: DB_USER
	    value: "{ .rds.db_user }"
	  - name: DB_PASSWORD
	    value: "{ .rds.db_password }"
```

## Application Build + Helm Chart
This showcases an example where an application build is declared, which builds a Docker image whose output reference is passed to Helm chart.   

```yaml
version: v1
resources:
- name: image_build
  driver: builder
  target:
    project: 1234
	registry: 4321
  config:
    method: buildpack
  resources:
  - name: web-example
    driver: helm
    source:
      kind: helm_registry
	  config:
	    url: https://porter-charts.getporter.dev
	    name: web
	    version: v0.11.0
    target:
      project: 1234
      cluster: 4321
	  namespace: default
    config:
      port: 8080
	  image:
	    repository: "{{ .image_build.image }}"
```
