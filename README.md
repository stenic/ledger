## Collecting versions

### Sending by installing the included agent in your cluster

It will collect changes to Deployments and Statefulsets

### Sending using the included client

```bash
curl -o ledger https://[ledger-installation]/download
chmod +x ledger
./ledger client new-version app env version
```

### Sending using the GraphQL API

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
