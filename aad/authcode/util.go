package authcode

// stringChange is a string which can tell you if it is changed when new value comes in
type stringChange string

// Changed receives the new string value and compares it with the current value.
// It will return true if they are not the same, and then update the current value.
func (c *stringChange) Changed(s string) bool {
	old := string(*c)
	if old != s {
		*c = stringChange(s)
		return true
	}
	return false
}

func stringSliceContains(sl []string, s string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}
