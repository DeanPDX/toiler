# Toiler

This is a task tracking app, written with Preact, Go, and Postgres. This version is very basic and doesn't even properly support auth yet so do not use. I'm using this as a holding tank for now to build out initial app.

To start API, run:

```
go run github.com/DeanPDX/toiler/api
```

## Testing docker build
Build a docker image with the tag toiler:

```bash
docker build . --tag toiler
```

Then run the image, expose a port and pass in environment variables:

```bash
# This assumes you have a postgres instance running on a different docker image:
docker run -p 8090:8090 --env DSN='postgres://todoapp:LetMeIn!@yourVEthernetAdapterIP:5432/todoapp' --env PORT=8090 --env SIGNINGSECRET=SomeCoolJWTSigningSecret toiler
```