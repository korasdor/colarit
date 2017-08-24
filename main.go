package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/korasdor/colarit/model"
	"github.com/korasdor/colarit/handler"
	"github.com/korasdor/colarit/services"
	"os"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	m := model.New()
	handler.Model = m

	services.InitDb()
	defer services.CloseDb()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handler.IndexHandler)
	r.HandleFunc("/test", handler.TestHandler)
	r.HandleFunc("/db_state", handler.DBStateHandler)
	r.HandleFunc("/books", handler.BooksHandler)
	r.HandleFunc("/update_books_template", handler.UpdateBooksHandler)
	r.HandleFunc("/get_temp/{file_name}", handler.GetTempHandler)

	r.HandleFunc("/create_serial/{serials_name}/{serials_count}", handler.CreateSerialsHandler)
	r.HandleFunc("/get_serial/{serials_name}/{serials_format}", handler.GetSerialsHandler)
	r.HandleFunc("/create_table/{serials_name}/{table_name}", handler.CreateTableHandler)
	r.HandleFunc("/fill_table/{serials_name}/{table_name}/{dealer_id}", handler.FillTableHandler)

	r.HandleFunc("/activate_serial/{table_name}/{serial_key}", handler.ActivateSerialsHandler)
	r.HandleFunc("/reset_serial/{table_name}/{serial_key}", handler.ResetSerialsHandler)
	r.HandleFunc("/about_serial/{table_name}/{serial_key}", handler.AboutSerialsHandler)

	r.HandleFunc("/sendmail", handler.SendMailHandler)

	//s := http.StripPrefix("/static/asset_bundles/", http.FileServer(http.Dir("./static/asset_bundles/")))
	//r.PathPrefix("/static/asset_bundles/").Handler(s).HandlerFunc(handler.ServeBundle);
	r.HandleFunc("/static/media/{file}", handler.ServeStaticFiles)

	fmt.Printf("listening on %s...\n", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}
