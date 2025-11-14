package db

import (
	"database/sql"
	"fmt"

	"rprj/be/models"
)

// CREATE (deprecated - use CreateGroupWithTransaction)
// func CreateGroup(g models.DBGroup) (string, error) {
// 	if g.ID == "" {
// 		newID, _ := uuid16HexGo()
// 		log.Print("newID=", newID)
// 		g.ID = newID
// 	}
// 	_, err := DB.Exec(
// 		"INSERT INTO "+tablePrefix+"groups (id, name, description) VALUES (?, ?, ?)",
// 		g.ID, g.Name, g.Description,
// 	)
// 	return g.ID, err
// }

// READ
func GetGroupByID(id string) (*models.DBGroup, error) {
	row := DB.QueryRow(
		"SELECT id, name, description FROM "+tablePrefix+"groups WHERE id = ?",
		id,
	)

	var g models.DBGroup
	err := row.Scan(&g.ID, &g.Name, &g.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Nessun gruppo trovato
		}
		return nil, err
	}
	return &g, nil
}

// UPDATE (deprecated - use UpdateGroupWithTransaction)
// func UpdateGroup(g models.DBGroup) error {
// 	_, err := DB.Exec(
// 		"UPDATE "+tablePrefix+"groups SET name = ?, description = ? WHERE id = ?",
// 		g.Name, g.Description, g.ID,
// 	)
// 	return err
// }

// DELETE
func DeleteGroup(id string) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete group-user associations
	_, err = tx.Exec("DELETE FROM "+tablePrefix+"users_groups WHERE group_id=?", id)
	if err != nil {
		return err
	}

	// Delete group
	_, err = tx.Exec("DELETE FROM "+tablePrefix+"groups WHERE id=?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// SEARCH
func SearchGroupsBy(search string, orderBy string) ([]models.DBGroup, error) {
	query := "SELECT id, name, description FROM " + tablePrefix + "groups"
	if search != "" {
		query += " WHERE name LIKE ? OR description LIKE ?"
	}
	if orderBy != "" {
		query += " ORDER BY " + orderBy
	}

	likePattern := "%" + search + "%"

	var rows *sql.Rows
	var err error
	if search != "" {
		rows, err = DB.Query(query, likePattern, likePattern)
	} else {
		rows, err = DB.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.DBGroup
	for rows.Next() {
		var g models.DBGroup
		err := rows.Scan(&g.ID, &g.Name, &g.Description)
		if err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

// CreateGroup creates a group and user associations in a single transaction
func CreateGroup(g models.DBGroup, userIDs []string) (*models.DBGroup, error) {
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Check that group with same name does not already exist
	var existingGroupID string
	err = tx.QueryRow("SELECT id FROM "+tablePrefix+"groups WHERE name = ?", g.Name).Scan(&existingGroupID)
	if err != sql.ErrNoRows {
		if err == nil {
			return nil, fmt.Errorf("group with name '%s' already exists", g.Name)
		}
		return nil, err
	}

	// Generate ID
	groupID, _ := uuid16HexGo()
	g.ID = groupID

	// Create group
	_, err = tx.Exec(
		"INSERT INTO "+tablePrefix+"groups (id, name, description) VALUES (?, ?, ?)",
		g.ID, g.Name, g.Description,
	)
	if err != nil {
		return nil, err
	}

	// Add users to group
	for _, userID := range userIDs {
		_, err = tx.Exec(
			"INSERT INTO "+tablePrefix+"users_groups (user_id, group_id) VALUES (?, ?)",
			userID, groupID,
		)
		if err != nil {
			return nil, err
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &g, nil
}

// UpdateGroup updates group and user associations in a single transaction
func UpdateGroup(g models.DBGroup, userIDs []string) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update group
	_, err = tx.Exec(
		"UPDATE "+tablePrefix+"groups SET name = ?, description = ? WHERE id = ?",
		g.Name, g.Description, g.ID,
	)
	if err != nil {
		return err
	}

	// Delete all existing user associations
	_, err = tx.Exec("DELETE FROM "+tablePrefix+"users_groups WHERE group_id=?", g.ID)
	if err != nil {
		return err
	}

	// Recreate user associations
	for _, userID := range userIDs {
		_, err = tx.Exec(
			"INSERT INTO "+tablePrefix+"users_groups (user_id, group_id) VALUES (?, ?)",
			userID, g.ID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
