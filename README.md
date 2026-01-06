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