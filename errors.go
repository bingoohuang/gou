package go_utils

func Error(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
