package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"Shittaka_back/internal/application/profile/dto"
	"Shittaka_back/internal/application/profile/usecases"
	"Shittaka_back/internal/domain/shared"
	presentationDTO "Shittaka_back/internal/presentation/dto"
)

// ProfileHandler はプロフィール関連のHTTPハンドラー
type ProfileHandler struct {
	profileUsecase *usecases.ProfileUsecase
}

// NewProfileHandler は新しいProfileHandlerを作成
func NewProfileHandler(profileUsecase *usecases.ProfileUsecase) *ProfileHandler {
	return &ProfileHandler{
		profileUsecase: profileUsecase,
	}
}

// GetProfileHandler はプロフィールを取得
func (h *ProfileHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLパスからuser_idを取得 (/api/profiles/{user_id})
	userID := r.URL.Path[len("/api/profiles/"):]
	if userID == "" {
		h.sendError(w, "User ID is required", http.StatusBadRequest)
		return
	}

	profile, err := h.profileUsecase.GetProfile(r.Context(), userID)
	if err != nil {
		h.handleUsecaseError(w, err)
		return
	}

	// レスポンスDTOに変換
	response := presentationDTO.ProfileResponse{
		ID:   profile.ID,
		Name: profile.Name,
	}

	h.sendJSON(w, response, http.StatusOK)
}

// CreateProfileHandler はプロフィールを作成
func (h *ProfileHandler) CreateProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req presentationDTO.CreateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// DTOの変換
	usecaseReq := dto.CreateProfileRequest{
		ID:   req.ID,
		Name: req.Name,
	}

	profile, err := h.profileUsecase.CreateProfile(r.Context(), usecaseReq)
	if err != nil {
		h.handleUsecaseError(w, err)
		return
	}

	// レスポンスDTOに変換
	response := presentationDTO.ProfileResponse{
		ID:   profile.ID,
		Name: profile.Name,
	}

	h.sendJSON(w, response, http.StatusCreated)
}

// UpdateProfileHandler はプロフィールを更新
func (h *ProfileHandler) UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLパスからuser_idを取得
	userID := r.URL.Path[len("/api/profiles/"):]
	if userID == "" {
		h.sendError(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var req presentationDTO.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// DTOの変換
	usecaseReq := dto.UpdateProfileRequest{
		Name: req.Name,
	}

	profile, err := h.profileUsecase.UpdateProfile(r.Context(), userID, usecaseReq)
	if err != nil {
		h.handleUsecaseError(w, err)
		return
	}

	// レスポンスDTOに変換
	response := presentationDTO.ProfileResponse{
		ID:   profile.ID,
		Name: profile.Name,
	}

	h.sendJSON(w, response, http.StatusOK)
}

// ヘルパー関数

// handleUsecaseError はユースケースエラーを適切なHTTPエラーに変換
func (h *ProfileHandler) handleUsecaseError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case shared.ValidationError:
		h.sendError(w, e.Message, http.StatusBadRequest)
	case shared.DomainError:
		switch e.Code {
		case "NOT_FOUND":
			h.sendError(w, e.Message, http.StatusNotFound)
		default:
			h.sendError(w, e.Message, http.StatusInternalServerError)
		}
	default:
		log.Printf("Profile usecase error: %v", err)
		h.sendError(w, "Internal server error", http.StatusInternalServerError)
	}
}

// sendJSON はJSONレスポンスを送信
func (h *ProfileHandler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

// sendError はエラーレスポンスを送信
func (h *ProfileHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	response := presentationDTO.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	h.sendJSON(w, response, statusCode)
}