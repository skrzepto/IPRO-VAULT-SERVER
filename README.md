# IPRO-VAULT-SERVER

## Installing

Requirements:
Go >=1.4


*NOTE: Make sure your $GOPATH is set*

``
go get -u github.com/skrzepto/IPRO-VAULT-SERVER
``

## Running the program

``
go run $GOPATH/src/github.com/skrzepto/IPRO-VAULT-SERVER/app.go
``

## Creating a binary

```
cd $GOPATH/src/github.com/skrzepto/IPRO-VAULT-SERVER/
go install
IPRO-VAULT-SERVER
```

now go to localhost:8082

## Running test suite
``
cd $GOPATH/src/github.com/skrzepto/IPRO-VAULT-SERVER/
go test
``
