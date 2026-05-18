package repository

import (
	"context"
	"database/sql"
	"fmt"
	profile2 "profile/internal"
	"profile/internal/models"

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

func (r *PostgresProfileRepository) GetAll(ctx context.Context) []models.Profile {
	//TODO CACHE!!!
	var profiles []models.Profile
	rows, err := r.db.Query(`select id, user_id, name, gender, description, date_created, photo_path from profiles`)
	if err != nil {
		return profiles
	}
	for rows.Next() {
		profile := &models.Profile{}
		err := rows.Scan(&profile.Id, &profile.UserId, &profile.Name, &profile.Gender, &profile.Description, &profile.DateCreated, &profile.PhotoPath)
		if err != nil {
			return profiles
		}
		profiles = append(profiles, *profile)
	}
	return profiles
}

// TODO all updates
func (r *PostgresProfileRepository) UpdateById(ctx context.Context, id string, newProfile *models.Profile) error {
	result, err := r.db.Exec(
		`update profiles set name = $1, gender = $2, description = $3 where user_id = $4`, newProfile.Name, newProfile.Gender, newProfile.Description, id,
	)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return profile2.ErrUserNotFound
	}

	return nil
}

func (r *PostgresProfileRepository) DeleteById(ctx context.Context, id string) error {
	result, err := r.db.Exec(
		`delete from profiles where user_id = $1`, id,
	)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return profile2.ErrUserNotFound
	}

	return nil
}
