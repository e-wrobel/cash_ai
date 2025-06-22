## CashCoach AI – Transactions Processor


## Key features
```text
• Pulls transaction batches from a mock provider, on schedule or via HTTP trigger
• Stores immutable transaction events (event‑sourcing) so that any point‑in‑time
state can be rebuilt ("replay")
• Publishes every domain event on an in‑memory Pub/Sub bus to decouple further
asynchronous processing (e.g. categorisation, ML, …)
• Clean, semi-hexagonal layout: core domain has no external dependencies; adapters live
under internal/adapter; wiring happens in cmd/processor/main.go
```

## Running locally
```bash
  make install  # Install linter and code formatter tools
  make fix      # Automatically fix linter errors and format code
  make build    # build executable for ./bin/processor
  make run      # starts  processor
  make mock     # starts mock provider only
  make test     # starts unit tests
```

Remark: Remember to run the mock provider first, as the processor depends on it.
## Usage
```bash
  curl -XPOST localhost:8080/trigger/fetch   # manual
  curl -XPOST localhost:8080/replay          # rebuild state from zero
```

## Environment
```text
PROVIDER_URL   – http://localhost:9000 (default)
PULL_INTERVAL  – 30s (Go duration, optional)
```

## Endpoints
```text

/trigger/fetch (POST)
Purpose: Manually or periodically triggers a fetch of transactions from the mock provider.
Steps:
	1.	Fetch transactions via HTTP from external provider.
	2.	Apply logic: compare with current state.

        New → created event
        Changed → updated event
        Missing locally → deleted event
	3.	Persist all events and publish via Pub/Sub.

/replay (POST)
Purpose: Rebuilds in-memory state from the complete event history.
Steps:
	1.	Fetch all events from EventStore.
	2.	Re-apply events in order:
	Created / updated → update state
	Deleted → remove from state
```

## Sequence diagram and swagger spec
Please take a look at the sequence diagram in the `docs` folder for a visual representation of the process flow.