package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/korasdor/colarit/handler"
	"github.com/korasdor/colarit/services"
	"os"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	services.InitDb()
	defer services.CloseDb()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handler.IndexHandler)
	r.HandleFunc("/db_state", handler.DBStateHandler)
	r.HandleFunc("/books", handler.BooksHandler)
	r.HandleFunc("/update_books_template", handler.UpdateBooksHandler)

	r.HandleFunc("/create_table/{table_name}/{access_token}", handler.CreateTableHandler)
	r.HandleFunc("/fill_table/{table_name}/{serials_count}/{range_id}/{dealer_id}/{access_token}", handler.FillTableHandler)

	r.HandleFunc("/get_serial/{serials_name}/{serials_format}", handler.GetSerialsHandler)

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
