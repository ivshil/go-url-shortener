package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	dbHost    = ""
	dbPort    = ""
	dbUser    = ""
	dbPass    = ""
	dbName    = ""
	goAppPort = ""
)

type User struct {
	ID        int       `json:"user_id"`
	Name      string    `json:"user_name"`
	Email     string    `json:"user_email"`
	BirthDate time.Time `json:"user_bdate"`
}

type Task struct {
	ID            int       `json:"task_id"`
	UserCreatorID int       `json:"user_creator_id"`
	Description   string    `json:"task_description"`
	StartDate     time.Time `json:"task_start_date"`
	DeadlineDate  time.Time `json:"task_deadline_date"`
}

type TaskContributor struct {
	ID           int       `json:"task_con_id"`
	UserID       int       `json:"user_id"`
	TaskID       int       `json:"task_id"`
	DateAssigned time.Time `json:"assigned_date"`
}

type URLShort struct {
	ID            int       `json:"url_short_id"`
	UserCreatorID int       `json:"user_creator_id"`
	URL           string    `json:"url_base"`
	URLS          string    `json:"url_short"`
	Created       time.Time `json:"url_created_date"`
}

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dbHost = os.Getenv("PGDB_HOST")
	dbPort = os.Getenv("PGDB_PORT")
	dbUser = os.Getenv("PGDB_USER")
	dbPass = os.Getenv("PGDB_PASS")
	dbName = os.Getenv("PGDB_NAME")
	goAppPort = os.Getenv("GOAPP_PORT")

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// check env variables
	// http.HandleFunc("/db/info", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method == http.MethodGet {
	// 		message := fmt.Sprintf("DB Host: %s\nDB Port: %s\nDB User: %s\nDB Password: %s\nDB Name: %s\n", dbHost, dbPort, dbUser, dbPass, dbName)
	// 		fmt.Println(message)
	// 		w.Header().Set("Content-Type", "text/plain")
	// 		w.WriteHeader(http.StatusOK)
	// 		w.Write([]byte(message))
	// 	} else {
	// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	}
	// })

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getUsersHandler(w, r, db)
		case http.MethodPost:
			addUserHandler(w, r, db)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getTasksHandler(w, r, db)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks_contributors", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getTaskContributorsHandler(w, r, db)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/urls", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getUrlShortsHandler(w, r, db)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start the HTTP server
	log.Printf("Server listening on port %s...\n", goAppPort)
	log.Fatal(http.ListenAndServe(":"+goAppPort, nil))
}

// openDB opens a connection to the PostgreSQL database
func openDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	return sql.Open("postgres", connStr)
}

// getUsersHandler retrieves all users from the Users table and returns them as a JSON response
func getUsersHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching users: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.BirthDate)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning users: %v", err), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Convert users slice to JSON and write the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding users to JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

func addUserHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON request: %v", err), http.StatusBadRequest)
		return
	}

	query := "INSERT INTO users (user_name, user_email, user_bdate) VALUES ($1, $2, $3) RETURNING user_id"
	err := db.QueryRow(query, user.Name, user.Email, user.BirthDate).Scan(&user.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting user into database: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response to JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

func getUrlShortsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query("SELECT * FROM url_shorts")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching url_shorts: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var urlShorts []URLShort
	for rows.Next() {
		var urlShort URLShort
		err := rows.Scan(&urlShort.ID, &urlShort.URL, &urlShort.URLS, &urlShort.UserCreatorID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning url_shorts: %v", err), http.StatusInternalServerError)
			return
		}
		urlShorts = append(urlShorts, urlShort)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(urlShorts)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding users to JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

func getTasksHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching tasks: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.UserCreatorID, &task.Description, &task.StartDate, &task.DeadlineDate)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning tasks: %v", err), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding tasks to JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

func getTaskContributorsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query("SELECT * FROM tasks_contributors")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching task contributors: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var taskContributors []TaskContributor
	for rows.Next() {
		var contributor TaskContributor
		err := rows.Scan(&contributor.ID, &contributor.UserID, &contributor.TaskID, &contributor.DateAssigned)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning task contributors: %v", err), http.StatusInternalServerError)
			return
		}
		taskContributors = append(taskContributors, contributor)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(taskContributors)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding task contributors to JSON: %v", err), http.StatusInternalServerError)
		return
	}
}
