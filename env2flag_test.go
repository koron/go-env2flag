package env2flag

import (
	"flag"
	"testing"
)

func TestStandardMap(t *testing.T) {
	ok := func(s, exp string) {
		act, ok := standardMap(s)
		if !ok {
			t.Fatalf("standardMap failed to map: %q", s)
		}
		if act != exp {
			t.Fatalf("standardMap mapped to unmatch: %q -> %q (expected:%q)", s, act, exp)
		}
	}

	ok("abc", "ABC")
	ok("ABC", "ABC")
	ok("Abc", "ABC")

	ok("abc123", "ABC123")
	ok("abc_123", "ABC_123")
	ok("abc-123", "ABC_123")
	ok("abc'123", "ABC_123")
	ok("abc#123", "ABC_123")

	ok("abc____123", "ABC_123")
	ok("abc-_-_123", "ABC_123")
	ok("____abc_123", "ABC_123")
	ok("abc____", "ABC_")
	ok("abc----", "ABC_")
}

func TestMap(t *testing.T) {
	mapped := func (m Mapper, k, exp string) {
		act, ok := m.Map(k)
		if !ok {
			t.Errorf("not mapped %q", k)
			return
		}
		if act != exp {
			t.Errorf("not match %q: %q (expected %q)", k, act, exp)
			return
		}
	}
	notMapped := func(m Mapper, k string) {
		v, ok := m.Map(k)
		if ok {
			t.Errorf("mapped %q", k)
			return
		}
		if v != "" {
			t.Errorf("not empty for %q: %q", k, v)
			return
		}
	}

	m := Map{"foo": "s1", "bar": "s2"}
	mapped(m, "foo", "s1")
	mapped(m, "bar", "s2")
	notMapped(m, "baz")
	notMapped(m, "quux")
}

func TestApplyMaps(t *testing.T) {
	fs := flag.NewFlagSet("cmdname", flag.ContinueOnError)
	v1 := fs.String("abc_123", "v1", "option v1")
	v2 := fs.String("foo", "v2", "option v2")
	v3 := fs.String("bar", "v3", "option v3")
	ApplyMaps(fs, DefaultMapper, Map{
		"ABC_123": "v1a",
		"FOO":     "v2a",
		"bar":     "v3b",
	})
	ok := func(name, value, exp string) {
		if value != exp {
			t.Errorf("*%s didn't match: %q (expected %q)", name, value, exp)
		}
	}
	ok("v1", *v1, "v1a")
	ok("v2", *v2, "v2a")
	ok("v3", *v3, "v3")
}
