package auth

type UserSessionData struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Provider  string `json:"provider"`
}
