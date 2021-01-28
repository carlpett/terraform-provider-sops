package ini

import (
	"gopkg.in/ini.v1"
)

func Unmarshal(in []byte, out *map[string]interface{}) error {
	f, err := ini.Load(in)
	if err != nil {
		return err
	}

	if *out == nil {
		*out = make(map[string]interface{})
	}

	for _, s := range f.Sections() {
		var m map[string]interface{}
		// The root section key-value pairs should go directly on the root of the
		// map, not under a default subkey
		if s.Name() == ini.DefaultSection {
			m = *out
		} else {
			m = make(map[string]interface{})
			(*out)[s.Name()] = m
		}

		for _, k := range s.Keys() {
			m[k.Name()] = k.Value()
		}
	}

	return nil
}
