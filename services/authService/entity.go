package authService

type RegistrationInputBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginInputBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegistrationConfirmationInputBody struct {
	Username string `json:"username"`
	Code     string `json:"code"`
}

type Claims struct {
	Sub       string `json:"id"`
	Iss       string `json:"iss"`
	ClientID  string `json:"client_id"`
	OriginJTI string `json:"origin_jti"`
	EventID   string `json:"event_id"`
	TokenUse  string `json:"token_use"`
	Scope     string `json:"scope"`
	AuthTime  int    `json:"auth_time"`
	Exp       int    `json:"exp"`
	Iat       int    `json:"iat"`
	Jti       string `json:"jti"`
	Username  string `json:"username"`
}
