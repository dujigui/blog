package utils

type Params map[string]interface{}

func (p Params) AddAll(p1 Params) Params {
	for k, v := range p1 {
		p[k] = v
	}
	return p
}
