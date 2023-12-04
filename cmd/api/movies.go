package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/thesoftwaremasons/greenlight/internal/data"
	"github.com/thesoftwaremasons/greenlight/internal/validator"
)

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filter
	}
	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCsv(qs, "genres", []string{})
	input.Filter.Page = app.readInt(qs, "page", 1, v)
	input.Filter.PageSize = app.readInt(qs, "size", 20, v)
	input.Filter.Sort = app.readString(qs, "sort", "id")

	input.Filter.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	if data.ValidateFilters(v, input.Filter); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	movies, metaData, err := app.model.Movies.GetAll(input.Title, input.Genres, input.Filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJson(w, http.StatusOK, envelope{"movies": movies, "metadata": metaData}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres,omitempty"`
	}

	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: int32(input.Runtime),
		Genres:  input.Genres,
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.model.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))
	err = app.writeJson(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)

	}
}
func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {

	//get id
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundError(w, r)
		return
	}
	//get movies by id
	movie, err := app.model.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// If the request contains a X-Expected-Version header, verify that the movie
	// version in the database matches the expected version specified in the header.
	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.FormatInt(int64(movie.Version), 32) != r.Header.Get("X-Expected-Version") {
			app.editConflictResponse(w, r)
			return
		}
	}

	//changing properties to pointers because we want nullable instead on the zeroth value
	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres,omitempty"`
	}

	err = app.readJson(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {

		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = int32(*input.Runtime)
	}
	if input.Genres != nil {

		movie.Genres = input.Genres
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.model.Movies.Update(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	err = app.writeJson(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)

	}

}
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	Id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundError(w, r)
		return
	}
	movie, err := app.model.Movies.Get(Id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundError(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}
	err = app.writeJson(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundError(w, r)
		return
	}
	err = app.model.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundError(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return a 200 OK status code along with a success message.
	err = app.writeJson(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
