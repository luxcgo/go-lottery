# go-gallery

```bash
docker pull postgres:9.6.21-alpine
docker run --name postgres9 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:9.6.21-alpine
docker exec -it postgres9 psql -U root
```

