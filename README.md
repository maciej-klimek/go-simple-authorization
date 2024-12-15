
# Instrukcja instalacji

## 1. Instalacja Go
# Instalacja Go na Windows

## Krok 1: Pobierz instalator Go
1. Przejdź na oficjalną stronę Go: [https://golang.org/dl/](https://golang.org/dl/).
2. Pobierz instalator `.msi` odpowiedni dla Twojej wersji systemu (zazwyczaj 64-bit).

## Krok 2: Uruchom instalator
1. Znajdź pobrany plik instalatora `.msi` w folderze `Pobrane` lub innym domyślnym folderze.
2. Kliknij dwukrotnie na plik, aby uruchomić instalator.
3. Postępuj zgodnie z instrukcjami wyświetlanymi na ekranie:
   - Akceptuj warunki licencji.
   - Wybierz domyślną lokalizację instalacji (zazwyczaj `C:\Program Files\Go`).
   - Kliknij przycisk **Install**, aby rozpocząć instalację.

## Krok 3: Dodaj Go do zmiennych środowiskowych (opcjonalne)
1. Instalator automatycznie dodaje Go do zmiennych środowiskowych. Możesz to sprawdzić:
   - Kliknij prawym przyciskiem na **Ten komputer** lub **Mój komputer**.
   - Wybierz **Właściwości** → **Zaawansowane ustawienia systemu** → **Zmienne środowiskowe**.
   - Sprawdź, czy w zmiennej `PATH` znajduje się ścieżka do katalogu Go, np. `C:\Program Files\Go\bin`.

## Krok 4: Sprawdź instalację
1. Otwórz terminal (np. PowerShell lub Wiersz polecenia).
2. Wpisz poniższą komendę, aby sprawdzić zainstalowaną wersję Go:
   ```bash
   go version


### Linux
1. Pobierz archiwum `.tar.gz` odpowiednie dla Twojej wersji systemu z [https://golang.org/dl/](https://golang.org/dl/).
2. Rozpakuj archiwum do `/usr/local`:
   ```bash
   sudo tar -C /usr/local -xzf go<wersja>.linux-amd64.tar.gz


 Instalacja Dockera
Pobierz Docker z oficjalnej strony: https://www.docker.com/products/docker-desktop
Zainstaluj Docker, postępując zgodnie z instrukcjami dla Twojego systemu operacyjnego.
Upewnij się, że Docker działa poprawnie, uruchamiając w terminalu:
`docker --version`



Zadania znajdują się w folderze `zadania/baza`

Jak odpalić apke do zadań:
Zainstaluj zależności Go:
`go mod download`

W folderze z zadaniami zbuildować apke
`go build -o ./tmp/main.exe .`.

Odpalić ją po zbuildowaniu 
`./tmp/main.exe`.


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

1. Postaw kontener bazy danych korzystając z gotowego obrazu dostępnego w Docker Hub (np. `postgres`, `mysql` lub innego). Za pomocą `docker-compose.yaml`.
2. Skonfiguruj bazę tak, aby aplikacja mogła się z nią połączyć.
3. Wykonaj w kontenerze bazy danych proste zapytanie SQL, np. `CREATE TABLE users`.

#### Wynik:
- Dołącz screen logów potwierdzających:
  1. Że baza danych działa (logi uruchomienia kontenera).
  2. Wykonanie zapytania SQL `CREATE TABLE users` w uruchomionej bazie danych.


### Zadanie 3 : Połączenie aplikacji i bazy danych w kontenerach

#### Twoje zadanie:
1. Przerób plik `docker-compose.yaml`, w którym:
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
