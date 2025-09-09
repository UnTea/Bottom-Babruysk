package main

import (
	"context"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"bottom_babruysk/database"
	"bottom_babruysk/logger"
)

func main() {
	l, err := logger.New()
	if err != nil {
		panic(err)
	}

	environmentDatabaseConnectionString := "postgres://admin:admin@localhost:5432/bottom_babruysk?sslmode=disable"

	dbConfig := database.Config{
		ConnectionString: environmentDatabaseConnectionString,
		Timeout:          time.Second * 30,
	}

	ctx := context.Background()
	db, err := database.New(ctx, dbConfig)
	if err != nil {
		l.Fatal("failed to connect to database", zap.Error(err))
	}

	defer db.Close()

	ctxExec, cancelExec := context.WithTimeout(ctx, db.QueryTimeout())
	defer cancelExec()

	testID := "11111111-1111-1111-1111-111111111111"
	testEmail := "test@example.local"
	testPasswordHash := "dummy-hash"
	testDisplayName := "Test User"

	insertSQL := `
		insert into users (id, email, password_hash, display_name)
		values ($1, $2, $3, $4)
		on conflict (id) do nothing;
	`

	if _, err := db.Driver().Exec(ctxExec, insertSQL, testID, testEmail, testPasswordHash, testDisplayName); err != nil {
		l.Fatal("failed to insert test user", zap.Error(err))
		return
	}

	request := &database.GetUserRequest{UserID: testID}

	response, err := request.GetUser(ctx, db)
	if err != nil {
		l.Error("failed to GetUser", zap.Error(err))
		return
	}

	if response != nil && response.Users != nil {
		uid := "<nil>"
		if response.Users.Id != nil {
			uid = *response.Users.Id
		}

		email := "<nil>"
		if response.Users.Email != nil {
			email = *response.Users.Email
		}

		display := "<nil>"
		if response.Users.DisplayName != nil {
			display = *response.Users.DisplayName
		}

		l.Info("got", zap.Any("user id", uid), zap.Any("user email", email), zap.Any("user display name", display))
	} else {
		l.Info("empty response")
	}
}
