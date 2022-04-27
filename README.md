# Simplebank

To run on local env:
Run mysql in the docker

    docker-compose up -d

Then run the app

    go mod tidy
    go install github.com/codegangsta/gin@latest
    make autoreload

If you get `gin: command not found` error, add this to `.profile` or `.bashrc` or `.zshrc`

    export GOPATH=$HOME/go
    export PATH=$GOPATH/bin:$PATH
