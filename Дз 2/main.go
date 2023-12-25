package main

import (
	"controller/stdhttp"
	"gate/psg"
	"net/http"
)

func main() {
	db, err := psg.NewPsg("psql://postgres:knyaz1234@localhost:5432/adressbook", "postgres", "knyaz1234")
	if err != nil {
		panic(err)
	}
	controller := stdhttp.NewController(":8080", db)

	http.HandleFunc("/record/add", controller.RecordAdd)
	http.HandleFunc("/records", controller.RecordsGet)
	http.HandleFunc("/record/update", controller.RecordUpdate)
	http.HandleFunc("/record/delete", controller.RecordDeleteByPhone)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
