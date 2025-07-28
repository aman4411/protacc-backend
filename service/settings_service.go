package service

import (
	"context"
	"fmt"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/repository"
)

type SettingsService struct {
	repo *repository.SettingsRepository
}

func NewSettingsService(repo *repository.SettingsRepository) *SettingsService {
	return &SettingsService{
		repo: repo,
	}
}

// GetAllSettings retrieves all system settings (admin only)
func (s *SettingsService) GetAllSettings(ctx context.Context) ([]models.SystemSetting, error) {
	return s.repo.GetAllSettings(ctx)
}

// GetPublicSettings retrieves only public settings
func (s *SettingsService) GetPublicSettings(ctx context.Context) ([]models.SystemSetting, error) {
	return s.repo.GetPublicSettings(ctx)
}

// GetSettingsByCategory retrieves settings grouped by category
func (s *SettingsService) GetSettingsByCategory(ctx context.Context) ([]models.SettingCategory, error) {
	// Get all categories
	categories, err := s.repo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	var settingCategories []models.SettingCategory
	categoryDisplayNames := map[string]string{
		"general":      "General Settings",
		"email":        "Email Configuration",
		"payment":      "Payment Settings",
		"notification": "Notifications",
		"security":     "Security Settings",
		"business":     "Business Settings",
		"ui":           "UI/UX Settings",
		"seo":          "SEO Settings",
	}

	categoryDescriptions := map[string]string{
		"general":      "Basic application settings and configuration",
		"email":        "Email server and delivery settings",
		"payment":      "Payment gateway and transaction settings",
		"notification": "Notification preferences and settings",
		"security":     "Security policies and authentication settings",
		"business":     "Business-specific settings and rules",
		"ui":           "User interface and experience settings",
		"seo":          "Search engine optimization settings",
	}

	for _, category := range categories {
		settings, err := s.repo.GetSettingsByCategory(ctx, category)
		if err != nil {
			return nil, err
		}

		displayName, exists := categoryDisplayNames[category]
		if !exists {
			displayName = category
		}

		description, exists := categoryDescriptions[category]
		if !exists {
			description = fmt.Sprintf("Settings for %s", category)
		}

		settingCategory := models.SettingCategory{
			Name:        category,
			DisplayName: displayName,
			Description: description,
			Settings:    settings,
		}

		settingCategories = append(settingCategories, settingCategory)
	}

	return settingCategories, nil
}

// GetSetting retrieves a specific setting
func (s *SettingsService) GetSetting(ctx context.Context, category, key string) (*models.SystemSetting, error) {
	return s.repo.GetSetting(ctx, category, key)
}

// UpdateSetting updates a specific setting
func (s *SettingsService) UpdateSetting(ctx context.Context, category, key, value string) error {
	// Validate the setting exists
	_, err := s.repo.GetSetting(ctx, category, key)
	if err != nil {
		return fmt.Errorf("setting not found: %s.%s", category, key)
	}

	return s.repo.UpdateSetting(ctx, category, key, value)
}

// UpdateMultipleSettings updates multiple settings in a transaction
func (s *SettingsService) UpdateMultipleSettings(ctx context.Context, updates []models.SettingUpdateItem) error {
	// Validate all settings exist before updating
	for _, update := range updates {
		_, err := s.repo.GetSetting(ctx, update.Category, update.SettingKey)
		if err != nil {
			return fmt.Errorf("setting not found: %s.%s", update.Category, update.SettingKey)
		}
	}

	return s.repo.UpdateMultipleSettings(ctx, updates)
}

// CreateSetting creates a new setting
func (s *SettingsService) CreateSetting(ctx context.Context, setting *models.SystemSetting) error {
	// Check if setting already exists
	existing, err := s.repo.GetSetting(ctx, setting.Category, setting.SettingKey)
	if err == nil && existing != nil {
		return fmt.Errorf("setting already exists: %s.%s", setting.Category, setting.SettingKey)
	}

	return s.repo.CreateSetting(ctx, setting)
}

// DeleteSetting deletes a specific setting
func (s *SettingsService) DeleteSetting(ctx context.Context, category, key string) error {
	return s.repo.DeleteSetting(ctx, category, key)
}

// GetSettingValue is a helper method to get a setting value with a default
func (s *SettingsService) GetSettingValue(ctx context.Context, category, key, defaultValue string) string {
	setting, err := s.repo.GetSetting(ctx, category, key)
	if err != nil || setting == nil {
		return defaultValue
	}
	return setting.SettingValue
}

// GetSettingValueAsInt gets a setting value as integer with a default
func (s *SettingsService) GetSettingValueAsInt(ctx context.Context, category, key string, defaultValue int) int {
	setting, err := s.repo.GetSetting(ctx, category, key)
	if err != nil || setting == nil {
		return defaultValue
	}

	value, err := setting.GetValueAsInt()
	if err != nil {
		return defaultValue
	}

	return value
}

// GetSettingValueAsBool gets a setting value as boolean with a default
func (s *SettingsService) GetSettingValueAsBool(ctx context.Context, category, key string, defaultValue bool) bool {
	setting, err := s.repo.GetSetting(ctx, category, key)
	if err != nil || setting == nil {
		return defaultValue
	}

	value, err := setting.GetValueAsBool()
	if err != nil {
		return defaultValue
	}

	return value
}
