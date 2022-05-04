package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"joyroku/todos/model"
	"net/http"

	"os"
	"strconv"
	"strings"
)

func (a *AppHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

var rd *render.Render = render.New()

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type AppHandler struct {
	http.Handler
	db model.DBHandler
}

var getSessionID = func(r *http.Request) string {
	session, err := store.Get(r, "session")
	if err != nil {
		return ""
	}
	val := session.Values["id"]
	if val == nil {
		return ""
	}
	return val.(string)
}

func (a *AppHandler) getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := getSessionID(r)
	//list := []*model.Todo{}
	//for _, v := range todoMap {
	//	list = append(list, v)
	//}
	list := a.db.GetTodos(sessionId)
	rd.JSON(w, http.StatusOK, list)
}
func (a *AppHandler) addTodoHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := getSessionID(r)
	name := r.FormValue("name")
	//id := len(todoMap) + 1
	//reqTodo := &Todo{id, name, false, time.Now()}
	//todoMap[id] = reqTodo
	reqTodo := a.db.AddTodo(name, sessionId)
	rd.JSON(w, http.StatusCreated, reqTodo)
}

type Success struct {
	Success bool `json:"success"`
}

func (a *AppHandler) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	ok := a.db.DeleteTodo(id)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
	//if _, ok := todoMap[id]; ok {
	//	delete(todoMap, id)
	//	rd.JSON(w, http.StatusOK, Success{true})
	//} else {
	//}
}
func (a *AppHandler) completeChangeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	complete := r.FormValue("complete") == "true"
	ok := a.db.CompleteTodo(id, complete)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
	//if todo, ok := todoMap[id]; ok {
	//	todo.Completed = complete
	//	rd.JSON(w, http.StatusOK, Success{true})
	//} else {
	//	rd.JSON(w, http.StatusOK, Success{false})
	//}
}

func (a *AppHandler) Close() {
	a.db.Close()
}
func CheckSignin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	//if request URL is /signin.html , then next()
	if strings.Contains(r.URL.Path, "/signin") || strings.Contains(r.URL.Path, "/auth") {
		next(w, r)
		fmt.Println("1")
		return
	}
	// if user already signed in
	sessionID := getSessionID(r)
	if sessionID != "" {
		next(w, r)
		fmt.Println("2")
		return
	}
	//if not user sign in
	http.Redirect(w, r, "/signin.html", http.StatusTemporaryRedirect)
}
func MakeHandler(filePath string) *AppHandler {
	r := mux.NewRouter()
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(CheckSignin),
		negroni.NewStatic(http.Dir("public")))
	n.UseHandler(r)

	a := &AppHandler{
		Handler: n,
		db:      model.NewDBHandler(filePath),
	}

	r.HandleFunc("/todos", a.getTodoListHandler).Methods("GET")
	r.HandleFunc("/todos", a.addTodoHandler).Methods("POST")
	r.HandleFunc("/todos/{id:[0-9]+}", a.deleteTodoHandler).Methods("DELETE")
	r.HandleFunc("/todoComplete/{id:[0-9]+}", a.completeChangeHandler).Methods("GET")
	r.HandleFunc("/auth/google/login", googleLoginHandler)
	r.HandleFunc("/auth/google/callback", googleAuthCallback)
	r.HandleFunc("/", a.indexHandler)
	return a
}
