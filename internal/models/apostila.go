package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type Apostila struct {
	Id         uuid.UUID `json:"id"`
	UserID     int64     `json:"user_id"`
	EditedHTML string    `json:"edited_raw_html"`
	CreatedAt  string    `json:"created_at"`
	EditedAt   string    `json:"edited_at"`
}

type EditedApostilaHTML struct {
	HTML string `json:"file"`
}

type ApostilaModel struct {
	DB *sql.DB
}

func (m *ApostilaModel) Insert(ctx context.Context, id uuid.UUID, userID int64) (*Apostila, error) {
	query := `
		INSERT INTO apostilas (id, user_id, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, user_id, created_at
	`

	var apostila Apostila
	err := m.DB.QueryRowContext(ctx, query, id, userID).Scan(
		&apostila.Id,
		&apostila.UserID,
		&apostila.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &apostila, nil
}

func (m *ApostilaModel) UpdateEditedHTMLByID(ctx context.Context, id uuid.UUID, editedHTML string, userID int64) error {
	query := `
		UPDATE apostilas
		SET edited_html = $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
	`

	result, err := m.DB.ExecContext(ctx, query, editedHTML, id, userID)
	if err != nil {
		log.Println("Error executing update:", err)
		return fmt.Errorf("UpdateEditedHTMLByID: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return fmt.Errorf("UpdateEditedHTMLByID: failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		log.Println("No rows were updated")
		return sql.ErrNoRows
	}

	log.Println("Rows affected:", rowsAffected)

	return nil
}

func (m *ApostilaModel) GetEditedHTMLByID(ctx context.Context, id uuid.UUID, userId int64) (*EditedApostilaHTML, error) {
	query := `
	SELECT edited_html 
	FROM apostilas 
	WHERE id = $1 AND user_id = $2
	`

	var editedApostilaHTML EditedApostilaHTML
	err := m.DB.QueryRowContext(ctx, query, id, userId).Scan(&editedApostilaHTML.HTML)
	if err != nil {
		editedApostilaHTML.HTML = ""
	}

	return &editedApostilaHTML, nil
}

func (m *ApostilaModel) GetAllByUserID(ctx context.Context, userID int64) ([]*Apostila, error) {
	query := `
		SELECT id, user_id, edited_html, created_at, COALESCE(updated_at, created_at) as edited_at
		FROM apostilas
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apostilas []*Apostila
	for rows.Next() {
		var apostila Apostila
		err := rows.Scan(
			&apostila.Id,
			&apostila.UserID,
			&apostila.EditedHTML,
			&apostila.CreatedAt,
			&apostila.EditedAt,
		)
		if err != nil {
			return nil, err
		}
		apostilas = append(apostilas, &apostila)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return apostilas, nil
}

func (m *ApostilaModel) GetByID(ctx context.Context, id uuid.UUID) (*Apostila, error) {
	query := `
		SELECT id, user_id, edited_html, created_at, COALESCE(updated_at, created_at) as edited_at
		FROM apostilas
		WHERE id = $1
	`

	var apostila Apostila
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&apostila.Id,
		&apostila.UserID,
		&apostila.EditedHTML,
		&apostila.CreatedAt,
		&apostila.EditedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("apostila not found")
	}

	if err != nil {
		return nil, err
	}

	return &apostila, nil
}

func (m *ApostilaModel) Delete(ctx context.Context, id uuid.UUID, userID int64) error {
	query := `
		DELETE FROM apostilas
		WHERE id = $1 AND user_id = $2
	`

	result, err := m.DB.ExecContext(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete apostila: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("apostila not found or unauthorized")
	}

	return nil
}
