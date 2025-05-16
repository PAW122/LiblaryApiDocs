# 📚 Dokumentacja projektu

---

🔗 [Jak używać systemu dokumentacji](https://github.com/PAW122/LiblaryApiDocs)

---

## 📁 ./markdowns

- Jeśli chcesz stworzyć plik `.md` w `/markdowns`, ale **nie chcesz**, aby był wyświetlany w kategoriach — **nazwa pliku musi zaczynać się od znaku `_`**.

- W katalogu `./markdowns/`:
  - Umieszczamy **foldery**, których nazwy są używane jako **nazwy kategorii** dla plików Markdown.
  - **Obsługiwane jest tylko jedno zagłębienie**: np. `/markdowns/api_docs/some_api_doc.md`.

---

## 📁 ./scripts

- Zawiera skrypty Lua do obsługi symulowania backendu.
- Nazwa pliku oraz funkcji w środku służą do rozpoznawania symulacji dla konkretnego endpointu.
  - **Maksymalnie jedna funkcja na jeden plik**.

---

## 📁 ./docs

- Plik `docs.json` zawiera listę JSON-ów, na podstawie których generowany jest kontent strony.

### 📌 Dostępne pola w `docs.json`:

| Pole         | Typ        | Opis                                                                 |
|--------------|------------|----------------------------------------------------------------------|
| `markdown`   | `[]string` | (opcjonalne) Ścieżki do plików `.md` przypisanych do tego endpointu |
| `method`     | `string`   | Metoda HTTP (np. GET, POST)                                         |
| `endpoint`   | `string`   | Endpoint API                                                        |
| `description`| `string`   | Opis działania endpointu                                            |
| `permissions`| `string`   | Wymagane uprawnienia (np. `admin`)                                  |
| `body`       | `string`   | Przykładowe ciało zapytania                                         |
| `headers`    | `string`   | Wymagane nagłówki HTTP                                              |
| `res`        | `string`   | Przykładowa odpowiedź                                               |
| `errors`     | `[]object` | Lista możliwych błędów                                              |
| `category`   | `string`   | Kategoria, pod którą będzie pogrupowany endpoint                    |
| `luaFunc`    | `string`   | (opcjonalne) Nazwa funkcji Lua do symulacji                         |
| `defaultDB`  | `[]object` | (opcjonalne) Domyślna struktura pokazowa bazy danych (tabelka)      |

---

## 📦 Przykład wpisu `docs.json`

```json
[
  {
    "markdown": ["api/_api_books_add.md"],
    "method": "POST",
    "endpoint": "/api/book/add",
    "description": "dodaj nową książkę do db",
    "permissions": "admin",
    "body": "{\"title\": \"example_title\", \"author\": \"ex_author\", \"publisher\": \"ex_publisher\", \"year\": \"2025\", \"isbn\": \"ex_isbn\", \"pages\": \"100\", \"description\": \"ex_description\", \"tags\": \"ex_tags\", \"shelf\": \"4C\", \"copies\": \"{copies json}\"}",
    "headers": "{\"Content-Type\":\"application/json\", \"auth_token\": \"user_token\", \"username\": \"example_username\"}",
    "res": "200 OK",
    "errors": [
      {
        "code": "401",
        "message": "Brak autoryzacji",
        "description": "użytkownik nie jest zalogowany lub nie ma uprawnień"
      },
      {
        "code": "500",
        "message": "Internal Server Error",
        "description": "błąd serwera"
      }
    ],
    "category": "LiblaryManagerApi",
    "luaFunc": "api_books_add",
    "defaultDB": [
      {
        "title": "example_title",
        "author": "author",
        "publisher": "publisher",
        "year": "year",
        "isbn": "isbn",
        "pages": "pages",
        "description": "description",
        "tags": "tags",
        "shelf": "shelf",
        "copies": "copies_data {json} ..."
      }
    ]
  }
]
