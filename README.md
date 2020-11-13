# Simple File Hosting

This is designed to be a simple Go server that just hosts files and does URL shortening.
Use at your own risk.

## Usage

Modify the values in `const.go` and `go build .` it.
Consider replacing `[this_domain]` in `static/index.html` with your domain as well.

You are recommended to keep this behind a reverse proxy that adds HTTPS.

## Defaults

- The server listens to `http://localhost:80`.
- Files are deleted once per minute to keep it under 5GB.
- No authentication required.
- Maximum file size of 100MB.
- 5 requests/IP/min.
- 5 MB/IP/min.

## Credits

Inspired by [0x0](https://github.com/mia-0/0x0) and [Better Motherf**king Website](http://bettermotherfuckingwebsite.com).

## License

This project is licensed under MIT.
