package integration

import (
	"context"
	"database/sql"
	grpc3 "date-bot-go/matching/client/grpc"
	"date-bot-go/matching/models"
	repository2 "date-bot-go/matching/repository"
	service2 "date-bot-go/matching/service"
	"date-bot-go/pkg/profilepb"
	"date-bot-go/profile/repository"
	grpc2 "date-bot-go/profile/server/grpc"
	"date-bot-go/profile/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
	"time"
)

func setupDB(t *testing.T, schema string) (*sql.DB, func()) {
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
	createSchema(t, db, schema)
	cleanup := func() {
		db.Close()
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	}
	return db, cleanup
}

func createSchema(t *testing.T, db *sql.DB, schema string) {
	_, err := db.Exec(schema)
	require.NoError(t, err)
}

func TestCall(t *testing.T) {
	schemaProfile := `
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
	schemaMatching := `
		CREATE TABLE IF NOT EXISTS public.likes
		(
			user_id text COLLATE pg_catalog."default" NOT NULL,
			liked_id text COLLATE pg_catalog."default" NOT NULL,
			UNIQUE (user_id, liked_id)
		);
	`
	db, cleanup := setupDB(t, schemaProfile)
	createSchema(t, db, schemaMatching)
	defer cleanup()
	lis := bufconn.Listen(1024 * 1024)
	profileRepository := repository.NewPostgresProfileRepository(db)
	profileService := service.NewProfileService(profileRepository)
	profileHandler := grpc2.NewProfileHandler(profileService)
	grpcServer := grpc.NewServer()
	profilepb.RegisterProfileServiceServer(grpcServer, profileHandler)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Fatalf("Server exited with error %v", err)
		}
	}()
	defer grpcServer.Stop()

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial with error %v", err)
	}
	defer conn.Close()
	matchingClient := profilepb.NewProfileServiceClient(conn)
	matchingRepository := repository2.NewPostgresMatchingRepository(db)
	profileProvider := grpc3.NewGrpcProfileProvider(matchingClient, conn)
	matchingService := service2.NewMatchingService(matchingRepository, profileProvider)

	//here
	err = profileService.Create(ctx, "123", "test1", "f", "test test test")
	assert.NoError(t, err)
	expProfile, err := matchingService.NextProfile(ctx, "123")
	assert.NoError(t, err)
	assert.Equal(t, "returned", expProfile.Username)

}

func TestProfileMatchingFlow(t *testing.T) {
	schemaProfile := `
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
	schemaMatching := `
		CREATE TABLE IF NOT EXISTS public.likes
		(
			user_id text COLLATE pg_catalog."default" NOT NULL,
			liked_id text COLLATE pg_catalog."default" NOT NULL,
			UNIQUE (user_id, liked_id)
		);
	`
	db, cleanup := setupDB(t, schemaProfile)
	createSchema(t, db, schemaMatching)
	defer cleanup()
	lis := bufconn.Listen(1024 * 1024)
	profileRepository := repository.NewPostgresProfileRepository(db)
	profileService := service.NewProfileService(profileRepository)
	profileHandler := grpc2.NewProfileHandler(profileService)
	grpcServer := grpc.NewServer()
	profilepb.RegisterProfileServiceServer(grpcServer, profileHandler)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Fatalf("Server exited with error %v", err)
		}
	}()
	defer grpcServer.Stop()

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial with error %v", err)
	}
	defer conn.Close()
	matchingClient := profilepb.NewProfileServiceClient(conn)
	matchingRepository := repository2.NewPostgresMatchingRepository(db)
	profileProvider := grpc3.NewGrpcProfileProvider(matchingClient, conn)
	matchingService := service2.NewMatchingService(matchingRepository, profileProvider)

	err = profileService.Create(ctx, "123", "test1", "f", "test test test")
	assert.NoError(t, err)
	err = profileService.Create(ctx, "456", "test2", "f", "test test test")
	assert.NoError(t, err)
	err = profileService.Create(ctx, "789", "test3", "f", "test test test")
	assert.NoError(t, err)
	err = profileService.Create(ctx, "321", "test4", "f", "test test test")
	assert.NoError(t, err)
	expProfile, err := matchingService.NextProfile(ctx, "123")
	assert.NoError(t, err)
	assert.IsType(t, expProfile, &models.Profile{})
	err = matchingService.Like(ctx, "456", "123")
	assert.NoError(t, err)
	err = matchingService.Like(ctx, "123", "456")
	assert.NoError(t, err)
	expProfile, err = matchingService.NextProfile(ctx, "123")
	assert.NoError(t, err)
	assert.IsType(t, expProfile, &models.Profile{})
	err = matchingService.Like(ctx, "123", expProfile.UserId)

	assert.NoError(t, err)
	err = matchingService.Like(ctx, expProfile.UserId, "123")
	assert.NoError(t, err)
}
