package util

// stringChange is a string which can tell you if it is changed when new value comes in
type StringChange string

// Changed receives the new string value and compares it with the current value.
// It will return true if they are not the same, and then update the current value.
func (c *StringChange) Changed(s string) bool {
	old := string(*c)
	if old != s {
		*c = StringChange(s)
		return true
	}
	return false
}

func StringSliceContains(sl []string, s string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}

type StringKVP struct {
	K, V string
}
