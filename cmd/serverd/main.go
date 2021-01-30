package main

import (
	"encoding/json"
	"github.com/helder-jaspion/go-springfield-bank/pkg/adapter/repository"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	accountMemRepo := repository.NewAccountMemoryRepository()
	accountUC := usecase.NewAccountUseCase(accountMemRepo)

	r := httprouter.New()
	r.HandlerFunc("POST", "/accounts", createAccoundHandler(accountUC))

	log.Fatal(http.ListenAndServe(":8080", r))
}

func createAccoundHandler(accountUC usecase.AccountUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		var input usecase.AccountCreateInput
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			//json.NewEncoder(w).Encode(errors.ServiceError{Message: "Error unmarshalling the request"})
			return
		}

		result, err2 := accountUC.Create(r.Context(), input)
		if err2 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			//json.NewEncoder(w).Encode(errors.ServiceError{Message: "Error saving the input"})
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(result)
	}
}
