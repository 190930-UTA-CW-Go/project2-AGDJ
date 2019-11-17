# JAG-D
Is an application which can be used to manage your multiple clients with a single master server machine. The application is able to install a list of selected applications from the apt package on multiple client machines.
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
git clone https://github.com/190930-UTA-CW-Go/project2-AGDJ.git
```
Or download the setup.sh script and it will do everything for you like setting up go and correct pathing.
## Download necessary packages on your device
```go
go get golang.org/x/crypto/ssh
go get github.com/gorilla/mux
go get github.com/lib/pq
```
```bash
sudo apt install sysstat
```
## On Server(master)
you migth have do edit the values for your IPs in the database
```bash
#server needs docker to host database
sudo apt install docker.io -y
sudo usermod -aG docker $USER
sudo service docker start
cd go/src/github.com/190930-UTA-CW-Go/project2-AGDJ/servergo
sudo go run main.go &
#the app will then start on port :80
```
## On clients
```bash
cd go/src/github.com/190930-UTA-CW-Go/project2-AGDJ/clientgo
sudo go run main.go &
```
