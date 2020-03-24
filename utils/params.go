package utils

type Params map[string]interface{}

func (p Params) AddAll(p1 Params) Params {
	for k, v := range p1 {
		p[k] = v
	}
	return p
}

func (p Params) Err(err error) Params {
	p["err"] = err
	return p
}

func (p Params) Add(key string, value interface{}) Params {
	p[key] = value
	return p
}
