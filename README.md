# Manga(nato) scraper API
As the name suggests, it's an api that you can use to fetch manga information as well as the panels.

## Routes
1. `/manga/get/{mangaID}` - Get manga information given mangaID.
2. `/manga/latest/` - Get the latest mangas.
3. `/chapter/panels/{mangaID}/{chapterID}` - Get panels for the chapter.
4. `/health` - Uhm, ping.

## Todos
1. Add functionality to fetch list of mangas.
2. Add proxy for panels.
3. ???

### Local
1. Build the project
```sh
make
```

2. Define environmental variables:
```env
HOST=...
PORT=...
```

3. Run the project
```sh
make start
```