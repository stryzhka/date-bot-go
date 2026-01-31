package repository

import (
	"context"
	"database/sql"
	matching2 "date-bot-go/matching"
	"date-bot-go/matching/models"
	"fmt"
	"github.com/lib/pq"
)

type PostgresMatchingRepository struct {
	db *sql.DB
}

func NewPostgresMatchingRepository(db *sql.DB) *PostgresMatchingRepository {
	return &PostgresMatchingRepository{db: db}
}

func (r *PostgresMatchingRepository) AddLike(ctx context.Context, like *models.Like) error {
	_, err := r.db.Exec(
		`
				insert into likes (user_id, liked_id) values 
				($1, $2)
				`,
		like.UserId,
		like.LikedId,
	)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			if e.Code == "23505" {
				return matching2.ErrAlreadyLiked
			}
		}
		return fmt.Errorf("failed to add like: %w", err)
	}
	return nil
}

func (r *PostgresMatchingRepository) DeleteLike(ctx context.Context, like *models.Like) error {
	//оди
	result, err := r.db.Exec(
		`DELETE FROM likes WHERE user_id = $1`,
		like.UserId, like.LikedId,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return matching2.ErrLikeNotFound
	}

	return nil
}

func (r *PostgresMatchingRepository) IsMutual(ctx context.Context, like *models.Like) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(
            SELECT 1
            FROM likes l1
            INNER JOIN likes l2 
                ON l1.user_id = l2.liked_id 
                AND l1.liked_id = l2.user_id
            WHERE l1.user_id = $1 AND l1.liked_id = $2
        )`, like.UserId, like.LikedId,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
