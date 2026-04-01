# slides

[Go present](https://pkg.go.dev/golang.org/x/tools/present) slide decks served over HTTP (default port **3999**).

## Decks

| Directory       | Topic                          |
|----------------|--------------------------------|
| `litefunctions/` | LiteFunctions / serverless talk |
| `quickgrpc/`     | gRPC quick intro                 |

Each deck has a `.slide` file; supporting `.html` assets live alongside it.

## Run locally

Install `present` (Go 1.25+ recommended; matches the Docker image):

```bash
go install golang.org/x/tools/cmd/present@latest
```

From the repo root, open a deck directory and start the server:

```bash
cd litefunctions   # or quickgrpc
present
```

Browse to the URL printed in the terminal (typically `http://127.0.0.1:3999`).

## Docker

Build and run:

```bash
docker build -t slides .
docker run --rm -p 3999:3999 slides
```

Published image: `ashupednekar535/slides` (see `Dockerfile`).

## CD

On push to `main`, GitHub Actions builds the image, pushes to Docker Hub, then SSHs to the deploy host and runs `docker pull` + `docker restart slides`.

Configure in the repo **Settings → Secrets and variables → Actions**:

- **Secrets:** `DOCKER_PAT`, `ashupednekar` (SSH private key for deploy)
- **Variables:** `SSH_HOST` (deploy server address)
