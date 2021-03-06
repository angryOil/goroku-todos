package model

import (
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

type pqHandler struct {
	db *sql.DB
}

func (s *pqHandler) GetTodos(session string) []*Todo {
	var todos []*Todo
	rows, err := s.db.Query("select id , name , completed , createdAt from todos where sessionId = $1", session)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Name, &todo.Completed, &todo.CreatedAt)
		todos = append(todos, &todo)
	}
	return todos
}

func (s *pqHandler) AddTodo(name string, sessionId string) *Todo {
	var stmt, err = s.db.Prepare("insert into todos (sessionId,name, completed, createdAt) values ($1,$2,$3,now()) returning id")
	if err != nil {
		panic(err)
	}
	var id int
	err = stmt.QueryRow(sessionId, name, false).Scan(&id)
	if err != nil {
		panic(err)
	}
	var todo Todo
	todo.ID = id
	todo.Name = name
	todo.CreatedAt = time.Now()
	todo.Completed = false
	return &todo
}

func (s *pqHandler) DeleteTodo(id int) bool {
	stmt, err := s.db.Prepare("delete from todos where id = $1")
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(id)
	cnt, _ := rst.RowsAffected()
	return cnt > 0
}

//git 커밋을 위한 주석

func (s *pqHandler) CompleteTodo(id int, complete bool) bool {
	stmt, err := s.db.Prepare("update todos set completed = $1 where id = $2")
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(id, complete)
	if err != nil {
		panic(err)
	}
	cnt, _ := rst.RowsAffected()
	return cnt > 0
}

func (s *pqHandler) Close() {
	s.db.Close()
}

func newPqHandler(dbConn string) DBHandler {
	database, err := sql.Open("postgres", dbConn)
	if err != nil {
		panic(err)
	}
	statement, err := database.Prepare(
		`create table if not exists todos(
				id serial primary key  ,
				sessionId varchar(256),
				name text,
				completed boolean,
				createdAt timestamp );`)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}
	statement, err = database.Prepare(`create index if not exists sessionIdIndexOnTodos on todos(sessionId asc);`)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}
	return &pqHandler{db: database}
}
