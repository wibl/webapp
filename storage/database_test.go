package storage

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/wibl/webapp/model"
)

func TestCreateGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	fakeDB := &dbStorage{db}
	newGroup := &model.Group{Title: "test_group"}
	//mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO group").WithArgs(newGroup.Title).WillReturnResult(sqlmock.NewResult(1, 1))
	//mock.ExpectCommit()

	// now we execute our method
	if err = fakeDB.CreateGroup(newGroup); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if newGroup.ID != 1 {
		t.Errorf("id expected 1 got %d", newGroup.ID)
	}
}
