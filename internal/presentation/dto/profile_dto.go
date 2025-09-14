package dto

// ProfileResponse はプロフィールレスポンスのHTTP DTO
type ProfileResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreateProfileRequest はプロフィール作成リクエストのHTTP DTO
type CreateProfileRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UpdateProfileRequest はプロフィール更新リクエストのHTTP DTO
type UpdateProfileRequest struct {
	Name string `json:"name"`
}