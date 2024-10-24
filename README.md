## Generic ETL Pipeline

This repository contains solution that solves the following problem:

Given a set of input files which contain JSOn data of same template with different key-valu pairs and a schema JSON whicn
represent the final format, convert each input file into an output file. The key mappings are provided in another mapping.json file.

This is a typical example of an **ETL - Extract Transform and Load** pipeline. Data pipelining from different data sources, be it a data warehouse or a data lake, need to be refined and transformed before storing it in a standardized format. This problem is a simpler version of exactly this.

[Problem Statement](./problem_statement)

## Repository structure:

```text
pkg/
├── buildjson/
│   ├── schema/
│   │   └── user.json
│   └── json.go
sample_data/
go.mod
main.go
mapping.json
problem_statement
README.md
```

`json.go` contains the implementation for the schema transformer. `mapping.json` contains a single JSON object for respective mapping for each input file. `user.json` is the schema for the final output. Note that nesting of the objects 
is **flattened** out in schema.json.

## Run

To run the pipeline, just head over to terminal and execute `go run main.go`
