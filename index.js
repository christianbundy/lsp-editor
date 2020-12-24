const childProcess = require("child_process");
const path = require("path");
const os = require("os");

const directory = `file://${__dirname}`;
const file = `file://${__dirname}/hello.go`;
const draw = require("./draw");

const gopls = childProcess.spawn("gopls", ["-vv"]);

gopls.on("error", (err) => {
  throw err;
});
gopls.stderr.on("data", (data) => {
  console.log("stderr", data.toString());
});
gopls.on("close", (code) => {
  process.exit(code);
});

let id = 0;

function request(method, params) {
  id += 1;
  let thisId = id;
  const data = JSON.stringify({
    jsonrpc: "2.0",
    id,
    method,
    params,
  });
  const payload = `Content-Length: ${data.length}\r\n\r\n${data}`;

  return new Promise((resolve) => {
    gopls.stdout.on("data", (data) => {
      const json = data.toString().split("\r\n\r\n")[1];

      try {
        const obj = JSON.parse(json);
        if (obj.id === thisId) {
          resolve(obj);
        }
      } catch {
        console.log("Bad JSON", data.toString());
      }
    });
    gopls.stdin.write(payload, "utf8", (err) => {
      if (err) throw err;
    });
  });
}

async function main() {
  await request("initialize", {
    rootUri: file,
    workspaceFolders: [
      {
        uri: directory,
        name: "lsp-editor",
      },
    ],
    capabilities: {
      textDocument: {
        semanticTokens: {
          requests: { full: true },
          tokenTypes: [
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
          ],
          tokenModifiers: [
            "declaration",
            "definition",
            "readonly",
            "static",
            "deprecated",
            "abstract",
            "async",
            "modification",
            "documentation",
            "defaultLibrary",
          ],
        },
      },
    },
    initializationOptions: { semanticTokens: true },
  });

  await request("initialized", {});

  const { result } = await request("textDocument/semanticTokens/full", {
    textDocument: {
      uri: file,
    },
  });

  draw(result.data);
  gopls.kill();
}

main();
