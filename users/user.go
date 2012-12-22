package users

import (
	"fmt"
	"github.com/YouthBuild-USA/godata/web"
	"github.com/gorilla/context"
	"net/http"
)

type contextKey int

const (
	userKey contextKey = iota
)

// User represents a web user.
type User struct {
	Id       int
	Username string
	password []byte
}

func (user User) String() string {
	return fmt.Sprintf("#%v %v/%v", user.Id, user.Username, user.password)
}

// CurrentUser returns the currently logged in user.  It attempts to load
// the user from the context first, and if that fails, it reads a userId
// from the session and loads the user, storing it in the context.  If
// there is no userId in the session (user is not authenticated) then nil
// is returned
func CurrentUser(r *http.Request) (*User, error) {
	if user := context.Get(r, userKey); user != nil {
		return user.(*User), nil
	} else {
		session := web.Session(r)
		userId, ok := session.Values["userId"]
		if !ok {
			// The current user is not authenticated....
			return nil, nil
		}

		user, err := LoadUser(userId.(int))
		if err != nil {
			return nil, err
		}
		context.Set(r, userKey, user)
		return user, nil
	}
	return nil, nil
}

var users = make([]*User, 0)

// LoadUser loads a user object from the database by Id.
func LoadUser(id int) (*User, error) {
	query := DB.QueryRow("SELECT id, username, password FROM users WHERE id = $1", id)
	var user User
	err := query.Scan(&user.Id, &user.Username, &user.password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// AuthenticateOrRedirect asserts that there is a user logged in.  If there
// is not a user logged in, then the user is redirected to the login page and
// the current URL is stored in the session.  Returns true if the user was
// redirected.
func AuthenticateOrRedirect(w http.ResponseWriter, r *http.Request, urlStr string) bool {
	user, err := CurrentUser(r)
	if user == nil || err != nil {
		session := web.Session(r)
		session.Values["loginDestinatation"] = r.URL.String()
		http.Redirect(w, r, urlStr, http.StatusFound)
		return true
	}
	return false
}
