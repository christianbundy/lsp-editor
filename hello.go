package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type Response struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	//  Result int `json:"result"`
}

type RequestMessage struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`

	Method string `json:"method"`
	Params string `json:"params"`
}

type InitializeRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`

	Method string           `json:"method"`
	Params InitializeParams `json:"params"`
}

type ResponseMessage struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
}

type InitializedRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`

	Method string            `json:"method"`
	Params InitializedParams `json:"params"`
}

type WorkspaceFolder struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
}

type InitializedParams struct {
}

type InitializeParams struct {
	RootUri               string                `json:"rootUri"`
	WorkspaceFolders      []WorkspaceFolder     `json:"workspaceFolders"`
	InitializationOptions InitializationOptions `json:"initializationOptions"`
	Capabilities          ClientCapabilities    `json:"capabilities"`
}

type ClientCapabilities struct {
	TextDocument TextDocumentClientCapabilities `json:"textDocument"`
}

type TextDocumentClientCapabilities struct {
	SemanticTokens SemanticTokensClientCapabilities `json:"semanticTokens"`
}

type SemanticTokensClientCapabilities struct {
	TokenTypes []string `json:"tokenTypes"`
}

type InitializationOptions struct {
	SemanticTokens bool `json:"semanticTokens"`
}

type TextDocumentIdentifier struct {
	Uri string `json:"uri"`
}

type SemanticTokensParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}
type SemanticTokensRequest struct {
	Jsonrpc string `json:"jsonrcp"`
	Id      int    `json:"id"`

	Method string               `json:"method"`
	Params SemanticTokensParams `json:"params"`
}

func main() {
	cmd := exec.Command("gopls")

	stdin, err := cmd.StdinPipe()

	if err != nil {
		log.Fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(stdout)

	err = cmd.Start()

	if err != nil {
		log.Fatal(err)
	}

	directory := "file:///home/christianbundy/src/lsp-editor/"
	file := directory + "hello.go"

	send := func(input interface{}) bool {
		request := input
		b, err := json.Marshal(request)
		if err != nil {
			log.Fatal(err)
		}

		size := len(b)
		out := "Content-Length: " + strconv.Itoa(size) + "\r\n\r\n" + string(b)
		// fmt.Println("=== <send>")
		// fmt.Println(string(out))

		_, err = stdin.Write([]byte(out))
		// fmt.Println("=== </send>")

		if err != nil {
			panic(err)
		}
		return true
	}

	receive := func(target int) []byte {
		var result []byte

		for true {
			// fmt.Println("receiving")

			data, err := reader.ReadString('\n')
			// fmt.Println("=== <receive.header>")
			// fmt.Println(string(data))

			// fmt.Println("=== </receive.header>")
			if err != nil {
				fmt.Println("invalid header")
				continue
			}
			prefix := "Content-Length: "
			if strings.HasPrefix(string(data), prefix) == false {
				fmt.Println("expecting header")
				continue
			}

			content_length, _ := strconv.Atoi(strings.TrimSpace(string(data)[len(prefix):]))

			// + 2 because newlines
			buf := make([]byte, content_length+2)
			if _, err := io.ReadFull(reader, buf); err != nil {
				log.Fatal(err)
			}

			decoded := ResponseMessage{}
			err = json.Unmarshal(buf, &decoded)
			if err != nil {
				fmt.Println("invalid json")
				continue
			}
			if decoded.Id == target {
				// fmt.Println("=== <receive.data>")
				// fmt.Println(string(buf))
				// fmt.Println("=== </receive.data>")
				result = buf
				break
			}
		}
		return result
	}
	send(InitializeRequest{Jsonrpc: "2.0",
		Id:     1,
		Method: "initialize",
		Params: InitializeParams{
			RootUri: file,
			WorkspaceFolders: []WorkspaceFolder{
				WorkspaceFolder{
					Name: "lsp-editor",
					Uri:  directory,
				},
			},
			InitializationOptions: InitializationOptions{
				SemanticTokens: true,
			},
			Capabilities: ClientCapabilities{
				TextDocument: TextDocumentClientCapabilities{
					SemanticTokens: SemanticTokensClientCapabilities{
						TokenTypes: []string{
							"namespace",
							"type",
							"class",
							"enum",
							"interface",
							"struct",
							"typeParameter",
							"parameter",
							"variable",
							"property",
							"enumMember",
							"event",
							"function",
							"method",
							"macro",
							"keyword",
							"modifier",
							"comment",
							"string",
							"number",
							"regexp",
							"operator",
						},
					},
				},
			},
		},
	})

	send(InitializedRequest{Jsonrpc: "2.0",
		Id:     2,
		Method: "initialized",
		Params: InitializedParams{},
	})

	send(SemanticTokensRequest{
		Id:      3,
		Jsonrpc: "2.0",
		Method:  "textDocument/semanticTokens/full",
		Params: SemanticTokensParams{
			TextDocument: TextDocumentIdentifier{
				Uri: file,
			},
		},
	})
	fmt.Println(string(receive(3)))
}
