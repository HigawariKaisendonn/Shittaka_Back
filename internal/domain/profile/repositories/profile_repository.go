package repositories

import (
	"context"
	"Shittaka_back/internal/domain/profile/entities"
)

// ProfileRepository はプロフィールのリポジトリインターフェース
type ProfileRepository interface {
	// GetByID はIDでプロフィールを取得
	GetByID(ctx context.Context, id string) (*entities.Profile, error)

	// Create は新しいプロフィールを作成
	Create(ctx context.Context, profile *entities.Profile) (*entities.Profile, error)

	// Update はプロフィールを更新
	Update(ctx context.Context, profile *entities.Profile) error
}