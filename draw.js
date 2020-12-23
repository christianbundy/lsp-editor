const chalk = require("chalk");
const fs = require("fs");

const sourceFileName = "hello.go";
const sourceFileString = fs.readFileSync(sourceFileName, "utf8");
const sourceLines = sourceFileString.split("\n");

module.exports = (lspData) => {
  let currentLine = 0;
  let currentChar = 0;
  const prettyLspData = lspData.reduce((acc, cur, idx) => {
    switch (idx % 5) {
      case 0:
        currentLine += cur;
        if (cur > 0) {
          currentChar = 0;
        }
        acc.push({
          line: currentLine,
        });
        break;
      case 1:
        currentChar += cur;
        acc[acc.length - 1].startChar = currentChar;
        break;
      case 2:
        acc[acc.length - 1].length = cur;
        break;
      case 3:
        acc[acc.length - 1].tokenType = cur;
        break;
      case 4:
        // TODO
        break;
    }
    return acc;
  }, []);

  sourceLines.map((line, lineIndex) => {
    prettyLspData
      .filter((token) => token.line === lineIndex)
      .reverse()
      .forEach((token) => {
        const before = line.slice(0, token.startChar);
        const during = line.slice(
          token.startChar,
          token.startChar + token.length
        );
        const after = line.slice(token.startChar + token.length);

        const hue = (360 / 21) * token.tokenType;
        const colored = chalk.hsv(hue, 50, 100)(during);
        line = `${before}${colored}${after}`;
      });
    process.stdout.write(`${line}\n`);
  });
};
