package gou

func Error(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
