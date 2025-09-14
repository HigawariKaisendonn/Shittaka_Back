package dto

// ProfileResponse はプロフィールレスポンス
type ProfileResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreateProfileRequest はプロフィール作成リクエスト
type CreateProfileRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UpdateProfileRequest はプロフィール更新リクエスト
type UpdateProfileRequest struct {
	Name string `json:"name"`
}