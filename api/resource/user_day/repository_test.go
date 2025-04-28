package userday_test

import (
	"fitapp-backend/api/resource/userday" // Adjust import path
	mockDB "fitapp-backend/mock/db"       // Adjust import path
	testUtil "fitapp-backend/util/test"   // Adjust import path
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Common columns for user_days table
var userDayColumns = []string{"id", "created_at", "updated_at", "deleted_at", "user_id", "user_date", "daily_kcal", "daily_proteins", "daily_carbs", "daily_fats"}

// Helper function to get a fixed time for date columns
func getTestDate(t *testing.T) time.Time {
	t.Helper()
	// Use a fixed date for consistent testing, time part is ignored by DATE type
	date, err := time.Parse("2006-01-02", "2024-03-15")
	testUtil.NoError(t, err)
	return date
}

func TestRepository_List(t *testing.T) {
	t.Parallel()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	repo := userday.NewRepository(db)

	testDate := getTestDate(t)
	userID1 := uuid.New()
	userID2 := uuid.New()

	mockRows := sqlmock.NewRows(userDayColumns).
		AddRow(uuid.New(), time.Now(), time.Now(), gorm.DeletedAt{}, userID1, testDate, 2000, 150, 200, 80).
		AddRow(uuid.New(), time.Now(), time.Now(), gorm.DeletedAt{}, userID2, testDate, 2500, 180, 250, 100)

	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "user_days" WHERE "user_days"."deleted_at" IS NULL`)
	mock.ExpectQuery(expectedSQL).WillReturnRows(mockRows)

	userDays, err := repo.List()
	testUtil.NoError(t, err)
	testUtil.Equal(t, 2, len(userDays))
	testUtil.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Create(t *testing.T) {
	t.Parallel()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	repo := userday.NewRepository(db)

	id := uuid.New()
	userID := uuid.New()
	testDate := getTestDate(t)
	newUserDay := &userday.UserDay{
		ID:            id,
		UserID:        userID,
		UserDate:      testDate,
		DailyKcal:     2200,
		DailyProteins: 160,
		DailyCarbs:    210,
		DailyFats:     90,
	}

	mock.ExpectBegin()
	// Adjust columns based on GORM behavior - it might omit zero values unless specified
	expectedSQL := regexp.QuoteMeta(`INSERT INTO "user_days" ("id","created_at","updated_at","deleted_at","user_id","user_date","daily_kcal","daily_proteins","daily_carbs","daily_fats") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`)
	mock.ExpectExec(expectedSQL).
		WithArgs(
			newUserDay.ID,
			mockDB.AnyTime{}, // CreatedAt
			mockDB.AnyTime{}, // UpdatedAt
			nil,              // DeletedAt
			newUserDay.UserID,
			newUserDay.UserDate, // GORM might send time.Time directly or formatted string
			newUserDay.DailyKcal,
			newUserDay.DailyProteins,
			newUserDay.DailyCarbs,
			newUserDay.DailyFats,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	created, err := repo.Create(newUserDay)
	testUtil.NoError(t, err)
	testUtil.NotNil(t, created)
	testUtil.Equal(t, id, created.ID)
	testUtil.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Read(t *testing.T) {
	t.Parallel()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	repo := userday.NewRepository(db)

	id := uuid.New()
	userID := uuid.New()
	testDate := getTestDate(t)
	expectedKcal := 2100

	mockRows := sqlmock.NewRows(userDayColumns).
		AddRow(id, time.Now(), time.Now(), gorm.DeletedAt{}, userID, testDate, expectedKcal, 155, 205, 85)

	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "user_days" WHERE id = $1 AND "user_days"."deleted_at" IS NULL ORDER BY "user_days"."id" LIMIT 1`)
	mock.ExpectQuery(expectedSQL).WithArgs(id).WillReturnRows(mockRows)

	found, err := repo.Read(id)
	testUtil.NoError(t, err)
	testUtil.NotNil(t, found)
	testUtil.Equal(t, id, found.ID)
	testUtil.Equal(t, expectedKcal, found.DailyKcal)
	testUtil.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Read_NotFound(t *testing.T) {
	t.Parallel()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	repo := userday.NewRepository(db)
	id := uuid.New()

	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "user_days" WHERE id = $1 AND "user_days"."deleted_at" IS NULL ORDER BY "user_days"."id" LIMIT 1`)
	mock.ExpectQuery(expectedSQL).WithArgs(id).WillReturnError(gorm.ErrRecordNotFound)

	_, err = repo.Read(id)
	testUtil.ErrorIs(t, err, gorm.ErrRecordNotFound)
	testUtil.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_FindByUserAndDate(t *testing.T) {
	t.Parallel()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	repo := userday.NewRepository(db)

	id := uuid.New()
	userID := uuid.New()
	testDate := getTestDate(t)
	dateStr := testDate.Format(userday.DateFormat)
	expectedKcal := 2300

	mockRows := sqlmock.NewRows(userDayColumns).
		AddRow(id, time.Now(), time.Now(), gorm.DeletedAt{}, userID, testDate, expectedKcal, 170, 220, 95)

	// Note: Query uses user_id and user_date in WHERE clause
	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "user_days" WHERE user_id = $1 AND user_date = $2 AND "user_days"."deleted_at" IS NULL ORDER BY "user_days"."id" LIMIT 1`)
	mock.ExpectQuery(expectedSQL).WithArgs(userID, dateStr).WillReturnRows(mockRows)

	found, err := repo.FindByUserAndDate(userID, testDate)
	testUtil.NoError(t, err)
	testUtil.NotNil(t, found)
	testUtil.Equal(t, id, found.ID) // Should find the correct primary ID
	testUtil.Equal(t, userID, found.UserID)
	testUtil.Equal(t, expectedKcal, found.DailyKcal)
	testUtil.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_FindByUserAndDate_NotFound(t *testing.T) {
	t.Parallel()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	repo := userday.NewRepository(db)

	userID := uuid.New()
	testDate := getTestDate(t)
	dateStr := testDate.Format(userday.DateFormat)

	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "user_days" WHERE user_id = $1 AND user_date = $2 AND "user_days"."deleted_at" IS NULL ORDER BY "user_days"."id" LIMIT 1`)
	mock.ExpectQuery(expectedSQL).WithArgs(userID, dateStr).WillReturnError(gorm.ErrRecordNotFound)

	_, err = repo.FindByUserAndDate(userID, testDate)
	testUtil.ErrorIs(t, err, gorm.ErrRecordNotFound)
	testUtil.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update(t *testing.T) {
	t.Parallel()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	repo := userday.NewRepository(db)

	id := uuid.New()
	userDayToUpdate := &userday.UserDay{
		ID:            id, // ID is crucial for WHERE clause
		DailyKcal:     2500,
		DailyProteins: 190,
		DailyCarbs:    260,
		DailyFats:     110,
		// UserID and UserDate should not be included in the SET clause per repository logic
	}

	mock.ExpectBegin()
	// Update should only set selected fields + UpdatedAt
	expectedSQL := regexp.QuoteMeta(`UPDATE "user_days" SET "daily_kcal"=$1,"daily_proteins"=$2,"daily_carbs"=$3,"daily_fats"=$4,"updated_at"=$5 WHERE id = $6`)
	mock.ExpectExec(expectedSQL).
		WithArgs(
			userDayToUpdate.DailyKcal,
			userDayToUpdate.DailyProteins,
			userDayToUpdate.DailyCarbs,
			userDayToUpdate.DailyFats,
			mockDB.AnyTime{}, // UpdatedAt
			id,               // WHERE id = ?
		).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
	mock.ExpectCommit()

	rowsAffected, err := repo.Update(userDayToUpdate)
	testUtil.NoError(t, err)
	testUtil.Equal(t, int64(1), rowsAffected)
	testUtil.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete(t *testing.T) {
	t.Parallel()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	repo := userday.NewRepository(db)

	id := uuid.New()

	mock.ExpectBegin()
	expectedSQL := regexp.QuoteMeta(`UPDATE "user_days" SET "deleted_at"=$1 WHERE id = $2 AND "user_days"."deleted_at" IS NULL`)
	mock.ExpectExec(expectedSQL).
		WithArgs(mockDB.AnyTime{}, id).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
	mock.ExpectCommit()

	rowsAffected, err := repo.Delete(id)
	testUtil.NoError(t, err)
	testUtil.Equal(t, int64(1), rowsAffected)
	testUtil.NoError(t, mock.ExpectationsWereMet())
}
