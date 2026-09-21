package structs

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	IsActive *bool  `json:"is_active"`
}

type LoginResponse struct {
	AccessToken string     `json:"access_token"`
	TokenType   string     `json:"token_type"`
	ExpiresIn   int64      `json:"expires_in"`
	User        UserPublic `json:"user"`
}

type UserPublic struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}
