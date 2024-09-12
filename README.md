# GORM DBML Gen

[dbml](https://dbml.dbdiagram.io/home/)

DBML (Database Markup Language) is an open-source DSL language designed to define and document database schemas and structures. It is designed to be simple, consistent and highly-readable.

## Install
```
go install github.com/alonelucky/gen/cmd/dbmlgen@latest
```

## Features
1. Easily record data changes with code.

## Usage of hmq:
```
Usage of dbmlgen
  -dbml string
        database mark languge
  -null
        generate model global configuration
  -output string
        curd and model output dir (default "./gen/curd")
  -type string
        specify data type mapping relationship (default "datetime=string"), exp "coltyp1=int,coltyp2=bool..."
```