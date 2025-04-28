package user_test

import (
	"fitapp-backend/api/resource/user"  // Adjust import path
	mockDB "fitapp-backend/mock/db"     // Adjust import path
	testUtil "fitapp-backend/util/test" // Adjust import path
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Updated list of columns reflecting the new schema and GORM mappings
var userColumns = []string{
	"id", "created_at", "updated_at", "deleted_at",
	"user_username", "user_full_name", "user_sex", "user_height", "user_weight", "user_age",
}

func TestRepository_List(t *testing.T) {
	t.Parallel()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	repo := user.NewRepository(db)

	now := time.Now()
	userID1 := uuid.New()
	userID2 := uuid.New()

	mockRows := sqlmock.NewRows(userColumns).
		AddRow(userID1, now, now, gorm.DeletedAt{}, "user1", "Full Name One", true, 180, 80, 30). // Sex=true (male)
		AddRow(userID2, now, now, gorm.DeletedAt{}, "user2", "Full Name Two", false, 165, 60, 25) // Sex=false (female)

	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL`)
	mock.ExpectQuery(expectedSQL).WillReturnRows(mockRows)

	users, err := repo.List()
	testUtil.NoError(t, err)
	testUtil.Equal(t, 2, len(users))
	testUtil.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Create(t *testing.T) {
	t.Parallel()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	repo := user.NewRepository(db)

	id := uuid.New()
	newUser := &user.User{
		ID:       id,
		Username: "newuser",
		FullName: "New User Full",
		Sex:      true, // Male
		Height:   175,
		Weight:   75,
		Age:      28,
	}

	mock.ExpectBegin()
	// Match the column order GORM uses for INSERT (check generated SQL if needed)
	expectedSQL := regexp.QuoteMeta(`INSERT INTO "users" ("id","created_at","updated_at","deleted_at","user_username","user_full_name","user_sex","user_height","user_weight","user_age") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`)
	mock.ExpectExec(expectedSQL).
		WithArgs(
			newUser.ID,
			mockDB.AnyTime{}, // CreatedAt
			mockDB.AnyTime{}, // UpdatedAt
			nil,              // DeletedAt
			newUser.Username,
			newUser.FullName,
			newUser.Sex,
			newUser.Height,
			newUser.Weight,
			newUser.Age,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	created, err := repo.Create(newUser)
	testUtil.NoError(t, err)
	testUtil.NotNil(t, created)
	testUtil.Equal(t, id, created.ID)
	testUtil.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Read(t *testing.T) {
	t.Parallel()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	repo := user.NewRepository(db)

	id := uuid.New()
	expectedUsername := "readuser"
	expectedHeight := 190

	mockRows := sqlmock.NewRows(userColumns).
		AddRow(id, time.Now(), time.Now(), gorm.DeletedAt{}, expectedUsername, "Read User Name", true, expectedHeight, 90, 40)

	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)
	mock.ExpectQuery(expectedSQL).WithArgs(id).WillReturnRows(mockRows)

	found, err := repo.Read(id)
	testUtil.NoError(t, err)
	testUtil.NotNil(t, found)
	testUtil.Equal(t, id, found.ID)
	testUtil.Equal(t, expectedUsername, found.Username)
	testUtil.Equal(t, expectedHeight, found.Height)
	testUtil.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update(t *testing.T) {
	t.Parallel()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	repo := user.NewRepository(db)

	id := uuid.New()
	userToUpdate := &user.User{
		ID:       id,
		Username: "updateduser",
		FullName: "Updated Full Name",
		Sex:      false, // Female
		Height:   168,
		Weight:   65,
		Age:      31,
	}

	mock.ExpectBegin()
	// Match the fields selected in Repository.Update
	expectedSQL := regexp.QuoteMeta(`UPDATE "users" SET "user_username"=$1,"user_full_name"=$2,"user_sex"=$3,"user_height"=$4,"user_weight"=$5,"user_age"=$6,"updated_at"=$7 WHERE id = $8`)
	mock.ExpectExec(expectedSQL).
		WithArgs(
			userToUpdate.Username,
			userToUpdate.FullName,
			userToUpdate.Sex,
			userToUpdate.Height,
			userToUpdate.Weight,
			userToUpdate.Age,
			mockDB.AnyTime{}, // UpdatedAt
			id,               // WHERE id = ?
		).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
	mock.ExpectCommit()

	rowsAffected, err := repo.Update(userToUpdate)
	testUtil.NoError(t, err)
	testUtil.Equal(t, int64(1), rowsAffected)
	testUtil.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete(t *testing.T) {
	t.Parallel()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	repo := user.NewRepository(db)

	id := uuid.New()

	mock.ExpectBegin()
	expectedSQL := regexp.QuoteMeta(`UPDATE "users" SET "deleted_at"=$1 WHERE id = $2 AND "users"."deleted_at" IS NULL`)
	mock.ExpectExec(expectedSQL).
		WithArgs(mockDB.AnyTime{}, id).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	rowsAffected, err := repo.Delete(id)
	testUtil.NoError(t, err)
	testUtil.Equal(t, int64(1), rowsAffected)
	testUtil.NoError(t, mock.ExpectationsWereMet())
}
