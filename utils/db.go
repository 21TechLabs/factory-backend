package utils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// JSONMap is a custom type to handle JSON data in map[string]interface{} format
type JSONMap[T any] map[string]T

// Value implements the driver.Valuer interface, converting the map to a JSON string.
func (m *JSONMap[T]) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// Scan implements the sql.Scanner interface, converting the JSON string from the DB to a map.
func (m *JSONMap[T]) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return json.Unmarshal(bytes, &m)
}

// GormDataType sets the GORM data type for the custom JSONMap
func (*JSONMap[T]) GormDataType() string {
	return "json"
}

// GormDBDataType ensures compatibility with different database dialects.
func (*JSONMap[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

// StringSlice is a custom type to handle []string data.
type StringSlice []string

// Value converts the slice to a JSON string for database storage.
func (s *StringSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan converts the JSON string from the database back to a slice.
func (s *StringSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB value")
	}
	return json.Unmarshal(bytes, s)
}

// GormDataType sets the GORM data type.
func (*StringSlice) GormDataType() string {
	return "json"
}
