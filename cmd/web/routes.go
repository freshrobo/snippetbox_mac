package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/justinas/alice"
	"snippetbox.alberttseng.net/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// fileServer := http.FileServer(http.Dir("./ui/static/"))
	// mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	mux.HandleFunc("GET /ping", ping)

	dynamic := alice.New(app.albert, app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /about", dynamic.ThenFunc(app.about))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("GET /account/view", protected.ThenFunc(app.accountView))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))
	mux.Handle("GET /account/password/update", protected.ThenFunc(app.accountPasswordUpdate))
	mux.Handle("POST /account/password/update", protected.ThenFunc(app.accountPasswordUpdatePost))

	mux.Handle("POST /snippet/delete/{id}", protected.ThenFunc(app.snippetDeletePost))

	mux.HandleFunc("POST /testing/bocheng", func(w http.ResponseWriter, r *http.Request) {
		var Input struct {
			Account  string `json:"user_number" binding:"required,email"`
			Password string `json:"secret" binding:"required"`
			Name     string `json:"user_name" binding:"required"`
			Age      int    `json:"age" binding:"required,min=18"`
		}
		var Output struct {
			Message string `json:"message"`
		}
		inputJsonBytes, err := io.ReadAll(r.Body)
		if err != nil {
			Output.Message = "Bad Request Bocheng"
			w.WriteHeader(http.StatusBadRequest)
			jsonString, err := json.Marshal(Output)
			if err != nil {
				log.Fatal(err)
			}
			_, err = fmt.Fprint(w, string(jsonString))
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		log.Println("inputJsonString===============>")
		log.Println(string(inputJsonBytes))
		err = json.Unmarshal(inputJsonBytes, &Input)
		if err != nil {
			Output.Message = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			jsonString, err := json.Marshal(Output)
			if err != nil {
				log.Fatal(err)
			}
			_, err = fmt.Fprint(w, string(jsonString))
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		fmt.Printf("%+v\n", Input)
		Output.Message = "OK"
		responseJsonBytes, _ := json.Marshal(Output)
		fmt.Fprint(w, string(responseJsonBytes))
		return
	})
	//====== following code is use GIN =====//
	// gr := gin.Default()
	// gr.POST("/testing/bocheng", func(c *gin.Context) {
	// 	var Input struct {
	// 		Account  string `json:"user_number" binding:"required,email"`
	// 		Password string `json:"secret" binding:"required"`
	// 		Name     string `json:"user_name" binding:"required"`
	// 		Age      int    `json:"age" binding:"required,min=18"`
	// 	}
	// 	var Output struct {
	// 		Message string `json:"message"`
	// 	}
	// 	if err := c.ShouldBindJSON(&Input); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{
	// 			"Message": err.Error(),
	// 		})
	// 		return
	// 	}

	// 	Output.Message = "OK"
	// 	c.JSON(http.StatusOK, Output)
	// 	return
	// })
	// mux.Handle("/testing/bocheng", gr)
	//===== GIN ======//

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
