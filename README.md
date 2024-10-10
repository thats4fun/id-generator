### ID Generator

#### Requirements

Should implement the following API:
- [x] Able to generate an ID that has not already been generated and output to the command line.
- [x] Able to free an already generated ID and reuse it the next time an ID is needed.
- [x] Should be concurrency-safe.
- [x] Should be able to be called from a bash script/ the command line.

#### API

#### generate:

```
go run main.go getid
```

#### free:

```
go run main.go freeid [ID]
```