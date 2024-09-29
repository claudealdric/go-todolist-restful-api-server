package data

import (
	"database/sql"
	"log"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"golang.org/x/crypto/bcrypt"
)

type SqliteStore struct {
	db *sql.DB
}

func NewSqliteStore(db *sql.DB) *SqliteStore {
	s := SqliteStore{db}
	return &s
}

func (s *SqliteStore) CreateTask(dto *models.CreateTaskDTO) (*models.Task, error) {
	result, err := s.db.Exec(`
		insert into tasks (title)
		values
			(?)
	`, dto.Title)
	if err != nil {
		return nil, err
	}

	taskId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return models.NewTask(int(taskId), dto.Title), nil
}

func (s *SqliteStore) CreateUser(dto *models.CreateUserDTO) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(dto.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	result, err := s.db.Exec(`
		insert into users (name, email, password)
		values
			(?, ?, ?)
	`, dto.Name, dto.Email, hashedPassword)
	if err != nil {
		return nil, err
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return models.NewUser(int(userId), dto.Name, dto.Email, dto.Password), nil
}

func (s *SqliteStore) DeleteTaskById(id int) error {
	_, err := s.db.Exec(`
		delete from tasks where id = ?
	`, id)
	return err
}

func (s *SqliteStore) GetTaskById(id int) (*models.Task, error) {
	var task models.Task
	err := s.db.QueryRow(`select * from tasks where id = ?`, id).Scan(
		&task.Id,
		&task.Title,
	)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *SqliteStore) GetTasks() ([]models.Task, error) {
	rows, err := s.db.Query(`select * from tasks`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.Id, &task.Title)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *SqliteStore) GetUsers() ([]models.User, error) {
	rows, err := s.db.Query(`select * from users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *SqliteStore) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(`select * from users where email = ?`, email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *SqliteStore) UpdateTask(task *models.Task) (*models.Task, error) {
	_, err := s.db.Exec(`
		update tasks
		set title = ?
		where id = ?
	`, task.Title, task.Id)
	if err != nil {
		return nil, err
	}
	updatedTask, err := s.GetTaskById(task.Id)
	if err != nil {
		return nil, err
	}

	return updatedTask, nil
}

func (s *SqliteStore) ValidateUserCredentials(email, password string) bool {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		log.Printf("error retrieving user for validation: %v\n", err)
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
}
