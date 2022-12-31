# ledger

## TL;DR;

```console
helm repo add ledger https://stenic.github.io/ledger/
helm install ledger --namespace ledger ledger/ledger
```

## Introduction

This chart installs `ledger` on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.12+
- Helm 3.0+

## Installing the Chart

To install the chart with the release name `my-release`:

```console
helm repo add ledger https://stenic.github.io/ledger/
helm install ledger --namespace ledger ledger/ledger
```

These commands deploy ledger on the Kubernetes cluster in the default configuration. The [Parameters](#parameters) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```console
helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following tables list the configurable parameters of the ledger chart and their default values.

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Affinity and anti-affinity |
| autoscaling.enabled | bool | `false` | Enable autoscaling |
| autoscaling.maxReplicas | int | `4` | Maximum number of replicas |
| autoscaling.minReplicas | int | `2` | Minimum number of replicas |
| autoscaling.targetCPUUtilizationPercentage | int | `80` | Target CPU utilization percentage |
| fullnameOverride | string | `""` | String to fully override fullname |
| image.pullPolicy | string | `"Always"` | ledger image pullPolicy |
| image.repository | string | `"ghcr.io/stenic/ledger"` | ledger image repository |
| image.tag | string | `""` | ledger image tag (immutable tags are recommended) Overrides the image tag whose default is the chart appVersion. |
| imagePullSecrets | list | `[]` | Docker registry secret names as an array |
| ingress.annotations | object | `{}` | Additional ingress annotations |
| ingress.className | string | `""` | Defines which ingress controller will implement the resource |
| ingress.enabled | bool | `false` | Enable an ingress resource |
| ingress.hosts | list | `[{"host":"chart-example.local","paths":[{"path":"/","pathType":"ImplementationSpecific"}]}]` | List of ingress hosts |
| ingress.tls | list | `[]` | Ingress TLS configuration |
| ledgerAgent.enabled | bool | `true` |  |
| ledgerAgent.endpoint | string | `"http://ledger-server"` |  |
| ledgerAgent.location | string | `""` |  |
| ledgerAgent.log.level | string | `"info"` |  |
| ledgerAgent.namespaced | bool | `false` |  |
| ledgerAgent.token | string | `nil` |  |
| ledgerServer.enabled | bool | `true` |  |
| ledgerServer.jwt.secret | string | `""` |  |
| ledgerServer.log.level | string | `"info"` |  |
| ledgerServer.mysql.dsn | string | `""` |  |
| ledgerServer.oidc.audience | string | `"account"` |  |
| ledgerServer.oidc.clientID | string | `"ledger"` |  |
| ledgerServer.oidc.issuerURL | string | `"https://[keycloak-url]/auth/realms/[realm-name]"` |  |
| nameOverride | string | `""` | String to partially override fullname |
| nodeSelector | object | `{}` | Node labels for controller pod assignment |
| podAnnotations | object | `{}` | Additional annotations for the pods. |
| podSecurityContext | object | `{}` | Enable Controller pods' Security Context |
| replicaCount | int | `2` | Desired number of pods |
| resources | object | `{"limits":{"cpu":"500m","memory":"64Mi"},"requests":{"cpu":"100m","memory":"64Mi"}}` | Resource requests and limits for the controller |
| securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"readOnlyRootFilesystem":true,"runAsGroup":100012,"runAsNonRoot":true,"runAsUser":100012}` | Enable Controller containers' Security Context |
| service.port | int | `80` | Service port |
| service.type | string | `"ClusterIP"` | Kubernetes Service type |
| serviceAccount.annotations | object | `{}` |  |
| serviceAccount.create | bool | `true` | Specifies whether a ServiceAccount should be created |
| serviceAccount.name | string | `""` | The name of the ServiceAccount to use. If not set and create is true, a name is generated using the fullname template |
| tolerations | list | `[]` | Node tolerations for server scheduling to nodes with taints |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`.

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```console
helm install my-release -f values.yaml ledger/ledger
```