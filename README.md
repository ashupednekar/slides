# slides
https://slides.ashudev.in/
[Go present](https://pkg.go.dev/golang.org/x/tools/present) slide decks served over HTTP (default port **3999**).

## Decks

| Directory       | Topic                          |
|----------------|--------------------------------|
| `litefunctions/` | LiteFunctions / serverless talk |
| `litefunctions-go/` | LiteFunctions / Go architecture talk |
| `quickgrpc/`     | gRPC quick intro                 |

Each deck has a `.slide` file; supporting `.html` assets live alongside it.

## Run locally

Install `present` (Go 1.25+ recommended; matches the Docker image):

```bash
go install golang.org/x/tools/cmd/present@latest
```

From the repo root, open a deck directory and start the server. Use the bundled templates so the default **Thank you** closing slide is omitted (upstream `present` always appends it):

```bash
cd litefunctions   # or quickgrpc
present -base=../present-assets
```

Browse to the URL printed in the terminal (typically `http://127.0.0.1:3999`).

## Docker

Build and run:

```bash
docker build -t slides .
docker run --rm -p 3999:3999 slides
```

Published image: `ashupednekar535/slides` (see `Dockerfile`).

## Static export

Render every top-level `.slide` deck into a self-contained static directory:

```bash
(cd render && go run .)
```

This writes `dist/<deck>/index.html` for each deck, copies each deck's local
assets, and copies the `present-assets/static/` runtime into each output
directory. Static files larger than 25 MiB are uploaded to the `slides` R2
bucket when Wrangler and `SLIDES_R2_PUBLIC_BASE_URL` are available; otherwise
their references fall back to `placeholder.png` so the Pages deploy stays under
the per-file size limit. To render only selected decks:

```bash
(cd render && go run . litefunctions quickgrpc)
```

Deploy the rendered output to Cloudflare Pages:

```bash
export SLIDES_R2_PUBLIC_BASE_URL="https://your-public-r2-host"
./deploy.sh
```

## CD

On push to `main`, GitHub Actions builds the image, pushes to Docker Hub, then SSHs to the deploy host and runs `docker pull` + `docker restart slides`.

Configure in the repo **Settings → Secrets and variables → Actions**:

- **Secrets:** `DOCKER_PAT`, `ashupednekar` (SSH private key for deploy)
- **Variables:** `SSH_HOST` (deploy server address)
