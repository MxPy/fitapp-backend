package user_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"

	"fitapp-backend/api/resource/user"
	mockDB "fitapp-backend/mock/db"
	testUtil "fitapp-backend/util/test"
)

func TestRepository_List(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := user.NewRepository(db)

	mockRows := sqlmock.NewRows([]string{"id", "username", "full_name"}).
		AddRow(uuid.New(), "user1", "full_name1").
		AddRow(uuid.New(), "user2", "full_name2")

	mock.ExpectQuery("^SELECT (.+) FROM \"users\"").WillReturnRows(mockRows)

	users, err := repo.List()
	testUtil.NoError(t, err)
	testUtil.Equal(t, len(users), 2)
}

func TestRepository_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := user.NewRepository(db)

	id := uuid.New()
	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO \"users\" ").
		WithArgs(id, "username", "full_name", mockDB.AnyTime{}, mockDB.AnyTime{}, nil).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	user := &user.User{ID: id, Username: "username", Full_name: "full_name"}
	_, err = repo.Create(user)
	testUtil.NoError(t, err)
}

func TestRepository_Read(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := user.NewRepository(db)

	id := uuid.New()
	mockRows := sqlmock.NewRows([]string{"id", "username", "full_name"}).
		AddRow(id, "user1", "full_name1")

	mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").
		WithArgs(id).
		WillReturnRows(mockRows)

	user, err := repo.Read(id)
	testUtil.NoError(t, err)
	testUtil.Equal(t, "user1", user.Username)
}

func TestRepository_Update(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := user.NewRepository(db)

	id := uuid.New()
	_ = sqlmock.NewRows([]string{"id", "username", "full_name"}).
		AddRow(id, "user1", "full_name1")

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"users\" SET").
		WithArgs("username", "full_name", "", mockDB.AnyTime{}, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	user := &user.User{ID: id, Username: "username", Full_name: "full_name"}
	rows, err := repo.Update(user)
	testUtil.NoError(t, err)
	testUtil.Equal(t, 1, rows)
}

func TestRepository_Delete(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := user.NewRepository(db)

	id := uuid.New()
	_ = sqlmock.NewRows([]string{"id", "username", "full_name"}).
		AddRow(id, "user1", "full_name1")

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"users\" SET \"deleted_at\"").
		WithArgs(mockDB.AnyTime{}, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	rows, err := repo.Delete(id)
	testUtil.NoError(t, err)
	testUtil.Equal(t, 1, rows)
}
