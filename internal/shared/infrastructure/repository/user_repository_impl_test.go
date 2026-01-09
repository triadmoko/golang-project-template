package repository

import (
	"app/internal/shared/domain/entity"
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db      *gorm.DB
	mock    sqlmock.Sqlmock
	repo    *userRepository
	ctx     context.Context
	sqlDB   *sql.DB
}

func (s *UserRepositoryTestSuite) SetupTest() {
	var err error
	s.sqlDB, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	dialector := postgres.New(postgres.Config{
		Conn:       s.sqlDB,
		DriverName: "postgres",
	})

	s.db, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(s.T(), err)

	s.repo = &userRepository{db: s.db}
	s.ctx = context.Background()
}

func (s *UserRepositoryTestSuite) TearDownTest() {
	s.sqlDB.Close()
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (s *UserRepositoryTestSuite) TestCreate_Success() {
	user := &entity.User{
		ID:        "user-123",
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "users" ("id","email","username","password","first_name","last_name","phone","status","birth_date","gender","role","provider","is_active","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`)).
		WithArgs(
			user.ID,
			user.Email,
			user.Username,
			user.Password,
			user.FirstName,
			user.LastName,
			nil,      // phone
			"active", // status (default value)
			nil,      // birth_date
			"",       // gender
			"user",   // role (default value)
			"",       // provider
			user.IsActive,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := s.repo.Create(s.ctx, user)

	assert.NoError(s.T(), err)
	assert.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *UserRepositoryTestSuite) TestCreate_Error() {
	user := &entity.User{
		ID:        "user-123",
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "users"`)).
		WillReturnError(sql.ErrConnDone)
	s.mock.ExpectRollback()

	err := s.repo.Create(s.ctx, user)

	assert.Error(s.T(), err)
}

func (s *UserRepositoryTestSuite) TestGetByID_Success() {
	userID := "user-123"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "email", "username", "password", "first_name", "last_name", "is_active", "created_at", "updated_at", "deleted_at"}).
		AddRow(userID, "test@example.com", "testuser", "hashedpassword", "Test", "User", true, now, now, nil)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(userID, 1).
		WillReturnRows(rows)

	user, err := s.repo.GetByID(s.ctx, userID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), user)
	assert.Equal(s.T(), userID, user.ID)
	assert.Equal(s.T(), "test@example.com", user.Email)
	assert.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *UserRepositoryTestSuite) TestGetByID_NotFound() {
	userID := "nonexistent-id"

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(userID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := s.repo.GetByID(s.ctx, userID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), user)
}

func (s *UserRepositoryTestSuite) TestGetByEmail_Success() {
	email := "test@example.com"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "email", "username", "password", "first_name", "last_name", "is_active", "created_at", "updated_at", "deleted_at"}).
		AddRow("user-123", email, "testuser", "hashedpassword", "Test", "User", true, now, now, nil)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(email, 1).
		WillReturnRows(rows)

	user, err := s.repo.GetByEmail(s.ctx, email)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), user)
	assert.Equal(s.T(), email, user.Email)
	assert.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *UserRepositoryTestSuite) TestGetByEmail_NotFound() {
	email := "nonexistent@example.com"

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(email, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := s.repo.GetByEmail(s.ctx, email)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), user)
}

func (s *UserRepositoryTestSuite) TestGetByUsername_Success() {
	username := "testuser"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "email", "username", "password", "first_name", "last_name", "is_active", "created_at", "updated_at", "deleted_at"}).
		AddRow("user-123", "test@example.com", username, "hashedpassword", "Test", "User", true, now, now, nil)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE username = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(username, 1).
		WillReturnRows(rows)

	user, err := s.repo.GetByUsername(s.ctx, username)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), user)
	assert.Equal(s.T(), username, user.Username)
	assert.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *UserRepositoryTestSuite) TestGetByUsername_NotFound() {
	username := "nonexistentuser"

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE username = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(username, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := s.repo.GetByUsername(s.ctx, username)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), user)
}

func (s *UserRepositoryTestSuite) TestUpdate_Success() {
	user := &entity.User{
		ID:        "user-123",
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "hashedpassword",
		FirstName: "Updated",
		LastName:  "User",
		IsActive:  true,
	}

	filter := entity.FilterUser{
		ID: "user-123",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "users" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := s.repo.Update(s.ctx, filter, user)

	assert.NoError(s.T(), err)
}

func (s *UserRepositoryTestSuite) TestUpdate_Error() {
	user := &entity.User{
		ID:        "user-123",
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "hashedpassword",
		FirstName: "Updated",
		LastName:  "User",
		IsActive:  true,
	}

	filter := entity.FilterUser{
		ID: "user-123",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "users" SET`)).
		WillReturnError(sql.ErrConnDone)
	s.mock.ExpectRollback()

	err := s.repo.Update(s.ctx, filter, user)

	assert.Error(s.T(), err)
}

func (s *UserRepositoryTestSuite) TestDelete_Success() {
	userID := "user-123"

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "users" SET "deleted_at"=$1 WHERE id = $2 AND "users"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := s.repo.Delete(s.ctx, userID)

	assert.NoError(s.T(), err)
	assert.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *UserRepositoryTestSuite) TestDelete_Error() {
	userID := "user-123"

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "users" SET "deleted_at"=$1 WHERE id = $2 AND "users"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), userID).
		WillReturnError(sql.ErrConnDone)
	s.mock.ExpectRollback()

	err := s.repo.Delete(s.ctx, userID)

	assert.Error(s.T(), err)
}

func (s *UserRepositoryTestSuite) TestList_Success() {
	now := time.Now()
	filter := entity.FilterUser{
		Offset:  5,
		PerPage: 10,
	}

	rows := sqlmock.NewRows([]string{"id", "email", "username", "password", "first_name", "last_name", "is_active", "created_at", "updated_at", "deleted_at"}).
		AddRow("user-1", "user1@example.com", "user1", "hashedpassword", "User", "One", true, now, now, nil).
		AddRow("user-2", "user2@example.com", "user2", "hashedpassword", "User", "Two", true, now, now, nil)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE deleted_at IS NULL AND "users"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`)).
		WithArgs(filter.PerPage, filter.Offset).
		WillReturnRows(rows)

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM "users" WHERE deleted_at IS NULL AND "users"."deleted_at" IS NULL`)).
		WillReturnRows(countRows)

	users, totalRows, err := s.repo.List(s.ctx, filter)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 2)
	assert.Equal(s.T(), 2, totalRows)
	assert.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *UserRepositoryTestSuite) TestList_Empty() {
	filter := entity.FilterUser{
		Offset:  0,
		PerPage: 10,
	}

	rows := sqlmock.NewRows([]string{"id", "email", "username", "password", "first_name", "last_name", "is_active", "created_at", "updated_at", "deleted_at"})

	// When offset is 0, GORM doesn't add OFFSET to the query
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE deleted_at IS NULL AND "users"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT $1`)).
		WithArgs(filter.PerPage).
		WillReturnRows(rows)

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM "users" WHERE deleted_at IS NULL AND "users"."deleted_at" IS NULL`)).
		WillReturnRows(countRows)

	users, totalRows, err := s.repo.List(s.ctx, filter)

	assert.NoError(s.T(), err)
	assert.Empty(s.T(), users)
	assert.Equal(s.T(), 0, totalRows)
	assert.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *UserRepositoryTestSuite) TestList_Error() {
	filter := entity.FilterUser{
		Offset:  0,
		PerPage: 10,
	}

	// When offset is 0, GORM doesn't add OFFSET to the query
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE deleted_at IS NULL AND "users"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT $1`)).
		WithArgs(filter.PerPage).
		WillReturnError(sql.ErrConnDone)

	users, totalRows, err := s.repo.List(s.ctx, filter)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), users)
	assert.Equal(s.T(), 0, totalRows)
}
