# Environment Variables

cmn-core is configured primarily via `etc/app.yaml`. There are no mandatory environment variables for the application itself.

## Docker Compose

For the dev environment, set platform in `.env.local` (not committed):

```bash
# Apple Silicon
echo "DOCKER_PLATFORM=linux/arm64" > .env.local

# Intel / Linux x86_64
echo "DOCKER_PLATFORM=linux/amd64" > .env.local
```
