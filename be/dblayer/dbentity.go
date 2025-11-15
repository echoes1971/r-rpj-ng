package dblayer

import (
	"sort"
	"strings"
)

type ForeignKey struct {
	Column    string
	RefTable  string
	RefColumn string
}

type Column struct {
	Name        string
	Type        string
	Constraints []string
}

type DBEntity struct {
	typename    string
	tablename   string
	columns     map[string]Column
	keys        []string
	foreignKeys []ForeignKey
	dictionary  map[string]any
}

func NewDBEntity(typename string, tablename string, columns []Column, keys []string, foreignKeys []ForeignKey, dictionary map[string]any) *DBEntity {
	columnsMap := make(map[string]Column)
	for _, col := range columns {
		columnsMap[col.Name] = col
	}
	return &DBEntity{
		typename:    typename,
		tablename:   tablename,
		columns:     columnsMap,
		keys:        keys,
		foreignKeys: foreignKeys,
		dictionary:  dictionary,
	}
}

/* Override */
func (dbEntity *DBEntity) NewInstance() *DBEntity {
	columns := make([]Column, 0, len(dbEntity.columns))
	for _, col := range dbEntity.columns {
		columns = append(columns, col)
	}
	return NewDBEntity(dbEntity.typename, dbEntity.tablename, columns, dbEntity.keys, dbEntity.foreignKeys, make(map[string]any))
}

func (dbEntity *DBEntity) GetColumnType(columnName string) string {
	if col, exists := dbEntity.columns[columnName]; exists {
		return col.Type
	}
	return ""
}
func (dbEntity *DBEntity) GetTypeName() string {
	return dbEntity.typename
}
func (dbEntity *DBEntity) GetTableName() string {
	return dbEntity.tablename
}
func (dbEntity *DBEntity) GetKeys() []string {
	return dbEntity.keys
}
func (dbEntity *DBEntity) GetForeignKeys() []ForeignKey {
	return dbEntity.foreignKeys
}
func (dbEntity *DBEntity) GetOrderBy() []string {
	return dbEntity.GetKeys()
}
func (dbEntity *DBEntity) GetOrderByString() string {
	return strings.Join(dbEntity.GetOrderBy(), ", ")
}
func (dbEntity *DBEntity) GetForeignKeysForTable(tableName string) []ForeignKey {
	var foreignKeysForTable []ForeignKey
	for _, fk := range dbEntity.foreignKeys {
		if fk.RefTable == tableName {
			foreignKeysForTable = append(foreignKeysForTable, fk)
		}
	}
	return foreignKeysForTable
}
func (dbEntity *DBEntity) GetForeignKeyDefinition(columnName string) *ForeignKey {
	for _, fk := range dbEntity.foreignKeys {
		if fk.Column == columnName {
			return &fk
		}
	}
	return nil
}

// TODO? Manage different types of values (int, date, etc.)
func (dbEntity *DBEntity) SetValue(columnName string, value string) {
	// if _, exists := dbEntity.dictionary[columnName]; exists {
	dbEntity.dictionary[columnName] = value
	// }
}
func (dbEntity *DBEntity) GetValue(columnName string) string {
	if val, exists := dbEntity.dictionary[columnName]; exists {
		return val.(string)
	}
	return ""
}
func (dbEntity *DBEntity) HasValue(columnName string) bool {
	_, exists := dbEntity.dictionary[columnName]
	return exists
}
func (dbEntity *DBEntity) ReadFKFrom(dbe *DBEntity) {
	fks := dbEntity.GetForeignKeysForTable(dbe.GetTableName())
	for _, fk := range fks {
		value := dbe.GetValue(fk.RefColumn)
		dbEntity.SetValue(fk.Column, value)
	}
}
func (dbEntity *DBEntity) WriteToFK(dbe *DBEntity) {
	fks := dbEntity.GetForeignKeysForTable(dbe.GetTableName())
	for _, fk := range fks {
		value := dbEntity.GetValue(fk.Column)
		dbe.SetValue(fk.RefColumn, value)
	}
}
func (dbEntity *DBEntity) IsPrimaryKey(columnName string) bool {
	for _, key := range dbEntity.keys {
		if key == columnName {
			return true
		}
	}
	return false
}
func (dbEntity *DBEntity) IsForeignKey(columnName string) bool {
	for _, fk := range dbEntity.foreignKeys {
		if fk.Column == columnName {
			return true
		}
	}
	return false
}

/*
Returns the dictionary keys which means all values set in the entity
*/
func (dbEntity *DBEntity) GetDictionaryKeys() []string {
	keys := make([]string, 0, len(dbEntity.dictionary))
	for key := range dbEntity.dictionary {
		keys = append(keys, key)
	}
	// Sort the keys alphabetically
	sort.Strings(keys)
	return keys
}
func (dbEntity *DBEntity) GetDictionaryValues() []string {
	keys := dbEntity.GetDictionaryKeys() // If I use this, the sorting of the keys may be unnecessary
	values := make([]string, 0, len(keys))
	for _, key := range keys {
		values = append(values, dbEntity.dictionary[key].(string))
	}
	return values
}

/*
Returns a dictionary of the keys set in the entity
*/
func (dbEntity *DBEntity) GetKeySetDictionary() map[string]string {
	result := make(map[string]string)
	for _, key := range dbEntity.keys {
		if val, exists := dbEntity.dictionary[key]; exists {
			result[key] = val.(string)
		}
	}
	return result
}

/*
Remove keys from dictionary
*/
func (dbEntity *DBEntity) RemoveKeysFromDictionary() {
	for _, key := range dbEntity.keys {
		delete(dbEntity.dictionary, key)
	}
}

/*
Returns true if all primary keys have not been set in the dictionary
*/
func (dbEntity *DBEntity) IsNew() bool {
	for _, key := range dbEntity.keys {
		if _, exists := dbEntity.dictionary[key]; exists {
			return false
		}
	}
	return true
}

func (dbEntity *DBEntity) beforeInsert(dbRepository *DBRepository) error {
	// Implement any logic needed before inserting the entity into the database
	return nil
}

func (dbEntity *DBEntity) afterInsert(dbRepository *DBRepository) error {
	// Implement any logic needed after inserting the entity into the database
	return nil
}

func (dbEntity *DBEntity) beforeUpdate(dbRepository *DBRepository) error {
	// Implement any logic needed before updating the entity in the database
	return nil
}

func (dbEntity *DBEntity) afterUpdate(dbRepository *DBRepository) error {
	// Implement any logic needed after updating the entity in the database
	return nil
}

func (dbEntity *DBEntity) beforeDelete(dbRepository *DBRepository) error {
	// Implement any logic needed before deleting the entity from the database
	return nil
}

func (dbEntity *DBEntity) afterDelete(dbRepository *DBRepository) error {
	// Implement any logic needed after deleting the entity from the database
	return nil
}
