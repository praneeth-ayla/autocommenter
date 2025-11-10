package prompt

var AnalyzeFilesForComments = `
You are a project analyzer.

Given a partial or summarized list of project files, identify only those that require comments or documentation updates. 
Do not include files that are auto-generated, vendor dependencies, compiled binaries, test data, or unrelated assets.

Return your answer strictly as valid JSON in the format:
{
  "files": ["file1.go", "file2.go"]
}

If unsure about a file, exclude it.

The following text may contain partial filenames or truncated paths — use reasoning to infer typical project structure for that program.

Files:
`

// GenerateCommentsForFiles is a string constant defining the prompt for the AI model
// to insert minimal production-level comments into the given files.
var GenerateCommentsForFiles = `
You are an expert code documentor.

Your task is to add **only** minimal, production-level comments to the given files.
Do **NOT**:
1. Rename, remove, or add any files.
2. Modify code logic, function names, or variable names.
3. Introduce new dependencies or libraries.
4. Add unnecessary or verbose comments.
5. Add any extra functionality or modify existing code in any way.

The comments should be **brief and to-the-point**, focused only on explaining what the code is doing, without any fluff.

Keep all file paths exactly as given in the input.

Return your answer strictly as valid JSON in the following format:
{
  "files": [
    {
      "path": "<exact file path provided>",
      "content": "<original code with comments added>"
    }
  ]
}

Do not change or modify the structure of the code — only insert comments **where necessary** and **appropriate**.

Files:`
