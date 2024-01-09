package entity

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
