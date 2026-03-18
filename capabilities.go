package ungo

type Capability interface {
	Type() string
	Is(string) bool
}

type Provider struct {
	reg *Registry[Capability]
}

func NewProvider(reg *Registry[Capability]) *Provider {
	return &Provider{
		reg: reg,
	}
}

func (p *Provider) HasCapability(cap string) bool {
	var found bool
	ForEach(p.reg.Keys(), func(c_name string) {
		c := p.reg.Get(c_name)
		if !c.HasValue() {
			return
		}

		if c.Value().Is(cap) {
			found = true
			return
		}
	})
	return found
}

func (p *Provider) GetCapability(cap string) Optional[Capability] {
	for _, c_name := range p.reg.Keys() {
		c := p.reg.Get(c_name)
		if !c.HasValue() {
			continue
		}

		if c.Value().Is(cap) {
			return c
		}
	}
	return EmptyOptional[Capability]()
}

func (p *Provider) Use(cap string, f func(Capability)) bool {
	c := p.GetCapability(cap)
	var used bool
	c.IfPresent(func(c Capability) {
		f(c)
		used = true
	})
	c.IfAbsent(func(*Capability) {
		used = false
	})
	return used
}
