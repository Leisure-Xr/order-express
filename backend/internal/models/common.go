package models

// I18nString is the multilingual text field used across the app.
// Stored as JSON text in SQLite: {"zh":"...","en":"..."}
type I18nString struct {
	ZH string `json:"zh"`
	EN string `json:"en"`
}
