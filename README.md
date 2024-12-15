
### Zadanie 1 : Skonteneryzowanie aplikacji

#### Kontekst:
Otrzymujesz bazową aplikację, która:
- Nie wykorzystuje bazy danych.
- Wszystkie dane są zapisywane w pliku `.json`.
- Pliki są przechowywane w lokalnym folderze aplikacji.


#### Twoje zadanie:
1. Napisz prosty plik `Dockerfile`, który umożliwi uruchomienie aplikacji w kontenerze.
2. Aplikacja w kontenerze powinna zachować pełną funkcjonalność, w tym:
   - Tworzenie użytkowników.
   - Logowanie użytkowników.
   - Uploadowanie plików do pliku lokalnego.

#### Wynik:
- Dołącz screen logów, które prezentują:
  1. Proces tworzenia użytkownika.
  2. Logowanie użytkownika.
  3. Upload plików.


### Zadanie 2 : Postawienie kontenera z bazą danych

1. Postaw kontener bazy danych korzystając z gotowego obrazu dostępnego w Docker Hub (np. `postgres`, `mysql` lub innego).
2. Skonfiguruj bazę tak, aby aplikacja mogła się z nią połączyć.
3. Wykonaj w kontenerze bazy danych proste zapytanie SQL, np. `CREATE TABLE users`.

#### Wynik:
- Dołącz screen logów potwierdzających:
  1. Że baza danych działa (logi uruchomienia kontenera).
  2. Wykonanie zapytania SQL `CREATE TABLE users` w uruchomionej bazie danych.


### Zadanie 3 : Połączenie aplikacji i bazy danych w kontenerach

#### Twoje zadanie:
1. Stwórz plik `docker-compose.yaml`, w którym:
   - Poprawnie zdefiniujesz oba kontenery (aplikacja i baza danych).
   - Ustalisz sposób uruchamiania tych kontenerów.
   - Zdefiniujesz wolumen, który umożliwi przechowywanie plików użytkowników w udostępnionym katalogu.
2. Upewnij się, że:
   - Dane użytkowników są zapisywane w bazie danych.
   - Pliki użytkowników są przechowywane w odpowiednich katalogach w udostępnionym wolumenie.

#### Wynik:
- Dołącz screeny pokazujące:
  1. Że dane użytkowników są poprawnie zapisywane w bazie danych (np. wynik zapytania SQL `SELECT * FROM users`).
  2. Że pliki użytkowników są przechowywane w odpowiednich katalogach w udostępnionym wolumenie.




NOTES: -------------------------------------------------------------

    Show all containers:
    docker ps -a

    Acces db container:
    docker exec -it DATABASE mysql -u admin -ppassword root_password

    List of volumes:
    docker volume ls

    Check the volume data:
    docker volume inspect go-simple-authorization_shared-data
