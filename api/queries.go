package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
)

type TaskItem struct {
	ID          int64      `json:"id"`
	UserID      int        `json:"userID"`
	Title       string     `json:"title"`
	CreatedAt   time.Time  `json:"createdAt"`
	CompletedAt *time.Time `json:"completedAt"`
}

func getTasks(userID int) []TaskItem {
	tasks := make([]TaskItem, 0, 10)
	rows, err := dbPool.Query(context.Background(), `select id,user_id,title,created_at,completed_at from tasks where user_id = $1 order by completed_at desc, created_at desc;`, userID)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var taskItem TaskItem
		rows.Scan(&taskItem.ID, &taskItem.UserID, &taskItem.Title, &taskItem.CreatedAt, &taskItem.CompletedAt)
		tasks = append(tasks, taskItem)
	}
	return tasks
}

func insertTask(taskName string, userID int) error {
	_, err := dbPool.Exec(context.Background(), `insert into tasks(user_id, title, created_at) values ($1, $2, $3);`, userID, taskName, time.Now())
	return err
}

func updateTask(taskID int64, userID int, completed bool) error {
	var completedAt *time.Time
	if completed == true {
		now := time.Now()
		completedAt = &now
	}
	_, err := dbPool.Exec(context.Background(), `update tasks set completed_at = $1 where id = $2 and user_id = $3;`, completedAt, taskID, userID)
	return err
}

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"_"`
	CreatedAt    time.Time `json:"createdAt"`
	LastLogin    time.Time `json:"lastLogin"`
}

func insertUser(email, passwordHash string) {
	dbPool.Exec(context.Background(), `INSERT INTO users
	(email, "password", created_at, last_login)
	VALUES($1, $2, $3, $4);`, email, passwordHash, time.Now(), time.Now())
}

func getUserByEmail(email string) User {
	var user User
	err := dbPool.QueryRow(context.Background(), `select id, email, "password", created_at, last_login from users where email = $1 limit 1;`, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.LastLogin)
	if err != nil && err != pgx.ErrNoRows {
		log.Fatal(err)
	}
	return user
}
