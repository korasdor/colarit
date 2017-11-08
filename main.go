package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/korasdor/colarit/handler"
	"github.com/korasdor/colarit/services"
	"os"
	"github.com/korasdor/colarit/utils"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	services.InitDb()
	defer services.CloseDb()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	result := utils.UpdateBooksTemplate()
	if result == false {
		utils.PrintOutput("Update book error")
	}

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handler.IndexHandler)
	muxRouter.HandleFunc("/db_state", handler.DBStateHandler)
	muxRouter.HandleFunc("/books", handler.BooksHandler)
	muxRouter.HandleFunc("/update_books_template", handler.UpdateBooksHandler)

	muxRouter.HandleFunc("/create_table/{table_name}/{access_token}", handler.CreateTableHandler)
	muxRouter.HandleFunc("/fill_table/{table_name}/{serials_count}/{range_id}/{dealer_id}/{access_token}", handler.FillTableHandler)
	muxRouter.HandleFunc("/fill_table_with_file/{table_name}/{file_name}/{range_id}/{dealer_id}/{access_token}", handler.FillTableWithFileHandler)
	muxRouter.HandleFunc("/get_serial/{table_name}/{range_id}/{dealer_id}/{serials_format}/{access_token}", handler.GetSerialsHandler)

	muxRouter.HandleFunc("/activate_serial/{table_name}/{serial_key}", handler.ActivateSerialsHandler)
	muxRouter.HandleFunc("/reset_serial/{table_name}/{serial_key}", handler.ResetSerialsHandler)
	muxRouter.HandleFunc("/about_serial/{table_name}/{serial_key}", handler.AboutSerialsHandler)

	muxRouter.HandleFunc("/sendmail", handler.SendMailHandler)

	//s := http.StripPrefix("/static/asset_bundles/", http.FileServer(http.Dir("./static/asset_bundles/")))
	//muxRouter.PathPrefix("/static/asset_bundles/").Handler(s).HandlerFunc(handler.ServeBundle);
	//muxRouter.HandleFunc("/static/media/{file}", handler.ServeStaticFiles)

	fmt.Printf("listening on %s...\n", port)

	muxWithTimeout := http.TimeoutHandler(muxRouter, time.Second*20, "{is_activated:false}")
	err := http.ListenAndServe(":"+port, muxWithTimeout)
	if err != nil {
		utils.PrintOutput(err.Error())
	}
}
