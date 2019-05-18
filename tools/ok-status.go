package tools

func IsOkStatus(status int) bool {
	if status == 200 || status == 203 {
		return true
	}

	return false
}
