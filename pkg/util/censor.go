package util

func CencorName(name string) string {
	if len(name) <= 2 {
		return name[:1] + "***"
	}

	firstChar := string(name[0])
	lastChar := string(name[len(name)-1])

	return firstChar + "***" + lastChar
}
