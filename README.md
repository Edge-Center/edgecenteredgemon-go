# edgecenteredgemon-go

## Testing

Unit tests (mock the `Requester`) and in-process integration tests run together:

```sh
make test
```

The `integration/` package wires the real `provider.Client` to an in-process
fake RMON server and drives every service through a full CRUD cycle - no network
or credentials required.

### Live integration tests

Tests against a real RMON API live behind the `integration` build tag and skip
themselves unless credentials are provided:

```sh
EDGECENTER_RMON_BASE_URL=https://api.example.com \
EDGECENTER_RMON_API_KEY=xxxxxxxx \
make test-integration
```
