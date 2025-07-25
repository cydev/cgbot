// Code generated by ent, DO NOT EDIT.

package telegramchannel

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the telegramchannel type in the database.
	Label = "telegram_channel"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldAccessHash holds the string denoting the access_hash field in the database.
	FieldAccessHash = "access_hash"
	// FieldTitle holds the string denoting the title field in the database.
	FieldTitle = "title"
	// FieldSaveRecords holds the string denoting the save_records field in the database.
	FieldSaveRecords = "save_records"
	// FieldSaveFavoriteRecords holds the string denoting the save_favorite_records field in the database.
	FieldSaveFavoriteRecords = "save_favorite_records"
	// FieldActive holds the string denoting the active field in the database.
	FieldActive = "active"
	// Table holds the table name of the telegramchannel in the database.
	Table = "telegram_channels"
)

// Columns holds all SQL columns for telegramchannel fields.
var Columns = []string{
	FieldID,
	FieldAccessHash,
	FieldTitle,
	FieldSaveRecords,
	FieldSaveFavoriteRecords,
	FieldActive,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the TelegramChannel queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByAccessHash orders the results by the access_hash field.
func ByAccessHash(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAccessHash, opts...).ToFunc()
}

// ByTitle orders the results by the title field.
func ByTitle(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTitle, opts...).ToFunc()
}

// BySaveRecords orders the results by the save_records field.
func BySaveRecords(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSaveRecords, opts...).ToFunc()
}

// BySaveFavoriteRecords orders the results by the save_favorite_records field.
func BySaveFavoriteRecords(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSaveFavoriteRecords, opts...).ToFunc()
}

// ByActive orders the results by the active field.
func ByActive(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldActive, opts...).ToFunc()
}
