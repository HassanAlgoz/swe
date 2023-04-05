package user

func validateUsername(username string) bool {
	if len(username) < 3 {
		return false
	}
	return true
}
