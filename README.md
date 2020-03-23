![Go](https://github.com/mchirico/mpubsub/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/mchirico/mpubsub/branch/master/graph/badge.svg)](https://codecov.io/gh/mchirico/mpubsub)



# mpubsub


Testing uses the montcoio Google project




## Build with vendor
```
export GO111MODULE=on
go mod init
# Below will put all packages in a vendor folder
go mod vendor



go test -v -mod=vendor ./...

# Don't forget the "." in "./cmd/script" below
go build -v -mod=vendor ./...
```


## Don't forget golint

```

golint -set_exit_status $(go list ./... | grep -v /vendor/)

```


# mpubsub
