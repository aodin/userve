package server

import ()

type Attrs map[string]interface{}

func (a Attrs) Merge(b map[string]interface{}) {
	// Do not overwrite attrs, but overwriting b is fine
	for key, value := range a {
		b[key] = value
	}
}
