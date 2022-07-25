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

> This is the implementation of [this](https://www.udemy.com/course/backend-master-class-golang-postgresql-kubernetes/) course, except it's using MySQL instead of Postgres.
