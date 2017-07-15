package env2flag

import (
	"flag"
	"os"
	"strings"
)

// Mapper provides a mapping capability.
type Mapper interface {
	// Map maps a string to another one.
	Map(name string) (string, bool)
}

type mappers []Mapper

func (maps mappers) Map(s string) (string, bool) {
	var ok bool
	if len(maps) == 0 {
		return "", false
	}
	for _, m := range maps {
		s, ok = m.Map(s)
		if !ok {
			return "", false
		}
	}
	return s, ok
}

// MapFunc is shorthand for Mapper.
type MapFunc func(string) (string, bool)

// Map maps a string using MapFunc.
func (f MapFunc) Map(name string) (string, bool) {
	return f(name)
}

var (
	// DefaultMapper provides standard rule to map argument names to
	// environment names.
	DefaultMapper Mapper = MapFunc(standardMap)

	// EnvMapper provides a map from name to environment value.
	EnvMapper = MapFunc(os.LookupEnv)
)

func standardMap(name string) (string, bool) {
	suppress := true
	n := strings.Map(func(r rune) rune {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			suppress = false
			return r
		}
		if suppress {
			return -1
		}
		suppress = true
		return '_'
	}, strings.ToUpper(name))
	return n, true
}

// ApplyMaps applies all maps to flag.FlagSet.
func ApplyMaps(fs *flag.FlagSet, maps ...Mapper) {
	if len(maps) == 0 {
		return
	}
	m := mappers(maps)
	fs.VisitAll(func(f *flag.Flag) {
		if v, ok := m.Map(f.Name); ok {
			f.Value.Set(v)
		}
	})
}

// Parse parses flag with considering environment variables.
func Parse() {
	ApplyMaps(flag.CommandLine, DefaultMapper, EnvMapper)
	flag.Parse()
}

// Map is a wrapper to use map[string]string as Mapper.
type Map map[string]string

// Map maps a string.
func (m Map) Map(s string) (string, bool) {
	if m == nil {
		return "", false
	}
	v, ok := m[s]
	return v, ok
}

// Parse parses flag with this map.
func (m Map) Parse(fs *flag.FlagSet, args []string) {
	ApplyMaps(fs, m, EnvMapper)
	fs.Parse(args)
}
