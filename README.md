# Ledger

[![Last release](https://github.com/stenic/ledger/actions/workflows/release.yaml/badge.svg)](https://github.com/stenic/ledger/actions/workflows/release.yaml)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/ledger)](https://artifacthub.io/packages/helm/ledger/ledger)

> Ledger allows you to record versions of applications you manage. It gives insights in your team's
> deployment habits.

```mermaid
flowchart TB
  subgraph ledger [ ]
    ui(server-ui) --> server(server-api)
    admin(admin-cli) --> server
    server --> database[(Database)]
  end
  subgraph external
    client(client-cli) ---> server
    agent(kubernetes-agent) ---> server
    curl([curl]) ---> server
  end
```

## Installation

```bash
helm repo add ledger https://stenic.github.io/ledger/
helm install ledger --namespace ledger ledger/ledger
```

Check the [helm chart](./charts/ledger/) for more configuration options.

## Collecting versions

Ledger lets you choose how you want to record versions.

- Client CLI
- Kubernetes agent
- Direct GraphQL API

### Client CLI

The client CLI simplifies the interaction with the API. This requires a TOKEN to be used, you
can use an OIDC token or generate a long-lived token using the admin CLI.

```bash
curl -o ledger https://[ledger-installation]/download
chmod +x ledger
./ledger client new-version app env version
```

### Kubernetes Agent

It will collect changes to Deployments and Statefulsets. When a change is detected, the agent will use the
image name without the repository as the application and the tag as the version. This can be overwritten
by using setting one or more of the annotations below.

| Aannotation                  | Description                    | Default                          |
| ---------------------------- | ------------------------------ | -------------------------------- |
| ledger.stenic.io/location    | Overwrite the location         | Set by the agent                 |
| ledger.stenic.io/environment | Overwrite the environment      | Resource namespace               |
| ledger.stenic.io/application | Overwrite the application name | The image without the repository |

### GraphQL API

Communication between any system is possible using the graphql API. This requires a TOKEN to be used, you
can use an OIDC token or generate a long-lived token using the admin CLI.

```bash
export TOKEN=ey...

export version=1.2.3

curl 'https://ledger.development.tbnlabs.be/query' \
 -H "authorization: Bearer $TOKEN" \
  -H 'content-type: application/json' \
  --data-raw "{\"query\":\"mutation { createVersion( input: { application:\\\"$version\\\", environment:\\\"test\\\", version:\\\"$version\\\" } ) { id } }\"}" \
 --compressed \
 --insecure
```

## Administration

```
kubectl exec -ti svc/ledger-server /app/ledger admin -h
```
