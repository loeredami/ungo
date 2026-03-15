package ungo

import "testing"

type TestStringRegistry struct{}

func (tsr *TestStringRegistry) Value(r *Registry[string]) string {
	return "bar"
}

func TestRegistry(t *testing.T) {
	tsr := &TestStringRegistry{}
	r := NewRegistry[string]()
	r.Register("foo", tsr)
	v := r.Get("foo")
	if v != "bar" {
		t.Errorf("expected bar, got %s", v)
	}

	v = r.GetOrDefault("bar", "")
	if v != "" {
		t.Errorf("expected empty, got %s", v)
	}

}
