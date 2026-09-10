# Go Uptime

A simple uptime checking system.

## Usage Overview

Create `conf.json` based on `README.conf.json`.

Create `.env` based on `README.env`.

Run a quick live test:

```bash
set -a; source .env; set +a; go run main.go --config=./conf.json
```

Run the tests:

```bash
go test
```

Build and run the container image using the Makefile, for example:

```bash
make build
make run
```

Deploying as an image to Google Cloud Artifact Registry:
```bash
sh ./deploy.sh --deploy
```