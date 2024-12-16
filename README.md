<!-- # Instrukcja instalacji

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

Jak odpalić apke do zadań:
Zainstaluj zależności Go:
`go mod download`

W folderze z zadaniami zbuildować apke
`go build -o ./tmp/main.exe .`.

Odpalić ją po zbuildowaniu
`./tmp/main.exe`. -->

# Instalacja środowiska Docker

> Tą część oczywiście można pominąć jeśli już masz zainstalowanego Dockera

- Windows: https://docs.docker.com/desktop/setup/install/windows-install/
- Linux:
  - Sam Docker Engine: https://docs.docker.com/engine/install/
  - Docker Desktop: https://docs.docker.com/desktop/setup/install/linux/
- Mac: https://docs.docker.com/desktop/setup/install/mac-install/

Upewnij się że środowisko działa poprawnie np. przez `docker --version`

# Zadanie 1 : Skonteneryzowanie prostego serwisu

## a) Kontekst:

Otrzymujesz prosty http web serverm który:

- Do komunikacji używa portu 8080
- Implementuje prosty system uwierzytelnienia - wszystkie dane są zapisywane w pliku `userData.json`.
- Umożliwia użytkownikom upload ich personalnych plików - w obecnie mocno okrojonej wersji wszystkie pliki trafiają do folderu `/uploads`.

## b) Twoje zadanie:

1. Napisz prosty plik `Dockerfile`, który umożliwi uruchomienie aplikacji w kontenerze.
2. Aplikacja w kontenerze powinna zachować pełną funkcjonalność, w tym:
   - Tworzenie użytkowników.
   - Logowanie użytkowników.
   - Uploadowanie plików do lokalnego folderu.
3. Postaraj się aby twój `Dockerfile` tworzył użytkownika _nieuprzywilejowanego_ w celu zapewnienia większego bezpieczeństwa naszego serwisu.

## c) Wynik:

Najlepiej screen Terminala/Powershella z widoczną:

1. Komendą która uruchamia twój kontener
   Przykładowo: `docker run --name zadanie_1_kontener -p 8080:8080 zadanie_1_image`
2. Logami aplikacji w których widać

   - Utorzenie nowego użytkownika
   - Zalogowanie się tego użytkownika
   - Upload jakiegoś pliku

   Logi powinny wyglądać mniej więcej tak:
   ![alt text](images/zad1_2.png)

### Uruchomiona aplikacja powinna wyglądać następująco:

![alt text](images/zad1_0.png)

# Zadanie 2: Postawienie kontenera z bazą danych

## a) Kontekst:

Musisz przygotować rozszerzenie serwisu:

- Aplikacja ma łączyć się z niezależną bazą danych
- Zamiast w pliku `userData.json` dane użytkownika będą trafiały właśnie do tej bazy

1. Postaw kontener bazy danych korzystając z gotowego obrazu dostępnego w Docker Hub (np. `postgres`, `mysql` lub innego). Za pomocą `docker-compose.yaml`.
2. Skonfiguruj bazę tak, aby aplikacja mogła się z nią połączyć.
3. Wykonaj w kontenerze bazy danych proste zapytanie SQL, np. `CREATE TABLE users`.

## Wynik:

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
