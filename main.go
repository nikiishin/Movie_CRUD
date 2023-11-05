package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
)

var logger *zap.Logger

func initZapLog() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
		defer logger.Sync()

	}
}

type Movie struct {
	ID       string    `json:"ID"`
	Isbn     string    `json:"Isbn"`
	Title    string    `json:"Title"`
	Director *Director `json:"Director"`
}

type Director struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	// Логирование начала запроса
	logger.Info("GetMovies request started")
	defer logger.Sync()
	logger.Info("GetMovies request recevied", zap.String("Method", r.Method), zap.String("Path", r.URL.Path))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)

}
func getMovie(w http.ResponseWriter, r *http.Request) {

	logger.Info("GetMovie request started")
	defer logger.Sync()
	logger.Info("GetMovie request received", zap.String("Method", r.Method), zap.String("Path", r.URL.Path))
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			logger.Info("Movie found", zap.String("ID", item.ID), zap.String("Title", item.Title))
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Movie not found", http.StatusNotFound)
	logger.Info("Movie not found", zap.String("ID", params["ID"]))
}

func createMovies(w http.ResponseWriter, r *http.Request) {
	logger.Info("CreateMovie request start")
	defer logger.Sync()
	logger.Info("CreateMovie request received", zap.String("Method", r.Method), zap.String("Path", r.URL.Path))
	var newMovie Movie
	if err := json.NewDecoder(r.Body).Decode(&newMovie); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		logger.Error("Failed to decode request body", zap.Error(err))
		return
	}
	for _, existingMovie := range movies {
		if existingMovie.Title == newMovie.Title {
			http.Error(w, "Movie with the same title already exists", http.StatusConflict)
			logger.Warn("Movie with the same title already exists", zap.String("Title", newMovie.Title))
			return
		}

	}
	newMovie.ID = strconv.Itoa(len(movies) + 1)
	movies = append(movies, newMovie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newMovie); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		logger.Error("Failed to encode response", zap.Error(err))
	}
	logger.Info("Movie created", zap.String("ID", newMovie.ID), zap.String("Title", newMovie.Title))
}

func upgradeMovies(w http.ResponseWriter, r *http.Request) {
	logger.Info("UpgradeMovie request started")
	logger.Info("UpgradeMovie request received", zap.String("Method", r.Method), zap.String("Path", r.URL.Path))
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	idToUpdate := params["id"]
	var movieToUpdete *Movie
	for i, item := range movies {
		if item.ID == idToUpdate {
			movieToUpdete = &movies[i]
			break
		}
	}
	if movieToUpdete == nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		logger.Info("Movie not found", zap.String("ID", idToUpdate))
		return
	}
	var updatedMovie Movie
	if err := json.NewDecoder(r.Body).Decode(&updatedMovie); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		logger.Error("Failed to decode request body", zap.Error(err))
	}
	movieToUpdete.Title = updatedMovie.Title
	movieToUpdete.Isbn = updatedMovie.Isbn
	if updatedMovie.Director != nil {
		movieToUpdete.Director = updatedMovie.Director
	}

	logger.Info("Movies updated", zap.String("ID", movieToUpdete.ID), zap.String("Title", movieToUpdete.Title))
	json.NewEncoder(w).Encode(movieToUpdete)

}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	logger.Info("deleteMovies request started")
	defer logger.Sync()
	logger.Info("deleteMovies request received", zap.String("Method", r.Method), zap.String("Path", r.URL.Path))

	w.Header().Set("Content-Type", "application/json")
	// Извлекаем параметр "id" из URL, чтобы определить, какой фильм нужно удалить
	params := mux.Vars(r)
	idToDelete := params["id"]
	// Создаем переменные для хранения информации о фильме, который будет удален,
	// и его индекса в срезе movies
	var deletedMovie Movie
	deletedIndex := -1
	// Проходим по срезу movies, чтобы найти фильм с указанным "id"
	for index, movie := range movies {
		if movie.ID == idToDelete {
			deletedMovie = movie
			deletedIndex = index
			break
		}
	}
	// Проверяем, был ли найден фильм для удаления
	if deletedIndex == -1 {
		// Если фильм не был найден, возвращаем ошибку 404 Not Found и записываем лог
		http.Error(w, "Movie not found", http.StatusNotFound)
		logger.Info("Movie not found", zap.String("ID", idToDelete))
		return
	}
	// Если фильм был найден, удаляем его из среза movies
	movies = append(movies[:deletedIndex], movies[deletedIndex+1:]...)
	// Записываем информацию о фильме, который был удален, в лог
	logger.Info("Deleted movie", zap.String("ID", idToDelete), zap.String("Title", deletedMovie.Title))
	// Возвращаем информацию об удаленном фильме в виде JSON
	json.NewEncoder(w).Encode(deletedMovie)

}

func main() {
	initZapLog()
	r := mux.NewRouter()
	movies = append(movies, Movie{"1", "438227", "Movie One", &Director{"John", "Don"}})
	movies = append(movies, Movie{"2", "455785", "Movie Two", &Director{"Lich", "Locher"}})
	movies = append(movies, Movie{"3", "4557896", "Movie Three", &Director{"Memas", "Malcom"}})
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovies).Methods("POST")
	r.HandleFunc("/movies/{id}", upgradeMovies).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")
	fmt.Printf("Start servise in :8000 port \n")
	log.Fatal(http.ListenAndServe(":8000", r))

}
