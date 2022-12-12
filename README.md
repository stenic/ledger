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
