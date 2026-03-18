package ungo

type Disposable interface {
	Dispose()
}

type ResourceArena struct {
	resources []Disposable
}

func NewResourceArena() *ResourceArena {
	return &ResourceArena{}
}

func (a *ResourceArena) Track(d Disposable) {
	a.resources = append(a.resources, d)
}

func (a *ResourceArena) Melt() {
	for _, r := range a.resources {
		r.Dispose()
	}
	a.resources = make([]Disposable, 0)
}
