package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"rprj/be/db"
	"rprj/be/models"

	"github.com/gorilla/mux"
)

// GET /groups
func GetAllGroupsHandler(w http.ResponseWriter, r *http.Request) {
	searchBy := r.URL.Query().Get("search")
	orderBy := r.URL.Query().Get("order_by")
	groups, err := db.SearchGroupsBy(searchBy, orderBy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(groups)
}

// GET /groups/{id}
func GetGroupHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	group, err := db.GetGroupByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if group == nil {
		http.NotFound(w, r)
		return
	}

	// Get group users
	groupUsers, err := db.GetUserGroupsByGroupID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build response with user IDs
	userIDs := make([]string, len(groupUsers))
	for i, gu := range groupUsers {
		userIDs[i] = gu.UserID
	}

	response := map[string]interface{}{
		"id":          group.ID,
		"name":        group.Name,
		"description": group.Description,
		"user_ids":    userIDs,
	}

	json.NewEncoder(w).Encode(response)
}

// POST /groups
func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		UserIDs     []string `json:"user_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request format"})
		return
	}

	// Validate required fields
	if req.Name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Group name is required"})
		return
	}

	g := models.DBGroup{
		Name:        req.Name,
		Description: req.Description,
	}

	// Create group with transaction
	createdGroup, err := db.CreateGroup(g, req.UserIDs)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		// Check if it's a duplicate name error
		if strings.Contains(err.Error(), "already exists") {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create group: " + err.Error()})
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdGroup)
}

// PUT /groups/{id}
func UpdateGroupHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing group ID"})
		return
	}

	var req struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		UserIDs     []string `json:"user_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request format"})
		return
	}

	// Validate required fields
	if req.Name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Group name is required"})
		return
	}

	g := models.DBGroup{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}

	// Update group with transaction
	if err := db.UpdateGroup(g, req.UserIDs); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update group: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /groups/{id}
func DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing group ID"})
		return
	}

	// Groups with negative ID cannot be deleted (system groups)
	if strings.HasPrefix(id, "-") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cannot delete system groups"})
		return
	}

	if err := db.DeleteGroup(id); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete group: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
