package salesforce_client

import "strconv"

type SoqlBuilder struct {
	fields string
	from   string
	where  string
	limit  int
}

func (b *SoqlBuilder) Select(fields ...string) {
	for _, field := range fields {
		if b.fields == "" {
			b.fields += field
		} else {
			b.fields += "," + field
		}
	}
}

func (b *SoqlBuilder) From(object string) {
	b.from = object
}

func (b *SoqlBuilder) Where(condition string) {
	b.where = condition
}

func (b *SoqlBuilder) Limit(limit int) {
	b.limit = limit
}

func (b *SoqlBuilder) Build() string {
	query := "SELECT " + b.fields + " FROM " + b.from
	if b.where != "" {
		query += " WHERE " + b.where
	}

	if b.limit != 0 {
		query += " LIMIT " + strconv.Itoa(b.limit)
	}

	return query
}
