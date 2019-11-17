#! /bin/bash
sudo apt update
sudo apt upgrade -y
sudo apt install golang-go -y
echo "export PATH=$PATH:/usr/lib/go/bin" >> .profile
source .profile
echo "export GOPATH=/home/$USERNAME/go" >> .bashrc
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

