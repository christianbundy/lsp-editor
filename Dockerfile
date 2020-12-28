FROM golang:latest
COPY main.go .
RUN GO111MODULE=on go get golang.org/x/tools/gopls@latest
RUN go build -o lsp-editor
CMD ./lsp-editor main.go
