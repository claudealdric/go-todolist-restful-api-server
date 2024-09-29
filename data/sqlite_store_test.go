package data

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/claudealdric/go-todolist-restful-api-server/testutils/assert"
)

func TestValidateUserCredentials(t *testing.T) {
	dbFile := "../tmp/sqlite_store_test.db"
	db, err := sql.Open("sqlite3", dbFile)
	assert.HasNoError(t, err)
	defer db.Close()
	InitDb(db)
	defer cleanSqliteDatabase(dbFile)
	store := NewSqliteStore(db)

	tests := []struct {
		name     string
		email    string
		password string
		want     bool
	}{
		{
			name:     "when provided valid credentials",
			email:    "cvaldric@gmail.com",
			password: "Caput Draconis",
			want:     true,
		},
		{
			name:     "when provided an invalid password",
			email:    "cvaldric@gmail.com",
			password: "password",
			want:     false,
		},
		{
			name:     "when provided an invalid email",
			email:    "doesnotexist@email.com",
			password: "Caput Draconis",
			want:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equals(
				t,
				store.ValidateUserCredentials(test.email, test.password),
				test.want,
			)
		})
	}

}

func cleanSqliteDatabase(path string) {
	os.Remove(path)
}
