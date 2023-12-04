package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	if err := app.writeJson(w, http.StatusOK, envelope{"health status": data}, nil); err != nil {

		app.serverErrorResponse(w, r, err)
		//http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

}
