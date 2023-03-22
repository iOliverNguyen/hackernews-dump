# HackerNews Dump

This is a simple script to download the HackerNews data from [HackerNews API](https://github.com/HackerNews/API). It
downloads the data in JSON format and stores it as a collection of files. The first time it runs, it will create the
directory `data` to store the data. Subsequent runs will update the data in the directory.

At the time of this writing, the HackerNew database has about 35 million items.

## Usage

```bash
cd hackernews
direnv allow
go run . sync --concurrent=30 --update-back=1

# for debugging
go run . sync --debug
```

Install [direnv](https://direnv.net/) if necessary.

## How it works

The script will download the data from the [HackerNews API](https://github.com/HackerNews/API) and store it in the
`data` directory.

- It starts with [max item id](https://hacker-news.firebaseio.com/v0/maxitem.json) and works its way down to 1.
- It spans 30 goroutines to download the data in parallel.
- Data are split into chunks of 1000 items a file, and 1000 files a directory. Every 100 items download will be stored
  into a file as `<number>.partial.jsonl`, and when all 1000 items of the chunk are downloaded, the file will be named to
  `<number>.json`. It will skip the items that are already downloaded.
- Subsequent runs will pick up the data from `.partial.jsonl` file and continue updating the data. It also take the new max item id and continue downloading the data from there.
- Can use `--update-back` flag to update the data from the last N days, to have scores updated.
- Use `--debug` flag to limit the goroutine, the number of items per chunk, etc. for easily seeing progress and debugging.

## License

MIT
