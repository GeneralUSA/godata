package users

import (
	"fmt"
	"github.com/YouthBuild-USA/godata/log"
	"github.com/YouthBuild-USA/godata/templates"
	"github.com/YouthBuild-USA/godata/web"
	"net/http"
)

var userLog = log.New("User")

var (
	loginFormTemplate = templates.OneColumn().Add("users/loginForm")
	userTemplate      = templates.OneColumn().Add("users/user")
)

func init() {
}

func (module userModule) AddRoutes(handleAdder web.HandlerAdder) error {

	handleAdder("/login", loginForm).Methods("GET")
	handleAdder("/login", login).Methods("POST")
	handleAdder("/user", userPage).Methods("GET")
	handleAdder("/logout", logout).Methods("GET")
	return nil
}

func loginForm(w http.ResponseWriter, r *http.Request) error {
	user, _ := CurrentUser(r)
	if user != nil {
		http.Redirect(w, r, "/user", http.StatusFound)
		return nil
	}
	loginFormTemplate.Render(w, r, "", nil)

	return nil
}

func login(w http.ResponseWriter, r *http.Request) error {
	username := r.FormValue("username")
	password := r.FormValue("password")

	userLog.Info("Log in: %v/%v", username, password)

	query := DB.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username)
	var user User
	err := query.Scan(&user.Id, &user.Username, &user.password)

	if err != nil {
		web.FlashWarning(r, "No such user found")
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil
	}

	if string(user.password) == password {
		session := web.Session(r)
		session.Values["userId"] = user.Id
		web.FlashInfo(r, fmt.Sprintf("Logged in as %v", user.Username))

		if dest, ok := session.Values["loginDestination"]; ok {
			http.Redirect(w, r, dest.(string), http.StatusFound)
		} else {
			http.Redirect(w, r, "/user", http.StatusFound)
		}
		return nil
	}
	web.FlashWarning(r, "Incorrect username or password")
	http.Redirect(w, r, "/login", http.StatusFound)
	return nil
}

func userPage(w http.ResponseWriter, r *http.Request) error {
	if AuthenticateOrRedirect(w, r, "/login") {
		return nil
	}
	user, _ := CurrentUser(r)
	return userTemplate.Render(w, r, "User Info", user)
}

func logout(w http.ResponseWriter, r *http.Request) error {
	session := web.Session(r)
	delete(session.Values, "userId")
	web.FlashInfo(r, "You have been logged out")
	http.Redirect(w, r, "/login", http.StatusFound)
	return nil
}
