package str

// Error truns err to string
func Error(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
