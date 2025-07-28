package models

import (
	"encoding/json"
	"strconv"
	"time"
)

type DataType string

const (
	DataTypeString  DataType = "string"
	DataTypeNumber  DataType = "number"
	DataTypeBoolean DataType = "boolean"
	DataTypeJSON    DataType = "json"
)

type SystemSetting struct {
	ID           int       `json:"id" db:"id"`
	Category     string    `json:"category" db:"category"`
	SettingKey   string    `json:"setting_key" db:"setting_key"`
	SettingValue string    `json:"setting_value" db:"setting_value"`
	DataType     DataType  `json:"data_type" db:"data_type"`
	Description  string    `json:"description" db:"description"`
	IsEncrypted  bool      `json:"is_encrypted" db:"is_encrypted"`
	IsPublic     bool      `json:"is_public" db:"is_public"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// GetValueAsString returns the setting value as string
func (s *SystemSetting) GetValueAsString() string {
	return s.SettingValue
}

// GetValueAsInt returns the setting value as integer
func (s *SystemSetting) GetValueAsInt() (int, error) {
	return strconv.Atoi(s.SettingValue)
}

// GetValueAsFloat returns the setting value as float64
func (s *SystemSetting) GetValueAsFloat() (float64, error) {
	return strconv.ParseFloat(s.SettingValue, 64)
}

// GetValueAsBool returns the setting value as boolean
func (s *SystemSetting) GetValueAsBool() (bool, error) {
	return strconv.ParseBool(s.SettingValue)
}

// GetValueAsJSON unmarshals the setting value into the provided interface
func (s *SystemSetting) GetValueAsJSON(v interface{}) error {
	return json.Unmarshal([]byte(s.SettingValue), v)
}

// SetValue sets the setting value from various types
func (s *SystemSetting) SetValue(value interface{}) error {
	switch v := value.(type) {
	case string:
		s.SettingValue = v
		s.DataType = DataTypeString
	case int:
		s.SettingValue = strconv.Itoa(v)
		s.DataType = DataTypeNumber
	case float64:
		s.SettingValue = strconv.FormatFloat(v, 'f', -1, 64)
		s.DataType = DataTypeNumber
	case bool:
		s.SettingValue = strconv.FormatBool(v)
		s.DataType = DataTypeBoolean
	default:
		// For complex types, marshal as JSON
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return err
		}
		s.SettingValue = string(jsonBytes)
		s.DataType = DataTypeJSON
	}
	return nil
}

// SettingCategory represents a group of related settings
type SettingCategory struct {
	Name        string          `json:"name"`
	DisplayName string          `json:"display_name"`
	Description string          `json:"description"`
	Settings    []SystemSetting `json:"settings"`
}

// SettingsUpdateRequest represents a request to update multiple settings
type SettingsUpdateRequest struct {
	Settings []SettingUpdateItem `json:"settings"`
}

type SettingUpdateItem struct {
	Category     string      `json:"category"`
	SettingKey   string      `json:"setting_key"`
	SettingValue interface{} `json:"setting_value"`
}
