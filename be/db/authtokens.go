package db

import (
	"log"
)

func SaveToken(userID string, tokenString string, expiry int64) error {
	_, err := DB.Exec(
		"INSERT INTO "+tablePrefix+"oauth_tokens (user_id, token_id, expires_at) VALUES (?, ?, ?)",
		userID, tokenString, expiry,
	)
	if err != nil {
		log.Println("Errore salvataggio token:", err)
	}
	return err
}

func IsTokenValid(tokenString string, userID string) bool {
	var count int
	err := DB.QueryRow(
		"SELECT COUNT(*) FROM "+tablePrefix+"oauth_tokens WHERE token_id = ? AND user_id = ?",
		tokenString, userID,
	).Scan(&count)
	if err != nil {
		log.Println("Errore verifica token:", err)
		return false
	}
	return count > 0
}

func DeleteToken(tokenString string) error {
	_, err := DB.Exec(
		"DELETE FROM "+tablePrefix+"oauth_tokens WHERE token_id = ?",
		tokenString,
	)
	if err != nil {
		log.Println("Errore cancellazione token:", err)
	}
	return err
}
