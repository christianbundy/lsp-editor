# LSP-Editor

My first contribution to open source software was me fixing some small problem
in a syntax highlighter. After a decade and some change, I'm banging my head
against Vim plugins trying to figure out how I can get syntax highlighting,
diagnostics (static analysis, linters, etc), and auto-complete in my editor.

At the same time, I've noticed a family of pain points around the fact that
**code can be serialized as text, but that's not what it is**. The underlying
data structure are a graph, the syntax is [usually] a tree, but instead of
editing either of those we find ourselves editing bytes and lines (i.e. text).

I've started to wonder whether there's a unified solution to both problems,
and whether we can move past text editors to **semantic editors**, or at least
we can move a bit in that direction.

Anyway, I've found myself wondering whether the [Language Server Protocol][lsp]
might be a worthwhile exploration. Most tools seem like they're [slowly]
starting to converge around this protocol (and its [index format][lsif]), and
it might be interesting to have a small software editor where all of the heavy
lifting is done by various language servers.

## Status

Not even close.

My first experiment: use [semantic tokens][semtok] to do some basic syntax
highlighting. My day basically went like this:

- Isn't `tsserver` a language server? I bet this will be easy.
- Oh... `tsserver` was invented before LSP, and doesn't support it.
- But there's a wrapper!
- But it's archived.
- But there's another wrapper!
- ...but `textDocument/semanticTokens` isn't available. :(
- Maybe there's some non-standard way to do it with `tsserver`?
- Nope, semantic tokens are mentioned in the Typescript source but not exported.
- I'll make an issue: https://github.com/microsoft/TypeScript/issues/42091
- What else is likely to have a language server? Python?
- No `textDocument/semanticTokens` in that language server either.
- Rust? Grepping for `semanticTokens` doesn't look good.
- ...Go?! Oh wow, there's an [official language server][gopls]?!?!
- AND IT SUPPORTS `textDocument/semanticTokens` HALLELUJAH

TL;DR: A normal morning as a software engineer.

My original prototype was written in JavaScript (Node.js), but I decided to
rewrite the project in Go so that I can [eventually] use my editor to edit
itself and also have an excuse to suck a bit less at Go. My code is probably
very bad and you shouldn't assume anything I'm doing is a best practice.

## Try

...but if you want to try it anyway.

1. Install Go
2. [Install gopls](https://github.com/golang/tools/blob/master/gopls/doc/user.md#installation)
4. `go build`
5. `./lsp-editor main.go`

If you're very lucky, you might see something like this:

![Screenshot of `hello.go` with syntax highlighting][screenshot]

## License

AGPL-3.0-Only

[lsp]: https://microsoft.github.io/language-server-protocol/
[lsif]: https://lsif.dev/
[semtok]: https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_semanticTokens
[gopls]: https://github.com/golang/tools/tree/master/gopls
[screenshot]: https://user-images.githubusercontent.com/537700/103044416-40f03780-4535-11eb-9671-9d1c3368116e.png
