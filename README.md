# Credhub Resource [![Concourse Resource](https://ci.starkandwayne.com/api/v1/pipelines/credhub-resource/jobs/test/badge)](https://ci.starkandwayne.com/teams/main/pipelines/credhub-resource)

A [concourse](https://concourse.ci) resource for [Credhub](https://github.com/cloudfoundry-incubator/credhub).

## Adding to your pipeline

To use the Credhub Resource, you must declare it in your pipeline as a resource type:

```
resource_types:
- name: credhub
  type: docker-image
  source:
    repository: cfcommunity/credhub-resource
```

## Source Configuration

* `server`: *Required.* The address of the Credhub server
* `username`: *Required.* The UAA client ID for authorizing with Credhub.
* `password`: *Required.* The UAA client secret for authorizing with Credhub.
* `path`: *Optional.* The Credhub path which needs to be watched.
* `skip_tls_validation`: *Optional.* Skip TLS validation for connections to Credhub and UAA.

### Example

``` yaml
- name: credhub
  type: credhub
  source:
    server: https://credhub.example.com
    username: admin
    password: admin
    path: /bosh/cf/
    skip_tls_validation: true
```

## Behaviour

This resource will create a new version when a credential Credhub changes.

### Example

``` yaml
jobs:
- name: deploy_cf
  plan:
  - aggregate:
  - get: cf-deployment
  - get: credhub
    trigger: true
...
```
