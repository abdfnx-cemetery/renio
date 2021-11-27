// Copyright 2021 @abdfnx. All rights reserved. Apache 2.0 license.
const IGNORED_DIAGNOSTICS = [
  2306,
  1375,
  1103,
  2691,
  5009,
  5055,
  5070,
  7016,
];

/*
  TS2306: File 'FILENAME.js' is not a module.

  TS1375: 'await' expressions are only allowed at the top level of a file
  when that file is a module, but this file has no imports or exports.
  Consider adding an empty 'export {}' to make this file a module.

  TS1103: 'for-await-of' statement is only allowed within an async function
  or async generator.

  TS2691: An import path cannot end with a '.ts' extension. Consider
  importing 'bad-module' instead.

  TS5009: Cannot find the common subdirectory path for the input files.

  TS5055: Cannot write file 'x.js' because it would overwrite input file.
  TypeScript is overly opinionated that only CommonJS modules kinds can
  support JSON imports.  Allegedly this was fixed in
  Microsoft/TypeScript#26825 but that doesn't seem to be working here,
  so we will ignore complaints about this compiler setting.

  TS7016: Could not find a declaration file for module '...'. '...'
  implicitly has an 'any' type.  This is due to `allowJs` being off by
  default but importing of a JavaScript module.
*/

const options = { allowNonTsExtensions: true };

function typeCheck(file, source) {
  const dummyFilePath = file;
  const textAst = ts.createSourceFile(
    dummyFilePath,
    source,
    ts.ScriptTarget.ES6
  );
  const dtsAST = ts.createSourceFile(
    "/lib.es6.d.ts",
    Asset("typescript/lib.es6.d.ts"),
    ts.ScriptTarget.ES6
  );

  const files = { [dummyFilePath]: textAst, "/lib.es6.d.ts": dtsAST };
  const host = {
    fileExists: (filePath) => {
      return files[filePath] != null || Renio.exists(filePath);
    },
    directoryExists: (dirPath) => dirPath === "/",
    getCurrentDirectory: () => Renio.cwd(),
    getDirectories: () => [],
    getCanonicalFileName: (fileName) => fileName,
    getNewLine: () => "\n",
    getDefaultLibFileName: () => "/lib.es6.d.ts",
    getSourceFile: (filePath) => {
      if (files[filePath] != null) return files[filePath];
      else {
        return ts.createSourceFile(
          filePath,
          Renio.readFile(filePath),
          ts.ScriptTarget.ES6
        );
      }
    },
    readFile: (filePath) => {
      return filePath === dummyFilePath ? text : Renio.readFile(filePath);
    },
    useCaseSensitiveFileNames: () => true,
    writeFile: () => {},
    resolveModuleNames,
  };
  const program = ts.createProgram({
    options,
    rootNames: [dummyFilePath],
    host,
  });

  let diag = ts.getPreEmitDiagnostics(program).filter(function ({ code }) {
    return code != 5023 && !IGNORED_DIAGNOSTICS.includes(code);
  });
  let diags = ts.formatDiagnosticsWithColorAndContext(diag, host);
  Report(diags);
}

function resolveModuleNames(moduleNames, containingFile) {
  const resolvedModules = [];
  for (const moduleName of moduleNames) {
    let fileName = join(containingFile, "..", moduleName);
    if (moduleName.startsWith("https://")) {
      fileName = moduleName.replace("https://", "/tmp/");
    }
    resolvedModules.push({ resolvedFileName: fileName });
  }
  return resolvedModules;
}

// Joins path segments.  Preserves initial "/" and resolves ".." and "."
// Does not support using ".." to go above/outside the root.
// This means that join("foo", "../../bar") will not resolve to "../bar"
function join(/* path segments */) {
  // Split the inputs into a list of path commands.
  var parts = [];
  for (var i = 0, l = arguments.length; i < l; i++) {
    parts = parts.concat(arguments[i].split("/"));
  }
  // Interpret the path commands to get the new resolved path.
  var newParts = [];
  for (i = 0, l = parts.length; i < l; i++) {
    var part = parts[i];
    // Remove leading and trailing slashes
    // Also remove "." segments
    if (!part || part === ".") continue;
    // Interpret ".." to pop the last segment
    if (part === "..") newParts.pop();
    // Push new path segments.
    else newParts.push(part);
  }
  // Preserve the initial slash if there was one.
  if (parts[0] === "") newParts.unshift("");
  // Turn back into a single string path.
  return newParts.join("/") || (newParts.length ? "/" : ".");
}
