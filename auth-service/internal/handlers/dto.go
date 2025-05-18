package handlers

// ====== Structures ======
type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type refreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type sdkTokenResponse struct {
	AccessToken string `json:"access_token"`
}
