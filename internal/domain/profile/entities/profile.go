package entities

import (
	"Shittaka_back/internal/domain/shared"
)

// Profile はユーザープロフィールのドメインエンティティ
type Profile struct {
	ID   string `json:"id"`   // Supabase Authのuuidと同一
	Name string `json:"name"` // 表示ユーザー名
}

// NewProfile は新しいProfileエンティティを作成
func NewProfile(id, name string) *Profile {
	return &Profile{
		ID:   id,
		Name: name,
	}
}

// Validate はProfileエンティティのバリデーションを行う
func (p *Profile) Validate() error {
	if p.ID == "" {
		return shared.NewValidationError("id", "id is required")
	}
	if p.Name == "" {
		return shared.NewValidationError("name", "name is required")
	}
	return nil
}