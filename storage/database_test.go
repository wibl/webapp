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

	mock.ExpectPrepare("INSERT INTO `group`").ExpectExec().WithArgs(newGroup.Title).WillReturnResult(sqlmock.NewResult(1, 1))

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

func TestGetAllGroups(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{"id", "title"}
	mock.ExpectQuery("SELECT \\* FROM `group`").WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 1)).RowsWillBeClosed()

	fakeDB := &dbStorage{db}
	_, err = fakeDB.GetAllGroups()
	if err != nil {
		t.Errorf("error was not expected while GetAllGroups: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var id int64 = 1
	testGroup := &model.Group{ID: id, Title: "test_group"}
	columns := []string{"id", "title"}

	mock.ExpectQuery("SELECT \\* FROM `group`").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(testGroup.ID, testGroup.Title)).
		RowsWillBeClosed()

	fakeDB := &dbStorage{db}
	fakeDB.CreateGroup(testGroup)
	gr, err := fakeDB.GetGroup(id)
	if err != nil {
		t.Errorf("error was not expected while GetGroup: %s", err)
	}

	if *gr != *testGroup {
		t.Errorf("testGroup is invalid: %v need: %v", gr, testGroup)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testGroup := &model.Group{ID: 1, Title: "test_group"}
	mock.ExpectExec("UPDATE `group`").WithArgs(testGroup.Title, testGroup.ID).WillReturnResult(sqlmock.NewResult(0, 0))

	fakeDB := &dbStorage{db}
	err = fakeDB.UpdateGroup(testGroup)
	if err != nil {
		t.Errorf("error was not expected while UpdateGroup: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testGroup := &model.Group{ID: 1, Title: "test_group"}
	mock.ExpectExec("DELETE FROM `group`").WithArgs(testGroup.ID).WillReturnResult(sqlmock.NewResult(0, 0))

	fakeDB := &dbStorage{db}
	err = fakeDB.DeleteGroup(testGroup)
	if err != nil {
		t.Errorf("error was not expected while UpdateGroup: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateTemplate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	fakeDB := &dbStorage{db}
	newGroup := &model.Group{Title: "test_group"}

	mock.ExpectPrepare("INSERT INTO `group`").ExpectExec().WithArgs(newGroup.Title).WillReturnResult(sqlmock.NewResult(1, 1))

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

func TestGetAllTemplate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testGroup := &model.Group{ID: 1, Title: "test_group"}
	columns := []string{"id", "groupid", "title", "queue", "body"}
	fakeDB := &dbStorage{db}

	mock.ExpectQuery("SELECT \\* FROM `template`").
		WithArgs(testGroup.ID).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 1, 1, 1, 1)).
		RowsWillBeClosed()	

	_, err = fakeDB.GetAllTemplates(testGroup)
	if err != nil {
		t.Errorf("error was not expected while GetAllGroups: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetTemplate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testTmp := &model.Template{ID: 1, GroupID: 1, Title: "testTemp", Queue: "testQueue", Body: "testBody"}
	columns := []string{"id", "groupid", "title", "queue", "body"}

	mock.ExpectQuery("SELECT \\* FROM `template`").
		WithArgs(testTmp.ID).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(testTmp.ID, testTmp.GroupID, testTmp.Title, testTmp.Queue, testTmp.Body)).
		RowsWillBeClosed()

	fakeDB := &dbStorage{db}
	fakeDB.CreateTemplate(testTmp)
	tmp, err := fakeDB.GetTemplate(testTmp.ID)
	if err != nil {
		t.Errorf("error was not expected while GetGroup: %s", err)
	}

	if *tmp != *testTmp {
		t.Errorf("testGroup is invalid: %v need: %v", tmp, testTmp)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateTemplate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testTmp := &model.Template{ID: 1, GroupID: 1, Title: "testTemp", Queue: "testQueue", Body: "testBody"}
	mock.ExpectExec("UPDATE `template`").WithArgs(testTmp.Title, testTmp.Queue, testTmp.ID).WillReturnResult(sqlmock.NewResult(0, 0))

	fakeDB := &dbStorage{db}
	err = fakeDB.UpdateTemplate(testTmp)
	if err != nil {
		t.Errorf("error was not expected while UpdateGroup: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteTemplate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testTmp := &model.Template{ID: 1, GroupID: 1, Title: "testTemp", Queue: "testQueue", Body: "testBody"}
	mock.ExpectExec("DELETE FROM `template`").WithArgs(testTmp.ID).WillReturnResult(sqlmock.NewResult(0, 0))

	fakeDB := &dbStorage{db}
	err = fakeDB.DeleteTemplate(testTmp)
	if err != nil {
		t.Errorf("error was not expected while UpdateGroup: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}