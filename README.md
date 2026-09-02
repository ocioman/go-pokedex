# go-pokedex

## Italiano

CLI in Go che usa [PokeAPI](https://pokeapi.co/) per esplorare aree, incontrare Pokémon e gestire un piccolo Pokédex locale.

### Features & Technical Highlights

- Esplorazione delle `location areas` di Pokémon con navigazione tra pagine precedenti e successive.
- Cattura dei Pokémon con probabilità dinamica basata sull'esperienza base, con countdown animato e gestione degli errori.
- Bufferizzazione dei salvataggi: i Pokémon catturati vengono accumulati in `PokemonsBuffer` e salvati in modo batch da una goroutine dedicata (`WritePokemonsBufferLoop`). Quando i segnali sul canale raggiungono un certo threshold, il codice chiama `savePokedex()`, usa `bufio.NewWriter` per scrivere sul file `save.json` e poi svuota il buffer. Questo riduce le scritture continue e mantiene il loop interattivo reattivo.
- Persistenza locale del Pokédex su file `save.json` per mantenere i Pokémon catturati tra sessioni.
- Dettagli completi per ciascun Pokémon: statistiche, tipi, peso, altezza e sprite in ASCII convertite da immagini ufficiali.
- Cache HTTP di breve durata: `cache.Cache` usa un `RWMutex` e una goroutine `ReapLoop` con ticker per scadere automaticamente le entry dopo 6 minuti, evitando chiamate ripetute a PokeAPI e migliorando la velocità di risposta.
- Architettura CLI in Go con goroutine, mutex e canali di coordinamento per sincronizzare l'accesso concorrente ai dati senza corrompere lo stato del Pokédex.

### Requisiti

- Go 1.26+

### Avvio

```bash
go run .
```

Dopo l'avvio, comparirà il prompt:

```text
Pokedex >
```

### Comandi disponibili

- `help`: mostra l'elenco dei comandi
- `exit`: chiude l'applicazione
- `map`: mostra le prossime aree
- `mapb`: mostra le aree precedenti
- `explore <area>`: mostra i Pokémon incontrabili in un'area
- `catch <pokemon>`: prova a catturare un Pokémon
- `inspect <pokemon>`: mostra i dettagli di un Pokémon catturato
- `pokedex`: mostra tutti i Pokémon catturati

### Esempio rapido

```text
Pokedex > map
Pokedex > explore canalave-city-area
Pokedex > catch pikachu
Pokedex > inspect pikachu
Pokedex > pokedex
```

---

## English

Go CLI that uses [PokeAPI](https://pokeapi.co/) to explore areas, encounter Pokémon, and manage a small local Pokédex.

### Features & Technical Highlights

- Exploration of Pokémon `location areas` with pagination for previous and next regions.
- Catch system with dynamic success probability based on base experience, including animated throw timing and escape handling.
- Buffered save pipeline: captured Pokémon are accumulated in `PokemonsBuffer` and flushed in batches by a dedicated goroutine (`WritePokemonsBufferLoop`). When the channel receives enough events, the app calls `savePokedex()`, writes through `bufio.NewWriter` into `save.json`, and then clears the buffer. This reduces continuous disk writes and keeps the interactive loop responsive.
- Local Pokédex persistence with `save.json`, so captured Pokémon remain available across sessions.
- Rich Pokémon inspection with stats, types, weight, height, and ASCII sprite rendering generated from the official sprite images.
- Short-lived HTTP caching: `cache.Cache` uses an `RWMutex` plus a `ReapLoop` goroutine with a ticker to expire entries after 6 minutes, avoiding repeated PokeAPI requests and improving response time.
- Go-based CLI architecture using goroutines, mutexes, and channels to synchronize concurrent access to state without corrupting the Pokédex data.

### Requirements

- Go 1.26+

### Run

```bash
go run .
```

After startup, the prompt will appear:

```text
Pokedex >
```

### Available commands

- `help`: show the list of commands
- `exit`: close the application
- `map`: show the next areas
- `mapb`: show the previous areas
- `explore <area>`: show encounterable Pokémon in an area
- `catch <pokemon>`: try to catch a Pokémon
- `inspect <pokemon>`: show details for a caught Pokémon
- `pokedex`: show all caught Pokémon

### Quick example

```text
Pokedex > map
Pokedex > explore canalave-city-area
Pokedex > catch pikachu
Pokedex > inspect pikachu
Pokedex > pokedex
```
