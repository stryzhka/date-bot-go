package repository

import (
	"context"
	"database/sql"
	profile2 "date-bot-go/profile"
	"date-bot-go/profile/models"
	"fmt"
	"github.com/lib/pq"
)

type PostgresProfileRepository struct {
	db *sql.DB
}

func NewPostgresProfileRepository(db *sql.DB) *PostgresProfileRepository {
	return &PostgresProfileRepository{db: db}
}

func (r *PostgresProfileRepository) Create(ctx context.Context, profile *models.Profile) error {
	_, err := r.db.Exec(
		`
				insert into profiles (id, user_id, name, gender, description, date_created, photo_path) values 
				($1, $2, $3, $4, $5, $6, $7)
				`,
		profile.Id,
		profile.UserId,
		profile.Name,
		profile.Gender,
		profile.Description,
		profile.DateCreated,
		profile.PhotoPath,
	)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			if e.Code == "23505" {
				return profile2.ErrUserAlreadyExists
			}
		}
		return fmt.Errorf("failed to create profile: %w", err)
	}
	return nil
}

func (r *PostgresProfileRepository) Get(ctx context.Context, id string) *models.Profile {
	profile := &models.Profile{}
	err := r.db.QueryRow(`
		select id, user_id, name, gender, description, date_created, photo_path from profiles where user_id = $1
	`, id).Scan(
		&profile.Id,
		&profile.UserId,
		&profile.Name,
		&profile.Gender,
		&profile.Description,
		&profile.DateCreated,
		&profile.PhotoPath,
	)
	if err != nil {
		profile = nil
	}
	return profile
}
