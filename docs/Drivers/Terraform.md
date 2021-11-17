## Source 
The Terraform driver can specify a source, which is applied as either a Terraform module or a Terraform program. The source can be read from:
1. A local directory
2. A public Github repository containing a Terraform module 
3. A Porter Github integration

### Local Directory

```yaml
source:
  kind: local
  path: ./relative/path
```

### Local Directory

TODO

### Porter Github Integration

TODO

## Target

There is no target configuration for the Terraform driver at the moment. 