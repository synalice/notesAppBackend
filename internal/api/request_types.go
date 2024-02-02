package api

// CredentialsRequest is received everytime user wants to log in or register.
type CredentialsRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// JWTRequest is received everytime there is a need to verify a JWT token.
type JWTRequest struct {
	JWT string `json:"jwt"`
}
