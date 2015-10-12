package utils

// create gravatar link
func Gravatar(email string) string {
	return "https://gravatar.com/avatar/" + Md5String(email)
}
