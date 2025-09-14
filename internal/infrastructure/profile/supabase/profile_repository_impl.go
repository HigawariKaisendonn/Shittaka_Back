package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"Shittaka_back/internal/domain/profile/entities"
	"Shittaka_back/internal/domain/profile/repositories"
	"Shittaka_back/internal/domain/shared"
)

// ProfileRepositoryImpl はSupabaseを使用したProfileRepositoryの実装
type ProfileRepositoryImpl struct{}

// NewProfileRepository は新しいProfileRepositoryImplを作成
func NewProfileRepository() repositories.ProfileRepository {
	return &ProfileRepositoryImpl{}
}

// GetByID はIDでプロフィールを取得
func (r *ProfileRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.Profile, error) {
	url := fmt.Sprintf("%s/rest/v1/profiles?id=eq.%s", os.Getenv("SUPABASE_URL"), id)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("find profile failed with status %d: %s", resp.StatusCode, string(body))
	}

	var profileList []map[string]interface{}
	if err := json.Unmarshal(body, &profileList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(profileList) == 0 {
		return nil, shared.NewDomainError("NOT_FOUND", "profile not found")
	}

	return mapToProfile(profileList[0]), nil
}

// Create は新しいプロフィールを作成
func (r *ProfileRepositoryImpl) Create(ctx context.Context, profile *entities.Profile) (*entities.Profile, error) {
	profileData := map[string]interface{}{
		"id":       profile.ID,
		"Username": profile.Name,
	}

	jsonData, err := json.Marshal(profileData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal profile data: %w", err)
	}

	url := os.Getenv("SUPABASE_URL") + "/rest/v1/profiles"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("Prefer", "return=representation")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create profile failed with status %d: %s", resp.StatusCode, string(body))
	}

	var profileList []map[string]interface{}
	if err := json.Unmarshal(body, &profileList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(profileList) == 0 {
		return nil, fmt.Errorf("no profile returned from create operation")
	}

	return mapToProfile(profileList[0]), nil
}

// Update はプロフィールを更新
func (r *ProfileRepositoryImpl) Update(ctx context.Context, profile *entities.Profile) error {
	profileData := map[string]interface{}{
		"Username": profile.Name,
	}

	jsonData, err := json.Marshal(profileData)
	if err != nil {
		return fmt.Errorf("failed to marshal profile data: %w", err)
	}

	url := fmt.Sprintf("%s/rest/v1/profiles?id=eq.%s", os.Getenv("SUPABASE_URL"), profile.ID)
	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("update profile failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// mapToProfile は map[string]interface{} を Profile エンティティに変換
func mapToProfile(m map[string]interface{}) *entities.Profile {
	return &entities.Profile{
		ID:   getString(m, "id"),
		Name: getString(m, "Username"),
	}
}

// getString は map から文字列を安全に取得
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}