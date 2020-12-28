package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Response struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
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

type SemanticTokensResponse struct {
	Jsonrpc string               `json:"jsonrpc"`
	Id      int                  `json:"id"`
	Result  SemanticTokensResult `json:"result"`
}

type SemanticTokensResult struct {
	ResultId string `json:"resultId"`
	Data     []int  `json:"data"`
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
	file, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	directory := filepath.Dir(file)

  var cmd *exec.Cmd

  if strings.HasSuffix(file, "Dockerfile") {
	  cmd = exec.Command("docker-langserver", "--stdio")
  }

  if strings.HasSuffix(file, ".go") {
	  cmd = exec.Command("gopls")
  }

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

	directoryUri := "file://" + directory
	fileUri := "file://" + file

	send := func(input interface{}) bool {
		request := input
		b, err := json.Marshal(request)
		if err != nil {
			log.Fatal(err)
		}

		size := len(b)
		out := "Content-Length: " + strconv.Itoa(size) + "\r\n\r\n" + string(b)

		_, err = stdin.Write([]byte(out))

		if err != nil {
			panic(err)
		}
		return true
	}

	receive := func(target int) []byte {
		var result []byte

		for true {

			data, err := reader.ReadString('\n')

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
			RootUri: fileUri,
			WorkspaceFolders: []WorkspaceFolder{
				WorkspaceFolder{
					Name: "lsp-editor",
					Uri:  directoryUri,
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
				Uri: fileUri,
			},
		},
	})

	out := receive(3)
	tokens := SemanticTokensResponse{}
	err = json.Unmarshal(out, &tokens)
	if err != nil {
		log.Fatal(err)
	}
	d := tokens.Result.Data

	f, err := os.Open(file)

	if err != nil {
		log.Fatal(err)
	}

	fileReader := bufio.NewReader(f)

	if err != nil {
		log.Fatal(err)
	}

	cursor := 0

	for i := 0; i < len(d); i += 5 {

		if d[i] > 0 {
			for skipLines := 0; skipLines < d[i]; skipLines++ {
				line, _, err := fileReader.ReadLine()
				cursor = 0
				if err != nil {
					fmt.Println("failed during skiplines")
					log.Fatal(err)
				}
				fmt.Println(string(line))
			}
		}

		if d[i+1] > 0 {
			skipChars := make([]byte, d[i+1]-cursor)
			cursor += len(skipChars)
			if _, err := io.ReadFull(fileReader, skipChars); err != nil {
				fmt.Println("failed during skipchars")
				log.Fatal(err)
			}
			fmt.Print(string(skipChars))
		}

		colorGreen := "\033[" + strconv.Itoa(31+(d[i+3]%7)) + "m"
		colorReset := "\033[0m"

		tokenChars := make([]byte, d[i+2])
		cursor = d[i+2]

		if _, err := io.ReadFull(fileReader, tokenChars); err != nil {
			fmt.Println("failed during token")
			log.Fatal(err)
		}
		fmt.Print(string(colorGreen) + string(tokenChars) + string(colorReset))
	}

	for {
		line, _, err := fileReader.ReadLine()
		if err != nil {
			break
		}
		fmt.Println(string(line))
	}
}
