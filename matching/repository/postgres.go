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
		`delete from likes where user_id = $1`,
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

func (r *PostgresMatchingRepository) GetUserLikes(ctx context.Context, userId string) []string {
	rows, err := r.db.Query(`select liked_id from likes where user_id = $1`, userId)
	if err != nil {
		return nil
	}
	var userLikes []string
	for rows.Next() {
		var likedId string
		if err := rows.Scan(&likedId); err != nil {
			return nil
		}
		userLikes = append(userLikes, likedId)
	}
	return userLikes
}

func (r *PostgresMatchingRepository) IsMutual(ctx context.Context, like *models.Like) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`select exists(
            select 1
            from likes l1
            inner join likes l2 
                on l1.user_id = l2.liked_id 
                and l1.liked_id = l2.user_id
            where l1.user_id = $1 and l1.liked_id = $2
        )`, like.UserId, like.LikedId,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
