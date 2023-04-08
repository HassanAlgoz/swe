package user

func isValidUsername(username string) (bool, string) {
	if !(len(username) >= 3) {
		return false, "must be 3+ characters long"
	}
	return true, ""
}
