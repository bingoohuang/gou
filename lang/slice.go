package lang

// M1 shims for 1 param return values
func M1(a, b interface{}) []interface{} {
	return []interface{}{a}
}

// M2 shims for 2 param return values
func M2(a, b interface{}) []interface{} {
	return []interface{}{a, b}
}

// M3 shims for 3 param return values
func M3(a, b, c interface{}) []interface{} {
	return []interface{}{a, b, c}
}

// M4 shim for 4 param return values
func M4(a, b, c, d interface{}) []interface{} {
	return []interface{}{a, b, c, d}
}

// M5 shim for 5 param return values
func M5(a, b, c, d interface{}) []interface{} {
	return []interface{}{a, b, c, d}
}
