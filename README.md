# Tekton HTTP Resolver

Custom Tekton resolver written in Go. It fetches a Tekton Pipeline or Task resource from a configured HTTP endpoint and returns the response through Tekton's resolution framework.

## How it works

```text
Tekton PipelineRun
        |
        v
HTTP resolver
        |
        v
RESOLVER_SOURCE_URL
        |
        v
Tekton Pipeline or Task resource
```

The resolver registers itself with resolver type `http` and serves requests through Tekton's resolver framework. The source endpoint is configured at runtime; no external tunnel or personal development URL is embedded in the code.

## Prerequisites

- Go 1.21 or newer
- Tekton Pipelines installed in a Kubernetes cluster
- A reachable HTTP endpoint that returns a Tekton resource
- Docker, for building the container image

## Run locally

1. Copy `.env.example` to `.env` and set `RESOLVER_SOURCE_URL`.
2. Export the variable:

   ```bash
   export RESOLVER_SOURCE_URL=https://example.com/tekton/pipelines/tasks
   ```

3. Start the resolver:

   ```bash
   go run .
   ```

On Windows PowerShell:

```powershell
$env:RESOLVER_SOURCE_URL = "https://example.com/tekton/pipelines/tasks"
go run .
```

## Test and build

```bash
go test ./...
go build -o tekton-resolver .
```

## Build the container

```bash
docker build -t tekton-http-resolver:local .
```

The image uses a multi-stage build and runs as the compiled resolver binary.

## Deploy to Kubernetes

Update the image and `RESOLVER_SOURCE_URL` in `deployment.yaml`, then apply the manifests:

```bash
kubectl apply -f rbac.yaml
kubectl apply -f deployment.yaml
```

The example Tekton manifests demonstrate resolver references and event-triggered PipelineRuns:

- `test-request.yaml`
- `pipelinrun.yaml`
- `event-listener.yaml`
- `trigger-binding.yaml`
- `trigger-template.yaml`

## Configuration

| Variable | Required | Description |
| --- | --- | --- |
| `RESOLVER_SOURCE_URL` | Yes | HTTP endpoint returning the resource to resolve |

The endpoint must return a successful HTTP response containing the Tekton resource data.
