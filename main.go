package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIError struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

type APIDoc struct {
	Method      string     `json:"method"`
	Endpoint    string     `json:"endpoint"`
	Description string     `json:"description"`
	Permissions string     `json:"permissions"`
	Body        string     `json:"body"`
	Response    string     `json:"res"`
	Errors      []APIError `json:"errors"`
	Category    string     `json:"category"`
}

var docs = []APIDoc{
	{
		Method:      "",
		Endpoint:    "",
		Description: "",
		Permissions: "",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "other",
	},
	{
		Method:      "GET",
		Endpoint:    "/api/books/get/<id>",
		Description: "Pobierz książkę o podanym id",
		Permissions: "user",
		Body:        `-`,
		Response:    `200 OK | body: (dane książki)`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "-"},
			{"400", "message: Brak poprawnego numeru inwentarzowego w ścieżce", "niepoprawne id"},
			{"404", "message: Książka o takim numerze nie istnieje", "brak id książki w db"},
			{"500", "Wystąpił błąd podczas sprawdzania książki", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "GET",
		Endpoint:    "/api/books/full/<id>",
		Description: "zwraca całego jsona z danymi o konkretnej książce",
		Permissions: "user",
		Body:        `-`,
		Response:    `200 OK | body: (dane książki)`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"400", "Nieprawidłowy numer inwentarzowy", "-"},
			{"404", "Nie znaleziono książki", "-"},
			{"404", "Brak danych książki", "-"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{ // niepełne errors od tego miejsca
		Method:      "POST",
		Endpoint:    "/api/book/add",
		Description: "dodaj nową książkę do db",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/book/borrow",
		Description: "wypożyczenie książki przez czytalnika w bibliotece",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "GET",
		Endpoint:    "/api/book/list",
		Description: "Lista wszystkich książek w bazie danych / bibliotece",
		Permissions: "user",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/book/update",
		Description: "aktualizacja danych książki - np poprzez edytowanie w formularzu",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/book/delete",
		Description: "usunięcie książki z db",
		Permissions: "",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/save_rules/borrow_time",
		Description: "zapisanie zasady maxymalnego czasu wypożyczenia książki",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "GET",
		Endpoint:    "/api/read_rules/borrow_time",
		Description: "odczytanie zasady maxymalnego czasu wypożyczenia książki",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/save_rules/overdue",
		Description: "zapisanie zasad dotyczących przetrzymania książki / opłat",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "GET",
		Endpoint:    "/api/load_rules/overdue",
		Description: "odczytanie zasad dotyczących przetrzymania książki / opłat",
		Permissions: "user",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/book/return",
		Description: "zwrócenie książki przez czytelnika",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/book/reserve",
		Description: "rezerwacja książki przez czytelnika",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/book/confirm_pickup",
		Description: "potwierdzenie odbioru wcześniej zarezerwowanej książki",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/book/save_cover",
		Description: "zapisanie okładki książki",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/book/get_cover",
		Description: "pobranie okładki książki",
		Permissions: "user",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/book/delete_cover",
		Description: "usunięcie okładki książki",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/book/upload_pdf_book",
		Description: "przesłąnie pliku pdf książki",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/book/delete_pdf_book",
		Description: "usunięcie pliku pdf książki",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/book/read_page",
		Description: "odczytanie konkretnej strony książki",
		Permissions: "user",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/stats/user",
		Description: "pobranie preferencji kategori książek czytanych przez użytkownika do rekomendacji",
		Permissions: "user",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/stats/time_piriod",
		Description: "odczytanie przeferencji w okresie czasu od month1 - month2",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/stats/reservations",
		Description: "statystyki rezerwacji książek w okresie czasu od month1 - month2\n informuje które książki najczęściej brakuje na pułkach",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	{
		Method:      "GET",
		Endpoint:    "/api/payments/read",
		Description: "odczytanie statysytk płatności",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryManagerApi",
	},
	// ==== następna kategoria ======================================
	{
		Method:      "POST",
		Endpoint:    "/api/user/get_acc_info",
		Description: "pobranie podstawowych informacji o użytkowniku",
		Permissions: "user",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryUserApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/user/get_borrowed_books",
		Description: "pobranie listy książek wypożyczonych przez użytkownika",
		Permissions: "user",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryUserApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/user/get_borrow_history",
		Description: "pobranie historii wypożyczeń książek przez użytkownika",
		Permissions: "user",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryUserApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/user/get_notifications",
		Description: "pobranie powiadomień dla użytkownika",
		Permissions: "user",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryUserApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/user/mark_notification_read",
		Description: "oznaczenie powiadomienia jako przeczytane",
		Permissions: "user",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryUserApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/user/reserve",
		Description: "rezerwacja książki na stronie(zdalnie) przez użytkownika",
		Permissions: "user",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryUserApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/accounts/update/",
		Description: "aktualizacja danych konta użytkownika",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryAccountsApi",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/accounts/register",
		Description: "rejestracja nowego konta użytkownika",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryAccountsApi",
	},
	{
		Method:      "GET",
		Endpoint:    "/api/accounts/getList",
		Description: "pobranie listy kont użytkowników",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryAccountsApi",
	},
	{
		Method:      "GET",
		Endpoint:    "/api/accounts/read",
		Description: "pobieranie danych konta użytkownika na podstawie username",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryAccountsApi",
	},
	{
		Method:      "GET",
		Endpoint:    "/api/accounts/details",
		Description: "pobieranie szczegółowych danych konta użytkownika",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryAccountsApi",
	},
	{
		Method:      "GET",
		Endpoint:    "/api/accounts/borrowed",
		Description: "Pobiera listę wypożyczonych książek przez użytkownika na podstawie indeksu",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryAccountsApi",
	},
	{
		Method:      "GET",
		Endpoint:    "/api/accounts/notifications",
		Description: "Pobiera listę powiadomień użytkownika na podstawie indeksu",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryAccountsApi",
	},
	{
		Method:      "GET",
		Endpoint:    "/api/accounts/preferences",
		Description: "Pobiera preferencje użytkownika na podstawie indeksu",
		Permissions: "admin",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "LiblaryAccountsApi",
	},
	{
		Method:      "GET",
		Endpoint:    "/",
		Description: "wczytywanie stron / StaticFileHandler",
		Permissions: "brak",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "other",
	},
	{
		Method:      "GET",
		Endpoint:    "/api/hello",
		Description: "testowy endpoint",
		Permissions: "brak",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "other",
	},
	{
		Method:      "POST?",
		Endpoint:    "/api/login",
		Description: "logowanie użytkownika",
		Permissions: "brak",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "other",
	},
	{
		Method:      "POST",
		Endpoint:    "/api/register",
		Description: "rejestracja użytkonika poprzez stronę testową - każdy może tworzyc konto admina\n call CreateAccountHandler()",
		Permissions: "brak",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "other",
	},
	{
		Method:      "",
		Endpoint:    "/",
		Description: "strona testowa do tworzenia kont usera i admina",
		Permissions: "",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"401", "Brak autoryzacji", "użytkownik nie jest zalogowany lub nie ma uprawnień"},
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "pages",
	},
	{
		Method:      "",
		Endpoint:    "/manager",
		Description: "strona admina do zarządzania książkami i użytkownikami",
		Permissions: "",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "pages",
	},
	{
		Method:      "",
		Endpoint:    "/liblary",
		Description: "strona gł dla użytkownika do przeglądania książek i innych akcji",
		Permissions: "",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "pages",
	},
	{
		Method:      "",
		Endpoint:    "/favicon.ico",
		Description: "favicon",
		Permissions: "",
		Body:        `-`,
		Response:    `200 OK`,
		Errors: []APIError{
			{"500", "Internal Server Error", "błąd serwera"},
		},
		Category: "pages",
	},
}

func main() {
	fmt.Println("server started on port 5899")
	http.HandleFunc("/api/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(docs)
	})

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.ListenAndServe(":5899", nil)
}
