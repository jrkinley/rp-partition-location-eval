# Readme
This project is a small utility designed to understand the partition-location mapping of a Redpanda cluster.

# Usage
To run the tool, use the following:

```shell
go run main.go --seed seed-redacted.fmc.prd.cloud.redpanda.com:9092 --username pmw --password redacted
```

You will receive output similar to the following:

```text
owlshop-frontend-events:0 has more than 1 replica in use2-az1
owlshop-frontend-events:1 has more than 1 replica in use2-az1
owlshop-frontend-events:2 has more than 1 replica in use2-az1
...
```

## Auth
If you require a sasl mechanism of SCRAM_SHA_512 (rather than SCRAM_SHA_256), run the tool as follows:

```shell
go run main.go --seed seed-redacted.fmc.prd.cloud.redpanda.com:9092 --username pmw --password redacted --use512 true
```