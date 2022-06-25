# 🧑‍💻 Backend

The backend manages the server it self. It provides the API and the files for the clients.

## How to build

### Requirements

- Golang 1.18
- Software is only tested on Linux

### Build steps

To Build with a version replace `{{ version_tag }}` with your version

```sh
go build -ldflags="-X 'cloudgobrrr/backend/pkg/env.version={{ version_tag }}'" -o server main.go
```

Otherwise just remove the ldflag

Before running set the environment variables. See "Docs -> Getting started -> Configuration"

Then run it with

```sh
./server
```
