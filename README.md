# JAG-D
Is an application which can be used to manage your multiple clients with a single point master server machine. The application is able to install a list of selected applications from the apt package on multiple client machines.
##
Download Go if not already installed. Here is an easy setup for Ubuntu machines.
```bash
sudo apt update
sudo apt upgrade -y
sudo apt install golang-go -y
echo "export PATH=$PATH:/usr/lib/go/bin" >> .profile
source .profile
echo "export GOPATH=/home/<$USERNAME>/go" >> .bashrc
source .bashrc
cd $HOME
mkdir go
cd go/
mkdir src/
cd src/
mkdir github.com/
cd github.com
mkdir 190930-UTA-CW-Go
cd 190930-UTA-CW-Go
git clone 
```
## Download packages
```go
go get golang.org/x/crypto/ssh
go get github.com/gorilla/mux
go get github.com/lib/pq
```

```bash
sudo apt install sysstat
```
