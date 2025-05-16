# LiblaryApiDocs

# todo
1. przenieść z json.go do faktyczne jsona
2. dodać symulacje db - coś w stylu: tabelki
    > defaultowa jakaś ustawiona w jsonie
    > lua może ją sobie edytować
    > można tak robić jakieś małe przedstawienia działania
3. ustawianie portu jako run arg + defaultowy
4. dodać ikonke logowania, po kliknięciu pojawi się pole na podanie
api key. w zależności od configu strony będzie wymagane lub nie
do wyświetlania endpointów
5. funkcja testowania faktycznego api
    > przełącznik Sim <-> Api
    > pole na podanie Top Api Url / domain
    > wyświetlać res z faktycznego api
    > & można dodać np autouzupełnianie authKey żeby się nie wyświtlał
    jako value w Headers albo body, zamiast tego będzie np
    {{api_key}}

- może kiedyś zrobić platformę gdzie ludzie mogą robić swoją dokumentację czy coś??

* jeżeli pliki .md mają na początku nazwy "_"
to oznaczenie że są tylko jako załaczniki do ednpointów
i nie mają się wyświetlać pod oddzielnym przyciskiem

# docs

- ./markdowns
    + jeżeli chcemy stworzyć plik .md w /markdowns ale nie chcemy aby wyświetlał się on w kategoriach to w jego nazwie pierwszym znakiem musi być "_". pliki z nazwą zaczynającą się od "_" nie będą wyświetlane w kategoriach.

    + w ./markdowns/
        > umieszczamy folder, jego nazwa będzie używana jako nazwa kategori do której należy plik markdown
        > używamy maxymalnie 1 folder np: /markdowns/api_docs/some_api_doc.md

- ./scripts
    > są tu skrypty lua do obsługi symulowania backendu
    + nazwy plików oraz funkcji wewnątrz pliku będą używane do rozpoznawania jakiego pliku i funkcji chcemy użyć do obsługi symulacji danego endpointu
        > maxymalnie 1 funkcja na 1 plik

- ./docs
    > w pliku docs.json trzymamy listę jsonów na podstawie których generowany jest kontent strony.

    + dane w jsonie takie jak:
        -  markdown (lista kategori i plików .md) (opcjonalna)
        - category (kategoria pod którą będzie endpoint)
        - luaFunc (nazwa funkcji / pliku lua) (opcjonalna)
        - defaultDB (defaultowa pokazowa baza danych / tabelka) (opcjonalna)

przykład:
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
  }, 
]
```
