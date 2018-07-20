# BBFT
Blockchain Byzantine Fault Torelance Consensus Algorithm based PBFT.

## environement
- go 1.10.3
- glide 0.13.1
- libprotoc 3.6.0

## previous install
```
$ glide install
```

## Demo

### Demo Server SetUp : (4 Peers)
```
$ docker-compose up
```

### Demo Transaction Send
```
$ make build-sender
$ ./bin/sender
```