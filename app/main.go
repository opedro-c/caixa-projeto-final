package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
)


type Movie struct {
	ID       string
	Title    string
	Year     int
	Genre    string
	Director string
}

type Pagination struct {
	CurrentPage int
	TotalPages  int
	HasPrev     bool
	PrevPage    int
	HasNext     bool
	NextPage    int
}

type IndexData struct {
	Movies []Movie
	Page   Pagination
}

type TopMovie struct {
	Title string
	Count int
}

type GenreRating struct {
	Genre   string
	Average float64
}

type CountryStats struct {
	Country string
	Count   int
}

type StatsData struct {
	TopMovies  []TopMovie
	BestGenre  GenreRating
	TopCountry CountryStats
}

var db *sql.DB

func main() {
	var err error
	connStr := "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"
	db, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco:", err)
	}
	defer db.Close()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/estatisticas", statsHandler)

	fmt.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSize := 20
	offset := (page - 1) * pageSize

	var totalMovies int
	err := db.QueryRow("SELECT COUNT(*) FROM movies").Scan(&totalMovies)
	if err != nil {
		http.Error(w, "Erro ao contar filmes", http.StatusInternalServerError)
		return
	}
	totalPages := int(math.Ceil(float64(totalMovies) / float64(pageSize)))

	rows, err := db.Query(`
		SELECT movie_id, title, year, genre, director 
		FROM movies ORDER BY title LIMIT $1 OFFSET $2`, pageSize, offset)
	if err != nil {
		http.Error(w, "Erro ao buscar filmes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		rows.Scan(&m.ID, &m.Title, &m.Year, &m.Genre, &m.Director)
		movies = append(movies, m)
	}

	data := IndexData{
		Movies: movies,
		Page: Pagination{
			CurrentPage: page,
			TotalPages:  totalPages,
			HasPrev:     page > 1,
			PrevPage:    page - 1,
			HasNext:     page < totalPages,
			NextPage:    page + 1,
		},
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, data)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	data := StatsData{}

	rows, _ := db.Query(`
		SELECT m.title, COUNT(r.movie_id) as total 
		FROM movies m JOIN ratings r ON m.movie_id = r.movie_id 
		GROUP BY m.title ORDER BY total DESC LIMIT 5`)
	for rows.Next() {
		var tm TopMovie
		rows.Scan(&tm.Title, &tm.Count)
		data.TopMovies = append(data.TopMovies, tm)
	}

	db.QueryRow(`
		SELECT m.genre, AVG(r.rating) as media 
		FROM movies m JOIN ratings r ON m.movie_id = r.movie_id 
		GROUP BY m.genre ORDER BY media DESC LIMIT 1`).Scan(&data.BestGenre.Genre, &data.BestGenre.Average)

	db.QueryRow(`
		SELECT country, COUNT(*) as total 
		FROM users u JOIN ratings r ON u.user_id = r.user_id 
		GROUP BY country ORDER BY total DESC LIMIT 1`).Scan(&data.TopCountry.Country, &data.TopCountry.Count)

	tmpl := template.Must(template.ParseFiles("templates/statistics.html"))
	tmpl.Execute(w, data)
}