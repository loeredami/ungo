package ungo

type Disposable interface {
	Dispose()
}

type Arena struct {
	resources []Disposable
}

func NewArena() *Arena {
	return &Arena{}
}

func (a *Arena) Track(d Disposable) {
	a.resources = append(a.resources, d)
}

func (a *Arena) Melt() {
	for _, r := range a.resources {
		r.Dispose()
	}
	a.resources = make([]Disposable, 0)
}
