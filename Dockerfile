FROM golang:1.18
WORKDIR /usr/local/go/src/kredit_plus
ADD . .
RUN export GO111MODULE=on
#GOPATH MUST BE OUTSIDE OF GOROOT directory!!!
RUN export GOPATH=/go/src/quickstart
RUN export PATH=$PATH:$GOPATH/bin

RUN export GOROOT=/usr/local/go/src/kredit_plus
RUN export PATH=$PATH:$GOROOT/bin
COPY . . 
RUN go get -d -v ./...
RUN go install -v ./...
EXPOSE 9333
COPY . . 
# Install server application
RUN go build -o main .
CMD ["go", "run", "main.go"]