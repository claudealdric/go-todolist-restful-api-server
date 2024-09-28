package data

import (
	"database/sql"
	"log"
)

func InitDb(db *sql.DB) {
	createUsersTable(db)
	seedUsersTable(db)
	createTasksTable(db)
	seedTasksTable(db)
}

func createTasksTable(db *sql.DB) {
	_, err := db.Exec(`
		create table tasks (
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
		create table users (
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
	_, err := db.Exec(`
		insert into users (name, email, password)
		values
			('Claude Aldric', 'cvaldric@gmail.com', 'Caput Draconis')
	`)
	if err != nil {
		log.Fatalln("failed seeding the users table:", err)
	}

	_, err = db.Exec(`
		insert into users (name, email, password)
		values
			('John Doe', 'john@email.com', 'password')
	`)
	if err != nil {
		log.Fatalln("failed seeding the users table:", err)
	}
}

func seedTasksTable(db *sql.DB) {
	_, err := db.Exec(`
		insert into tasks (title)
		values
			('This is the first task')
	`)
	if err != nil {
		log.Fatalln("failed seeding the users table:", err)
	}
}
