# SalesForce REST API Golang client

[![Build Status](https://travis-ci.org/ARACOOOL/salesforce_client.svg?branch=master)](https://travis-ci.org/ARACOOOL/salesforce_client)

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

**Authentication**

You have to call the `Auth()` method after you created a client
```go
err := client.Auth(Auth{
	"Username",
	"Password",
	"ClientID",
	"ClientSecret"
})
```

**Retrieve a specific object**
```go
account := &struct{
	Id,
	Name,
	Status
}{}
err := client.Find("Account", "0030x00000N1vJ0AAJ", account)

// account.Id == 0030x00000N1vJ0AAJ
```

**Update a specific object**
```go
params := &Params{}
params.AddField("Name", "New name")
err := client.Update("Account", "0030x00000N1vJ0AAJ", params)
```

**Delete a specific object**
```go
err := client.Delete("Account", "0030x00000N1vJ0AAJ")
```

**Make a query**
```go
builder := &SoqlBuilder{}
builder.Select("Id", "Name")
builder.From("Account")
builder.Where("Name='Test'")
builder.Limit(1)

accounts := &struct{
	Records []struct{
		Id string,
		Name string
	} `json:"records"`
}{}
err := client.Query(builder.Build(), accounts)
```