# gRPCCalc

This is a gRPC-based microservice, that provides two methods: add and multiplyF.
The add method allows adding two integers and the multiplyF method allows
multiplying two floating point numbers.

The microservice can be deployed to Kubernetes, with the help of the provided
Dockerfile and the manifests in the kubernetes directory.

## Service architecture

The repository contains both a server and a client, within cmd/server/ and cmd/client/,
respectively. The client, which is for testing the service, is a CLI program that
supports two subcommands: add and multiplyF. For an example of using it, see the
section *Running*.

## Requirements

* Go >= v1.14
* protobuf: `brew install protobuf && go get -u github.com/golang/protobuf/protoc-gen-go`

## Running

```
go generate ./...
go run cmd/server/main.go
# In another terminal
go run cmd/client/main.go add 1 2
go run cmd/client/main.go multiplyF 1.2 2.5
```

## Testing

The service has a test suite in cmd/server/server_test.go and cmd/client/client_test.go. Run it as follows:

```
go test ./...
```

## Continuous integration

The project is continuously tested, on every push, through GitHub Actions.

## Building Docker image

To build the Docker image ariannavespri/grpccalc run the following:
```
./build-image.sh
```

## Deploying to Kubernetes

There are manifests in the *kubernetes* directory, that deploy the
ariannavespri/grpccalc image to Kubernetes. To deploy them to your cluster, run
the following:

```
kubectl apply -r -f kubernetes
```

This will create a Deployment and corresponding Service in the cluster.

## External access to the service

The service is currently not accessible to clients outside of the cluster.
To enable such access, one would have to add an ingress controller and a service
ingress to the Kubernetes cluster. I've decided it's not necessary for now.

### Expanding the service to use an event store

The service could conceptually be expanded to use an event store, for example by
logging every call in the event store. Technically, this could be supported by
detecting if the environment variable EVENTSTORE_URL were set. If this were the
case, the service could connect to the event store at said URL.

## 12 factor compliance

1. There is just one codebase for the service.
2. Dependencies are explicitly declared in go.mod and the README.
3. Configuration variables are read from the environment.
4. There are no backing services.
5. Build and run stages are separate through building a Docker image and
    deploying it to Kubernetes.
6. The service is stateless.
7. The service binds a port.
8. The service allows for concurrency.
9. The service is easily disposable.
10. The service has dev/prod parity through the use of Kubernetes.
11. The service logs in the standard manner.
12. Admin tasks are included as scripts.

## Cloud native

The service is well-adapted to being deployed in Kubernetes, for example by
being stateless, so it is in its nature cloud native.
