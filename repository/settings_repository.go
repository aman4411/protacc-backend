package repository

import (
	"context"
	"fmt"

	"github.com/aman4411/protacc-backend/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SettingsRepository struct {
	db *pgxpool.Pool
}

func NewSettingsRepository(db *pgxpool.Pool) *SettingsRepository {
	return &SettingsRepository{db: db}
}

// GetAllSettings retrieves all system settings
func (r *SettingsRepository) GetAllSettings(ctx context.Context) ([]models.SystemSetting, error) {
	query := `
		SELECT id, category, setting_key, setting_value, data_type, description, is_encrypted, is_public, created_at, updated_at
		FROM system_settings
		ORDER BY category, setting_key`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []models.SystemSetting
	for rows.Next() {
		var setting models.SystemSetting
		err := rows.Scan(
			&setting.ID,
			&setting.Category,
			&setting.SettingKey,
			&setting.SettingValue,
			&setting.DataType,
			&setting.Description,
			&setting.IsEncrypted,
			&setting.IsPublic,
			&setting.CreatedAt,
			&setting.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		settings = append(settings, setting)
	}

	return settings, nil
}

// GetSettingsByCategory retrieves all settings for a specific category
func (r *SettingsRepository) GetSettingsByCategory(ctx context.Context, category string) ([]models.SystemSetting, error) {
	query := `
		SELECT id, category, setting_key, setting_value, data_type, description, is_encrypted, is_public, created_at, updated_at
		FROM system_settings
		WHERE category = $1
		ORDER BY setting_key`

	rows, err := r.db.Query(ctx, query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []models.SystemSetting
	for rows.Next() {
		var setting models.SystemSetting
		err := rows.Scan(
			&setting.ID,
			&setting.Category,
			&setting.SettingKey,
			&setting.SettingValue,
			&setting.DataType,
			&setting.Description,
			&setting.IsEncrypted,
			&setting.IsPublic,
			&setting.CreatedAt,
			&setting.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		settings = append(settings, setting)
	}

	return settings, nil
}

// GetPublicSettings retrieves all public settings
func (r *SettingsRepository) GetPublicSettings(ctx context.Context) ([]models.SystemSetting, error) {
	query := `
		SELECT id, category, setting_key, setting_value, data_type, description, is_encrypted, is_public, created_at, updated_at
		FROM system_settings
		WHERE is_public = true
		ORDER BY category, setting_key`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []models.SystemSetting
	for rows.Next() {
		var setting models.SystemSetting
		err := rows.Scan(
			&setting.ID,
			&setting.Category,
			&setting.SettingKey,
			&setting.SettingValue,
			&setting.DataType,
			&setting.Description,
			&setting.IsEncrypted,
			&setting.IsPublic,
			&setting.CreatedAt,
			&setting.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		settings = append(settings, setting)
	}

	return settings, nil
}

// GetSetting retrieves a specific setting by category and key
func (r *SettingsRepository) GetSetting(ctx context.Context, category, key string) (*models.SystemSetting, error) {
	query := `
		SELECT id, category, setting_key, setting_value, data_type, description, is_encrypted, is_public, created_at, updated_at
		FROM system_settings
		WHERE category = $1 AND setting_key = $2`

	var setting models.SystemSetting
	err := r.db.QueryRow(ctx, query, category, key).Scan(
		&setting.ID,
		&setting.Category,
		&setting.SettingKey,
		&setting.SettingValue,
		&setting.DataType,
		&setting.Description,
		&setting.IsEncrypted,
		&setting.IsPublic,
		&setting.CreatedAt,
		&setting.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &setting, nil
}

// UpdateSetting updates a specific setting
func (r *SettingsRepository) UpdateSetting(ctx context.Context, category, key, value string) error {
	query := `
		UPDATE system_settings 
		SET setting_value = $1, updated_at = NOW()
		WHERE category = $2 AND setting_key = $3`

	result, err := r.db.Exec(ctx, query, value, category, key)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("setting not found: %s.%s", category, key)
	}

	return nil
}

// UpdateMultipleSettings updates multiple settings in a transaction
func (r *SettingsRepository) UpdateMultipleSettings(ctx context.Context, updates []models.SettingUpdateItem) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, update := range updates {
		// Convert the value to string based on its type
		var valueStr string
		switch v := update.SettingValue.(type) {
		case string:
			valueStr = v
		case bool:
			if v {
				valueStr = "true"
			} else {
				valueStr = "false"
			}
		case float64:
			valueStr = fmt.Sprintf("%v", v)
		default:
			valueStr = fmt.Sprintf("%v", v)
		}

		query := `
			UPDATE system_settings 
			SET setting_value = $1, updated_at = NOW()
			WHERE category = $2 AND setting_key = $3`

		_, err = tx.Exec(ctx, query, valueStr, update.Category, update.SettingKey)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// CreateSetting creates a new setting
func (r *SettingsRepository) CreateSetting(ctx context.Context, setting *models.SystemSetting) error {
	query := `
		INSERT INTO system_settings (category, setting_key, setting_value, data_type, description, is_encrypted, is_public)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		setting.Category,
		setting.SettingKey,
		setting.SettingValue,
		setting.DataType,
		setting.Description,
		setting.IsEncrypted,
		setting.IsPublic,
	).Scan(&setting.ID, &setting.CreatedAt, &setting.UpdatedAt)

	return err
}

// DeleteSetting deletes a specific setting
func (r *SettingsRepository) DeleteSetting(ctx context.Context, category, key string) error {
	query := `DELETE FROM system_settings WHERE category = $1 AND setting_key = $2`

	result, err := r.db.Exec(ctx, query, category, key)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("setting not found: %s.%s", category, key)
	}

	return nil
}

// GetCategories retrieves all unique categories
func (r *SettingsRepository) GetCategories(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT category FROM system_settings ORDER BY category`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}
