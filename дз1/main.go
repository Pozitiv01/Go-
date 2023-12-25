package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Note struct {
	ID      int
	Name    string
	Surname string
	Text    string
}

var notes []Note

func main() {
	fmt.Println("Консольное приложение")
	http.HandleFunc("/", handleRequest)
	go func() {
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			panic(err)
		}
	}()
	menu()
}

func menu() {
	for {
		fmt.Println("\nМеню:")
		fmt.Println("1) Добавить заметку")
		fmt.Println("2) Посмотреть заметку")
		fmt.Println("3) Удалить заметку")
		fmt.Println("4) Завершить работу")

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Выберите пункт меню: ")
		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1\n":
			addNote()
		case "2\n":
			viewNote()
		case "3\n":
			deleteNote()
		case "4\n":
			fmt.Println("Работа приложения завершена.")
			return
		default:
			fmt.Println("Некорректный выбор. Попробуйте еще раз.")
		}
	}
}

func addNote() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите имя: ")
	name, _ := reader.ReadString('\n')
	name = name[:len(name)-1]

	fmt.Print("Введите фамилию: ")
	surname, _ := reader.ReadString('\n')
	surname = surname[:len(surname)-1]

	fmt.Print("Введите заметку: ")
	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1]

	id := len(notes) + 1
	note := Note{ID: id, Name: name, Surname: surname, Text: text}
	notes = append(notes, note)

	fmt.Println("Заметка добавлена. ID заметки:", id)
}

func viewNote() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите ID заметки: ")
	idStr, _ := reader.ReadString('\n')
	var id int
	fmt.Sscanf(strings.TrimSpace(idStr), "%d", &id)

	found := false

	for _, note := range notes {
		if note.ID == id {
			fmt.Println("Автор:")
			fmt.Println(note.Name, note.Surname)
			fmt.Println("Заметка:")
			fmt.Println(note.Text)

			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Заметка с ID %d не найдена.\n", id)
	}
}

func deleteNote() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите ID: ")
	idStr, _ := reader.ReadString('\n')
	var id int
	fmt.Sscanf(strings.TrimSpace(idStr), "%d", &id)

	found := false

	for i, note := range notes {
		if note.ID == id {
			notes = append(notes[:i], notes[i+1:]...)
			fmt.Printf("Заметка с ID %d удалена.\n", id)

			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Заметка с ID %d не найдена.\n", id)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/add":
		addNoteHTTP(w, r)
	case "/view":
		viewNoteHTTP(w, r)
	case "/delete":
		deleteNoteHTTP(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Страница не найдена")
	}
}

func addNoteHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Неподдерживаемый метод")
		return
	}
	reader := bufio.NewReader(r.Body)

	// Чтение данных из тела запроса
	name, _ := reader.ReadString('\n')
	name = name[:len(name)-1]

	surname, _ := reader.ReadString('\n')
	surname = surname[:len(surname)-1]

	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1]

	// Добавление заметки
	id := len(notes) + 1
	note := Note{ID: id, Name: name, Surname: surname, Text: text}
	notes = append(notes, note)

	// Ответ клиенту с ID новой заметки
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, strconv.Itoa(id))
}

func viewNoteHTTP(w http.ResponseWriter, r *http.Request) {
	// Обработка GET-запроса для просмотра заметки
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Неподдерживаемый метод")
		return
	}

	// Получение ID заметки из параметра запроса
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "ID заметки не указан или некорректный")
		return
	}

	var note Note
	found := false

	// Поиск заметки по ID
	for _, n := range notes {
		if n.ID == id {
			note = n
			found = true
			break
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Заметка не найдена")
		return
	}

	// Отправка данных заметки клиенту
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Автор: %s %s\nЗаметка: %s", note.Name, note.Surname, note.Text)
}

func deleteNoteHTTP(w http.ResponseWriter, r *http.Request) {
	// Обработка POST-запроса для удаления заметки
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Неподдерживаемый метод")
		return
	}

	// Получение ID заметки из параметра запроса
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "ID заметки не указан или некорректный")
		return
	}
	found := false
	// Поиск заметки по ID и удаление
	for i, note := range notes {
		if note.ID == id {
			notes = append(notes[:i], notes[i+1:]...)
			found = true
			break
		}
	}

	// Ответ клиенту о результате удаления
	if found {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Заметка с ID %d удалена", id)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Заметка с ID %d не найдена", id)

	}
}
