# project2-AGDJ

# Running DB
```bash
docker build -t project2 .
docker run -p 5432:5432 -d --rm --name runningproject2 project2
```

**OPTIONAL COMMAND**

If you want to look into your table in postgres use this command
```bash
docker exec -it runningproject2 psql -U postgres
```
