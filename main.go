package main

//import package and library go
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	bcrypt "golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
)

var db *sql.DB

var err error

var tpl *template.Template

//conect db and set template
func init() {
	db, err = sql.Open("mysql", "root:jadir123@tcp(127.0.0.1:3306)/go_db")
	checkErr(err)
	err = db.Ping()
	checkErr(err)
	tpl = template.Must(template.ParseGlob("templates/*"))
}

//function checked error
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

//main function first execute
func main() {

	defer db.Close()
	//get from index() function
	http.HandleFunc("/", index)

	//get from userForm() function
	http.HandleFunc("/userForm", userForm)

	//get from createUsers() function
	http.HandleFunc("/createUsers", createUsers)

	//get from editUsers() function
	http.HandleFunc("/editUsers", editUsers)

	//get from updateUsers() function
	http.HandleFunc("/updateUsers", updateUsers)

	//get from deleteUsers() function
	http.HandleFunc("/deleteUsers", deleteUsers)

	//run server in 127.0.0.1:9000
	log.Println("Server is up on 9000 port")
	log.Fatalln(http.ListenAndServe(":9000", nil))
}

//deklarasi users variabel properti
type user struct {
	ID        int64
	Username  string
	FirstName string
	LastName  string
	Password  []byte
}

//list user
func index(w http.ResponseWriter, req *http.Request) {
	rows, e := db.Query(
		`SELECT id,
		username,
		first_name,
		last_name,
		password
		FROM users;`)

	if e != nil {
		log.Println(e)
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	users := make([]user, 0)
	for rows.Next() {
		usr := user{}
		rows.Scan(&usr.ID, &usr.Username, &usr.FirstName, &usr.LastName, &usr.Password)
		users = append(users, usr)
	}
	log.Println(users)
	tpl.ExecuteTemplate(w, "index.html", users)
}

//form create user
func userForm(w http.ResponseWriter, req *http.Request) {
	err = tpl.ExecuteTemplate(w, "userForm.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//action create users
func createUsers(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		usr := user{}
		usr.Username = req.FormValue("username")
		usr.FirstName = req.FormValue("firstName")
		usr.LastName = req.FormValue("lastName")
		bPass, e := bcrypt.GenerateFromPassword([]byte(req.FormValue("password")), bcrypt.MinCost)

		if e != nil {
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}
		usr.Password = bPass

		_, e = db.Exec("INSERT INTO users (username,first_name, last_name, password) VALUES (?,?,?,?)", usr.Username,
			usr.FirstName,
			usr.LastName,
			usr.Password,
		)

		if e != nil {
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Methof not supported", http.StatusMethodNotAllowed)
}

//form edut users
func editUsers(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	rows, err := db.Query(
		`SELECT id,
	 	username,
		first_name,
		last_name
		FROM users
		WHERE id = ` + id + `;`)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	usr := user{}
	for rows.Next() {
		rows.Scan(&usr.ID, &usr.Username, &usr.FirstName, &usr.LastName)
	}
	tpl.ExecuteTemplate(w, "editUser.html", usr)
}

//action edit users
func updateUsers(w http.ResponseWriter, req *http.Request) {
	_, er := db.Exec("UPDATE users set username = ?, first_name = ?, last_name = ? WHERE id = ?",
		req.FormValue("username"),
		req.FormValue("firstName"),
		req.FormValue("lastName"),
		req.FormValue("id"),
	)

	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

//action deleted users
func deleteUsers(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")

	if id == "" {
		http.Error(w, "Please Send ID", http.StatusBadRequest)
		return
	}

	_, er := db.Exec("DELETE FROM users WHERE id = ?", id)

	if er != nil {
		http.Error(w, er.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}
