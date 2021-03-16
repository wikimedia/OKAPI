package utils

// FilterNs check whether page in allowed namespace
func FilterNs(ns int) bool {
	return ns == 0
}
