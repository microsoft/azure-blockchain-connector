package proxy

type Proxy struct {
	*Params
}

func (p *Proxy) Handler() Handler {
	return Handler{Params: p.Params}
}
