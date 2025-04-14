package auth

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Login    string `jsom:"login"`
	Password string `json:"-"`
}

type Token struct {
	Token string `json:"token"`
}
