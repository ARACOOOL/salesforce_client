package salesforce_client

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSoqlBuilder_Build(t *testing.T) {
	test := assert.New(t)

	builder := &SoqlBuilder{}
	builder.Select("Id", "Name")
	builder.From("Account")
	builder.Where("Id=0010x00000CfFAnAAN")

	test.Equal("SELECT Id,Name FROM Account WHERE Id=0010x00000CfFAnAAN", builder.Build())

	builder.Where("Name=Test")
	builder.Limit(2)

	test.Equal("SELECT Id,Name FROM Account WHERE Name=Test LIMIT 2", builder.Build())
}
