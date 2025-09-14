package di

// container.goは依存関係のコンテナを定義
// コンンテナとは、依存関係の解決を行う
// 必要な部品を正しい順で作って配線する工場

import (
	"Shittaka_back/internal/application/auth/usecases"
	profileUsecases "Shittaka_back/internal/application/profile/usecases"
	"Shittaka_back/internal/domain/auth/services"
	"Shittaka_back/internal/infrastructure/auth/supabase"
	"Shittaka_back/internal/infrastructure/config"
	profileSupabase "Shittaka_back/internal/infrastructure/profile/supabase"
	"Shittaka_back/internal/presentation/http/handlers"
)

// Container は依存関係のコンテナ
type Container struct {
	Config         *config.Config
	AuthHandler    *handlers.AuthHandler
	ProfileHandler *handlers.ProfileHandler
}

// NewContainer は新しいコンテナを作成
func NewContainer() *Container {
	// 設定を読み込み
	cfg := config.LoadConfig()

	// 依存関係を構築（外側から内側へ）
	// Auth関連
	userRepo := supabase.NewUserRepository()
	authService := services.NewAuthService(userRepo)
	authUsecase := usecases.NewAuthUsecase(authService)
	authHandler := handlers.NewAuthHandler(authUsecase)

	// Profile関連
	profileRepo := profileSupabase.NewProfileRepository()
	profileUsecase := profileUsecases.NewProfileUsecase(profileRepo)
	profileHandler := handlers.NewProfileHandler(profileUsecase)

	return &Container{
		Config:         cfg,
		AuthHandler:    authHandler,
		ProfileHandler: profileHandler,
	}
}
