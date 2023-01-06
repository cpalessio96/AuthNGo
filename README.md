# AuthNGo

App per autenticazione tramite jwt

## Build & Deploy
Per buildare e deployare l'app basta eseguire i comandi:
``` docker compose build ```
``` docker compose up -d ```
L'applicazione sta in ascolto nella porta 8080.
Il build dell'applicazione crea anche una instanza postgres

## Configurazione db
TODO: inserire qui i passaggi per la creazione del database e delle tabelle

## API
Le api create sono le seguenti:
### Registration
Questa api ha path /registration e metodo in POST, serve per la registazione dell'utente e accetta nome, cognome, email e password dell'utente.
Di seguito ecco esempi di payload:
input
```json
{
    "email": "catania.alessio96@gmail.com",
    "password": "12345rf6",
    "firstName": "alessio",
    "lastName": "catania"
}
```
output
```json
{
    "error": false,
    "message": "Logged in user 12",
    "data": {
        "firstname": "alessio",
        "lastname": "catania",
        "email": "catania.alessio96@gmail.com"
    }
}
```