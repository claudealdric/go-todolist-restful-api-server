package data

import (
	"database/sql"
	"log"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"golang.org/x/crypto/bcrypt"
)

func InitDb(db *sql.DB) {
	createUsersTable(db)
	seedUsersTable(db)
	createTasksTable(db)
	seedTasksTable(db)
}

func createTasksTable(db *sql.DB) {
	_, err := db.Exec(`
		create table if not exists tasks (
			id integer primary key autoincrement,
			title text not null
		)
	`)
	if err != nil {
		log.Fatalln("failed creating the tasks table:", err)
	}
}

func createUsersTable(db *sql.DB) {
	_, err := db.Exec(`
		create table if not exists users (
			id integer primary key autoincrement,
			name text not null,
			email text not null unique,
			password text not null
		)
	`)
	if err != nil {
		log.Fatalln("failed creating the users table:", err)
	}
}

func seedUsersTable(db *sql.DB) {
	dto := models.NewCreateUserDTO(
		"Claude Aldric",
		"cvaldric@gmail.com",
		"Caput Draconis",
	)
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(dto.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Fatalln("failed at hashing the password:", err)
	}
	_, err = db.Exec(`
		insert into users (name, email, password)
		select ?, ?, ?
		where not exists (select 1 from users)
	`, dto.Name, dto.Email, hashedPassword)
	if err != nil {
		log.Fatalln("failed seeding the users table:", err)
	}
}

func seedTasksTable(db *sql.DB) {
	_, err := db.Exec(`
		insert into tasks (title)
		select 'This is the first task'
		where not exists (select 1 from tasks)
	`)
	if err != nil {
		log.Fatalln("failed seeding the users table:", err)
	}
}
