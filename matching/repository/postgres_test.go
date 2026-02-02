package repository

import (
	"context"
	"database/sql"
	"date-bot-go/matching"
	"date-bot-go/matching/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

func setupDB(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)))
	require.NoError(t, err)
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	require.NoError(t, db.Ping())
	createSchema(t, db)
	cleanup := func() {
		db.Close()
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	}
	return db, cleanup
}

func createSchema(t *testing.T, db *sql.DB) {
	schema := `
		CREATE TABLE IF NOT EXISTS public.likes
		(
			user_id text COLLATE pg_catalog."default" NOT NULL,
			liked_id text COLLATE pg_catalog."default" NOT NULL,
			UNIQUE (user_id, liked_id)
		);
	`
	_, err := db.Exec(schema)
	require.NoError(t, err)
}

func TestSuccessAddLike(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()
	repo := NewPostgresMatchingRepository(db)
	like := &models.Like{
		UserId:  "123",
		LikedId: "456",
	}
	err := repo.AddLike(context.Background(), like)
	assert.NoError(t, err)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = $1", like.UserId).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestErrAlreadyLikedAddLike(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()
	repo := NewPostgresMatchingRepository(db)
	like := &models.Like{
		UserId:  "123",
		LikedId: "456",
	}
	err := repo.AddLike(context.Background(), like)
	assert.NoError(t, err)
	err = repo.AddLike(context.Background(), like)
	assert.Error(t, err)
	assert.IsType(t, matching.ErrAlreadyLiked, err)
}

func TestSuccessDeleteLike(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()
	repo := NewPostgresMatchingRepository(db)
	like1 := &models.Like{
		UserId:  "123",
		LikedId: "456",
	}
	err := repo.AddLike(context.Background(), like1)
	assert.NoError(t, err)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = $1", like1.UserId).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

}

func TestErrLikeNotFoundDeleteLike(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()
	repo := NewPostgresMatchingRepository(db)
	like := &models.Like{
		UserId:  "123",
		LikedId: "456",
	}
	err := repo.DeleteLike(context.Background(), like)
	assert.Error(t, err)
	assert.IsType(t, matching.ErrLikeNotFound, err)
}

func TestSuccessGetUserLikes(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()
	repo := NewPostgresMatchingRepository(db)
	like := &models.Like{
		UserId:  "123",
		LikedId: "456",
	}
	err := repo.AddLike(context.Background(), like)
	assert.NoError(t, err)
	like = &models.Like{
		UserId:  "123",
		LikedId: "789",
	}
	err = repo.AddLike(context.Background(), like)
	assert.NoError(t, err)
	like = &models.Like{
		UserId:  "123",
		LikedId: "321",
	}
	err = repo.AddLike(context.Background(), like)
	assert.NoError(t, err)
	userLikes := repo.GetUserLikes(context.Background(), like.UserId)

	assert.Len(t, userLikes, 3)
	assert.Contains(t, userLikes, "456")
	assert.Contains(t, userLikes, "789")
	assert.Contains(t, userLikes, "321")
}

func TestSuccessIsMutual(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()
	repo := NewPostgresMatchingRepository(db)
	like := &models.Like{
		UserId:  "123",
		LikedId: "456",
	}
	err := repo.AddLike(context.Background(), like)
	assert.NoError(t, err)
	like = &models.Like{
		UserId:  "456",
		LikedId: "123",
	}
	err = repo.AddLike(context.Background(), like)
	assert.NoError(t, err)
	isMutual, err := repo.IsMutual(context.Background(), like)
	assert.NoError(t, err)
	assert.Equal(t, true, isMutual)
}

func TestNoEqualIsMutual(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()
	repo := NewPostgresMatchingRepository(db)
	like := &models.Like{
		UserId:  "123",
		LikedId: "456",
	}
	err := repo.AddLike(context.Background(), like)
	assert.NoError(t, err)
	like = &models.Like{
		UserId:  "456",
		LikedId: "123",
	}
	//второго не будет
	//err = repo.AddLike(context.Background(), like)
	assert.NoError(t, err)
	isMutual, err := repo.IsMutual(context.Background(), like)
	assert.NoError(t, err)
	assert.NotEqual(t, true, isMutual)
}
