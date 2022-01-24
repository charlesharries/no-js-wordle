package forms

// Define a new errors type, which we will use to hold the validation error
// messages for forms. The name of the form field will be used as the key
// in this map.
type errors map[string][]string

// Add adds an error message to a given field.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Implement a Get() method to retrieve the first error message for a given
// field from the map.
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}

	return es[0]
}
