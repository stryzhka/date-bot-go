package repository

import (
	"context"
	"database/sql"
	"date-bot-go/profile/models"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
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
		CREATE TABLE IF NOT EXISTS public.profiles
		(
			id uuid NOT NULL,
			user_id text COLLATE pg_catalog."default" NOT NULL,
			name text COLLATE pg_catalog."default" NOT NULL,
			gender text COLLATE pg_catalog."default" NOT NULL,
			description text COLLATE pg_catalog."default",
			date_created date NOT NULL,
			photo_path text COLLATE pg_catalog."default",
			CONSTRAINT id PRIMARY KEY (id),
			CONSTRAINT unique_profile UNIQUE (user_id)
		);
	`
	_, err := db.Exec(schema)
	require.NoError(t, err)
}

func TestPostgresSuccessCreate(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()
	repo := NewPostgresProfileRepository(db)
	profile := &models.Profile{
		Id:          uuid.New(),
		UserId:      "123",
		Name:        "test",
		Gender:      "f",
		Description: "test test test",
		Topics:      nil,
		DateCreated: time.Now(),
		PhotoPath:   "",
	}
	err := repo.Create(context.Background(), profile)
	assert.NoError(t, err)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM profiles WHERE user_id = $1", profile.UserId).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestPostgresEmptyFieldsCreate(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()
	repo := NewPostgresProfileRepository(db)
	profile := &models.Profile{
		Id:          uuid.New(),
		UserId:      "",
		Name:        "",
		Gender:      "",
		Description: "",
		Topics:      nil,
		DateCreated: time.Now(),
		PhotoPath:   "",
	}
	err := repo.Create(context.Background(), profile)
	assert.Error(t, err)
}

func TestPostgresErrUserAlreadyExistsCreate(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()
	repo := NewPostgresProfileRepository(db)
	profile := &models.Profile{
		Id:          uuid.New(),
		UserId:      "123",
		Name:        "test",
		Gender:      "f",
		Description: "test test test",
		Topics:      nil,
		DateCreated: time.Now(),
		PhotoPath:   "",
	}
	err := repo.Create(context.Background(), profile)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), profile)
	assert.Error(t, err)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM profiles").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestPostgresSuccessGet(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()
	repo := NewPostgresProfileRepository(db)
	newUuid := uuid.New()
	newYork, _ := time.LoadLocation("America/New_York")
	timeNow := time.Date(1920, 01, 01, 01, 01, 01, 01, newYork)
	profile := &models.Profile{
		Id:          newUuid,
		UserId:      "123",
		Name:        "test",
		Gender:      "f",
		Description: "test test test",
		Topics:      nil,
		DateCreated: timeNow,
		PhotoPath:   "",
	}
	err := repo.Create(context.Background(), profile)
	assert.NoError(t, err)
	gotProfile := repo.Get(context.Background(), "123")
	// Сравниваем отдельно поля
	assert.Equal(t, profile.Id, gotProfile.Id)
	assert.Equal(t, profile.UserId, gotProfile.UserId)
	assert.Equal(t, profile.Name, gotProfile.Name)
	assert.Equal(t, profile.Gender, gotProfile.Gender)
	assert.Equal(t, profile.Description, gotProfile.Description)

	// Сравниваем только даты (год, месяц, день)
	assert.Equal(t, profile.DateCreated.Year(), gotProfile.DateCreated.Year())
	assert.Equal(t, profile.DateCreated.Month(), gotProfile.DateCreated.Month())
	assert.Equal(t, profile.DateCreated.Day(), gotProfile.DateCreated.Day())
}

func TestPostgresNilGet(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()
	repo := NewPostgresProfileRepository(db)
	profile := repo.Get(context.Background(), "123")
	assert.Nil(t, profile)
}
