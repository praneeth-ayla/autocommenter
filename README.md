# AutoCommenter 🧠

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**AutoCommenter** is a command-line tool that leverages AI to automate the tedious process of writing code comments and project documentation. It scans your codebase, builds a contextual understanding using Google's Gemini model, and then uses that context to generate high-quality, meaningful comments and a comprehensive `README.md` file.

---

## ✨ Features

*   **AI-Powered Context Generation**: Scans your entire project to understand the purpose of each file, its exports, and its relationship with other modules.
*   **Automated Code Commenting**: Intelligently inserts detailed, context-aware comments directly into your source files.
*   **README.md Generation**: Automatically creates a well-structured `README.md` for your project based on the overall code context.
*   **Two-Step Process**: Ensures high-quality output by first building a project-wide context before generating comments or documentation.
*   **Modular AI Backend**: Built with a provider interface to easily support different AI models in the future (currently implemented with Google Gemini).

---

## ⚙️ How It Works

The tool operates in a two-step process to ensure high-quality, context-aware output:

1.  **Generate Context**: First, `AutoCommenter` scans your project to understand the purpose of each file. This information is summarized by the AI and stored locally in a `.autocommenter/context.json` file. This step builds a holistic understanding of your entire codebase.

2.  **Generate Artifacts**: With the project context established, you can generate code comments or a README. This approach ensures that the generated content is not just based on a single file, but on the project as a whole, leading to more accurate and relevant results.

---

## 🚀 Getting Started

### Prerequisites

*   Go 1.21 or later.
*   A Google Gemini API key.

### 1. Installation

You can install `AutoCommenter` directly using `go install`:

go install github.com/praneeth-ayla/AutoCommenter/cmd/autocommenter@latest

### 2. Configuration

Export your Google Gemini API key as an environment variable:

export GEMINI_API_KEY="YOUR_API_KEY_HERE"

---

## 🧰 Usage

All commands should be run from the root directory of your project.

### Step 1: Generate Project Context

This is the first and most important step. It creates the knowledge base the tool uses for all other operations.

autocommenter context gen

This command will:
*   Recursively scan your source files.
*   Call the Gemini API to generate a summary for each file.
*   Save this context to a `.autocommenter/context.json` file in your project root.

### Step 2 (Option A): Generate Code Comments

Once the context is generated, you can add comments to your files. The tool will identify which files need comments and generate them.

autocommenter comments gen

### Step 2 (Option B): Generate a README File

To generate a new `README.md` for your project based on the code context:

autocommenter readme gen

If a `README.md` file already exists, the tool will attempt to merge the AI-generated content with your existing file.

---

## 🧱 Project Structure

A brief overview of the `AutoCommenter` internal structure:

AutoCommenter/
├── cmd/                  # Cobra CLI command definitions (root, context, comments, readme)
├── internal/
│   ├── ai/               # AI provider interface and Gemini implementation
│   ├── contextstore/     # Logic for saving and loading project context from JSON
│   ├── prompt/           # Builders for constructing prompts sent to the AI
│   └── scanner/          # File system scanning, batching, and I/O utilities
├── main.go               # Main application entry point
└── go.mod

---

## 🌐 Future Enhancements

*   Support for more programming languages.
*   Integration with other AI providers (e.g., OpenAI, Anthropic).
*   Fine-grained control over commenting style and verbosity.
*   A VS Code extension for a more integrated workflow.

Contributions are welcome! Feel free to open an issue or submit a pull request.

---

## 📜 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.