package prompt

const TemplateFileContext = `
Analyze the Go file and output a single JSON object with exactly these fields (no extras):

* path: string
* file_name: string
* summary: short (<=50 words) summary describing the file's purpose and its runtime logic. Include important behavior only: flags and default values, file reads/writes, external calls (providers, stores, scanners, etc.), control-flow decisions, and observable side effects or error returns.
* exports: array of exported identifiers (names only)
* imports: array of imported package paths (literal strings as in the import block)

Rules:

1. Do not include any fields other than the five listed.
2. Keep the summary concise and focused; avoid listing local variables or low-level implementation details.
3. List exported symbols exactly as they appear (identifiers with capitalized names).
4. List imports as the package paths shown in the source.
5. Return valid JSON only. Do not include explanations, commentary, code fences, or extra text.

Path:
%s

Content:
%s
`
