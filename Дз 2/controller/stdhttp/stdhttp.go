package stdhttp

import (
	"adreesbook/models/dto"
	"adreesbook/pkg"
	"adreesbook/psg"
	"encoding/json"
	"fmt"
	"net/http"
)

// Controller обрабатывает HTTP запросы для адресной книги.
type Controller struct {
	DB  *psg.Psg
	Srv *http.Server
}

// NewController создает новый Controller.
func NewController(addr string, db *psg.Psg) *Controller {
	return &Controller{
		DB:  db,
		Srv: &http.Server{Addr: addr},
	}
}

// RecordAdd обрабатывает HTTP запрос для добавления новой записи.
func (c *Controller) RecordAdd(w http.ResponseWriter, r *http.Request) {
	record := dto.Record{}
	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record.Phone = pkg.PhoneNormalize(record.Phone)

	_, err = c.DB.RecordAdd(record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Record added successfully")
}

// RecordsGet обрабатывает HTTP запрос для получения записей на основе предоставленных полей Record.
func (c *Controller) RecordsGet(w http.ResponseWriter, r *http.Request) {
	record := dto.Record{}
	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	records, err := c.DB.RecordsGet(record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(records)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// RecordUpdate обрабатывает HTTP запрос для обновления записи.
func (c *Controller) RecordUpdate(w http.ResponseWriter, r *http.Request) {
	record := dto.Record{}
	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record.Phone = pkg.PhoneNormalize(record.Phone)

	err = c.DB.RecordUpdate(record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Record updated successfully")
}

// RecordDeleteByPhone обрабатывает HTTP запрос для удаления записи по номеру телефона.
func (c *Controller) RecordDeleteByPhone(w http.ResponseWriter, r *http.Request) {
	// Извлекаем номер телефона из URL-параметра
	phone := r.URL.Query().Get("phone")

	// Проверяем, что номер телефона был указан
	if phone == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Удаляем запись с указанным номером телефона из базы данных
	err := c.DB.DeleteRecordByPhone(phone)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Возвращаем статус 200 OK в случае успеха
	w.WriteHeader(http.StatusOK)
}
