# env-setup-local-container/ — invariants

## Binary must be compiled inside Docker
`env-setup-local-container/docker-compose.yml` uses `build: { dockerfile: env-setup-local-container/Dockerfile.infra }` — not `image: newrelic/infrastructure:latest` with a volume-mounted binary.
Do not revert to the image+volume pattern. A macOS-compiled binary mounted into a Linux container produces `permission denied` at exec time regardless of chmod.

## /host/proc errors in newrelic-infra logs are harmless
On macOS+Docker Desktop the infra agent cannot see the host's /proc or /sys. This only breaks host samplers (CPU, disk, network). The HAProxy integration runs independently and is unaffected.
