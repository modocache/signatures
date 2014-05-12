# signatures

A simple RESTful API in Go.

See [this blog post](http://modocache.svbtle.com/restful-go) for an explanation on how to build your own.

## Running the Server

Make sure you have MongoDB installed and running on a standard port.

```
src/signatures/ $ go install
src/signatures/ $ signatures
[martini] listening on :3000 (development)
```

## Running the Tests

You'll need MongoDB running for these as well.

```
src/signatures/ $ ginkgo -r --randomizeAllSpecs -cover
```
