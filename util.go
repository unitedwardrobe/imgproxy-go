package imgproxy

func boolAsNumberString(i bool) string {
	if i {
		return "1"
	}

	return "0"
}
