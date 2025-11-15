package dblayer

import (
	"database/sql"
	"log"
)

type DBContext struct {
	UserID   string
	GroupIDs []string

	Schema string // prefix to add to a table name
}

func (dbctx *DBContext) IsInGroup(groupID string) bool {
	for _, gid := range dbctx.GroupIDs {
		if gid == groupID {
			return true
		}
	}
	return false
}
func (dbctx *DBContext) AddGroup(groupID string) {
	for _, gid := range dbctx.GroupIDs {
		if gid == groupID {
			return
		}
	}
	dbctx.GroupIDs = append(dbctx.GroupIDs, groupID)
}
func (dbctx *DBContext) IsUser(userID string) bool {
	return dbctx.UserID == userID
}

type DBRepository struct {
	Verbose   bool
	DbContext *DBContext
	factory   *DBEFactory

	/* Can be a connection to mysql, postgresql, sqlite, etc. */
	DbConnection *sql.DB
}

func NewDBRepository(dbContext *DBContext, factory *DBEFactory, dbConnection *sql.DB) *DBRepository {
	return &DBRepository{
		Verbose:      false,
		DbContext:    dbContext,
		factory:      factory,
		DbConnection: dbConnection,
	}
}

func (dbr *DBRepository) GetInstanceByClassName(classname string) *DBEntity {
	return dbr.factory.GetInstanceByClassName(classname)
}
func (dbr *DBRepository) GetInstanceByTableName(tablename string) *DBEntity {
	return dbr.factory.GetInstanceByTableName(tablename)
}

func (dbr *DBRepository) buildTableName(dbe *DBEntity) string {
	if dbr.DbContext != nil && dbr.DbContext.Schema != "" {
		return dbr.DbContext.Schema + "_" + dbe.GetTableName()
	}
	return dbe.GetTableName()
}

func (dbr *DBRepository) Search(dbe *DBEntity, useLike bool, caseSensitive bool, orderBy string) ([]*DBEntity, error) {
	if dbr.Verbose {
		log.Print("DBRepository::Search: dbe=", dbe)
	}

	keys := make([]string, 0)
	values := make([]string, 0)
	for key, value := range dbe.dictionary {
		keys = append(keys, key)
		values = append(values, value.(string))
	}

	return nil, nil
}
