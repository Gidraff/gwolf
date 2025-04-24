package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-sql-driver/mysql"
)

var (
	db          *sql.DB
	clientID    = "log-backend"
	keyCloakURL = "http://localhost:8080/realms/myrealm"
)

type WorkoutLog struct {
	ID           int64
	Name         string
	Comment      string
	Date         time.Time
	NumberOfSets int64
	NumberOfReps int64
	Weight       int64
	Effort       int64
}

type UserClaims struct {
	Email       string `json:"email"`
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

func authMiddleware(next http.Handler) http.Handler {
	provider, err := oidc.NewProvider(context.Background(), keyCloakURL)
	if err != nil {
		panic(err)
	}

	// verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
	verifier := provider.Verifier(&oidc.Config{SkipClientIDCheck: true})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		idToken, err := verifier.Verify(r.Context(), tokenStr)
		if err != nil {
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		var claims UserClaims
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "Failed to parse claims: "+err.Error(), http.StatusUnauthorized)
			return
		}
		log.Printf("%+v", claims)

		// Add claims to context
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getWorkoutLogs query for logs
func getWorkoutLogs() ([]WorkoutLog, error) {
	var workoutLogs []WorkoutLog
	rows, err := db.Query("SELECT * FROM log")
	if err != nil {
		return nil, fmt.Errorf("log %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var log WorkoutLog
		if err := rows.Scan(&log.ID, &log.Name, &log.Comment, &log.Date, &log.NumberOfSets, &log.NumberOfReps, &log.Weight, &log.Effort); err != nil {
			return nil, fmt.Errorf("log :%v", err)
		}
		workoutLogs = append(workoutLogs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("log: %v", err)
	}
	return workoutLogs, nil
}

func getWorkoutLogById(id int) (WorkoutLog, error) {
	var log WorkoutLog
	row := db.QueryRow("SELECT * FROM log WHERE id = ?", id)
	if err := row.Scan(&log.ID, &log.Name, &log.Comment, &log.Date, &log.NumberOfSets, &log.NumberOfReps, &log.Weight, &log.Effort); err != nil {
		if err == sql.ErrNoRows {
			return log, fmt.Errorf("log ID %d: no such entry", id)
		}
		return log, fmt.Errorf("log ID %d: %v", id, err)
	}
	return log, nil
}

func addLogEntry(log WorkoutLog) (int64, error) {
	result, err := db.Exec("INSERT INTO log (workout_name, comment, date, number_of_sets, number_of_reps, weight, effort) VALUES (?,?,?,?,?,?,?)",
		log.Name, log.Comment, log.Date, log.NumberOfSets, log.NumberOfReps, log.Weight, log.Effort)
	if err != nil {
		return 0, fmt.Errorf("add log entry :%v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("add log entry :%v", err)
	}
	return id, nil
}
func main() {

	cfg := mysql.Config{
		User:      os.Getenv("DB_USER"),
		Passwd:    os.Getenv("DB_PASSWD"),
		Net:       "tcp",
		Addr:      "127.0.0.1:3306",
		DBName:    "recordings",
		ParseTime: true,
	}

	var err error

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(authMiddleware)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Get("/logs", func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("user").(UserClaims)
		if !ok {
			http.Error(w, "Invalid user context", http.StatusInternalServerError)
			return
		}

		for _, role := range claims.RealmAccess.Roles {
			if role == "admin" {
				log.Printf("Checking role: %s", role)
				wlogs, err := getWorkoutLogs()
				if err != nil {
					http.Error(w, "failed to get logs", http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json") // ✅ fix typo here too
				json.NewEncoder(w).Encode(wlogs)
				return
			}
		}

		http.Error(w, "Forbidden", http.StatusForbidden)
	})
	http.ListenAndServe(":3000", r)

}
