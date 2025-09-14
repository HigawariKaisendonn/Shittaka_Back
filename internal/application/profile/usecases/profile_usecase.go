package usecases

import (
	"context"

	"Shittaka_back/internal/application/profile/dto"
	"Shittaka_back/internal/domain/profile/entities"
	"Shittaka_back/internal/domain/profile/repositories"
)

// ProfileUsecase はプロフィールに関するユースケース
type ProfileUsecase struct {
	profileRepo repositories.ProfileRepository
}

// NewProfileUsecase は新しいProfileUsecaseを作成
func NewProfileUsecase(profileRepo repositories.ProfileRepository) *ProfileUsecase {
	return &ProfileUsecase{
		profileRepo: profileRepo,
	}
}

// GetProfile はプロフィールを取得する
func (u *ProfileUsecase) GetProfile(ctx context.Context, id string) (*dto.ProfileResponse, error) {
	profile, err := u.profileRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.ProfileResponse{
		ID:   profile.ID,
		Name: profile.Name,
	}, nil
}

// CreateProfile は新しいプロフィールを作成する
func (u *ProfileUsecase) CreateProfile(ctx context.Context, req dto.CreateProfileRequest) (*dto.ProfileResponse, error) {
	// プロフィールエンティティを作成
	profile := entities.NewProfile(req.ID, req.Name)

	// エンティティレベルでのバリデーション
	if err := profile.Validate(); err != nil {
		return nil, err
	}

	// リポジトリに保存
	createdProfile, err := u.profileRepo.Create(ctx, profile)
	if err != nil {
		return nil, err
	}

	return &dto.ProfileResponse{
		ID:   createdProfile.ID,
		Name: createdProfile.Name,
	}, nil
}

// UpdateProfile はプロフィールを更新する
func (u *ProfileUsecase) UpdateProfile(ctx context.Context, id string, req dto.UpdateProfileRequest) (*dto.ProfileResponse, error) {
	// 既存のプロフィールを取得
	existingProfile, err := u.profileRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// プロフィールを更新
	existingProfile.Name = req.Name

	// エンティティレベルでのバリデーション
	if err := existingProfile.Validate(); err != nil {
		return nil, err
	}

	// リポジトリで更新
	if err := u.profileRepo.Update(ctx, existingProfile); err != nil {
		return nil, err
	}

	return &dto.ProfileResponse{
		ID:   existingProfile.ID,
		Name: existingProfile.Name,
	}, nil
}