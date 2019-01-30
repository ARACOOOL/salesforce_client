# SalesForce REST API Golang client

## Installation
```bash
go get github.com/aracoool/salesforce_client
```

## Usage

**Create a client**

```go
// Production 
client := NewClient(EnvProduction, "44.0")

// Staging 
client := NewClient(EnvStaging, "44.0")
```

**Retrieve a specific object**
```go
account := &struct{
	Id,
	Name,
	Status
}{}
client.Find("Account", "0030x00000N1vJ0AAJ", account)

// account.Id == 0030x00000N1vJ0AAJ
```