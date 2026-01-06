# FilePecker 

**FilePecker** is a lightweight, robust CLI tool written in Go. It recursively scans your current directory and consolidates the content of all files into a single text file (`file.txt`).

It is designed to be **smart**: it automatically ignores binary files (images, executables), `.git` folders, and system files, making it perfect for preparing codebases for LLM context, documentation, or archiving.

---

## 🚀 Quick Install

Run this single command in your terminal to download and install `filepecker` automatically (Linux & Mac):

```bash
curl -fsSL https://raw.githubusercontent.com/Sarthak160/filepecker/main/install.sh | bash
```

## ⚙️ Advanced Usage

You can customize the behavior using flags:

| Flag | Description | Example |
| :--- | :--- | :--- |
| `-o` | Set the output filename (default: file.txt) | `filepecker -o mycode.txt` |
| `-ignore` | Comma-separated list of extensions to skip | `filepecker -ignore .json,.css` |

**Example:**
Scan the directory, but output to `backup.txt` and ignore all markdown and JSON files:

```bash
filepecker -o backup.txt -ignore .md,.json