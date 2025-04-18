package auth

type UserSessionData struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Provider  string `json:"provider"`
	Avatar    string `json:"avatar"`
}
