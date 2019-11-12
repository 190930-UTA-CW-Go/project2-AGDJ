# project2-AGDJ

## Download packages
```go
go get golang.org/x/crypto/ssh
go get github.com/gorilla/mux
```

```bash
sudo apt install sysstat
```

```bash
To create ssh key:
ssh-keygen -t rsa
ssh-copy-id agent1@192.168.56.102
```

```bash
sudo -S, --stdin
Write the prompt to the standard error and read the password from the standard input instead of using the terminal device. The password must be followed by a newline character.
```

## Running DB
```bash
docker build -t project2 .
docker run -p 5432:5432 -d --rm --name runningproject2 project2
```

**OPTIONAL COMMAND**

If you want to look into your table in postgres use this command
```bash
docker exec -it runningproject2 psql -U postgres
```