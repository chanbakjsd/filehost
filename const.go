package main

const (
	// The domain prepended to the result of /upload.
	domain = "http://localhost"
	// The host to listen to. Default is any address with port 80.
	host = ":80"
	// Try to keep total file size under 5GB.
	maxStorage = 5 * 1024 * 1024 * 1024
	// Wait 1 minute before checking if files need to be deleted.
	gcCooldown = 60
	// By default, GC also cleans up redirect URLs. Set this to `true` to make it skip them.
	skipRedirects = false

	// Password if you prefer to limit its access. Blank means not required.
	password = ""

	// Request limit of 1 request/12 sec/IP or 5 req/min/IP.
	requestPerSecond = 1.0 / 12
	// Size limit of 5MB/IP/min.
	sizePerSecond = 5 * 1024 * 1024 / 60
	// Maximum file size of 100MB at any point.
	burstSize = 100 * 1024 * 1024
	// Cache up to 1MB in memory per upload request.
	// Larger files are stored temporarily on disk.
	maxMemoryPerRequest = 1024 * 1024
)
