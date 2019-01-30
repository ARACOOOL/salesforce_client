package salesforce_client

type Params struct {
	fields map[string]interface{}
}

func (p *Params) AddField(name string, value interface{}) {
	if p.fields == nil {
		p.fields = map[string]interface{}{}
	}

	p.fields[name] = value
}

func (p *Params) GetFields() map[string]interface{} {
	return p.fields
}
