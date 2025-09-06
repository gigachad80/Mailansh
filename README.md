
🚀 Project Name : Mailansh
===============


### Mailansh : OSINT tool to extract emails of contributors from Git repositories across multiple platforms with advanced filtering capabilities.


 
![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-purple.svg)
<a href="https://github.com/gigachad80/Mailansh/issues"><img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat"></a>

## Table of Contents

* [📌 Overview](#-overview)
* [✨ Features](#-features)
* [📚 Requirements & Dependencies](#-requirements--dependencies)
* [📥 Installation Guide](#-installation-guide)
* [🚀 Usage](#-usage)
  - [Basic Usage](#basic-usage)
* [📋 Command Line Options](#-command-line-options)
* [📝Roadmap](#-roadmap)
* [🔧 Technical Details](#-technical-details)
* [🤔 Why This Name?](#-why-this-name)
* [⌚ Development Time](#-development-time)
* [🙃 Why I Created This](#-why-i-created-this)
* [🤝 Contributing](#-contributing)
* [📞 Contact](#-contact)
* [📄 License](#-license)

### 📌 Overview

**Mailansh** is an OSINT tool to extract emails  of contributors from Git repositories across multiple platforms. It provides advanced filtering capabilities and concurrent processing for optimal performance, making it perfect for repository analysis, team insights, and contributor research.

**Key Capabilities:**
* Multi-platform repository support (GitHub, GitLab, Gitea, Bitbucket)
* High-performance concurrent processing with worker pools
* Advanced email filtering and categorization
* Smart deduplication with preference for real names
* Multiple output formats (formatted table, CSV)
* Real-time progress tracking

### ✨ Features

### 🚀 Multi-Platform Support
- **GitHub** 
- **GitLab** 
- **Gitea** 
- **Bitbucket** 

### ⚡ High Performance
- **Concurrent Processing** - 8 worker goroutines for optimal speed
- **Dual Extraction** - Processes both commit logs and patch content simultaneously
- **Progress Tracking** - Real-time progress updates with percentage completion
- **Smart Deduplication** - Prefers real names over noreply usernames
- **Memory Efficient** - Streaming processing with automatic cleanup

### 🎯 Advanced Filtering
- **GitHub NoReply** (`-g`) - Show only GitHub noreply emails
- **Popular Domains** (`-p`) - Filter for Gmail, Outlook, Yahoo, ProtonMail, etc.
- **Custom Domains** (`-cd`) - Show only custom/corporate domains
- **Quiet Mode** (`-q`) - CSV output perfect for scripting and data processing

### 🛡️ Smart Email Detection
- Platform-specific noreply patterns recognition
- Comprehensive email validation and sanitization
- Invalid email filtering (test domains, localhost, malformed addresses)
- Thread-safe contributor management with mutex protection

### 📚 Requirements & Dependencies

* **Go 1.19+** - Latest version recommended for optimal performance
* **Git** - Must be installed and accessible in system PATH


### 📥 Installation Guide

### ⚡ Quick Install

**Method 1: Build from Source**
```bash
git clone https://github.com/gigachad80/Mailansh
cd Mailansh/cmd/mailansh/
go build -o mailansh main.go
```
OR 

**Method 3: Download Binary**
Download the latest binary from the [releases page](https://github.com/gigachad80/Mailansh/releases) and add it to your PATH.

### 🚀 Usage

### Basic Usage

```bash
# Extract all contributors
./mailansh https://github.com/user/repo

# G#itHub noreply emails only
./mailansh -g https://github.com/user/repo 

# Popular email domains only
./mailansh -p  https://github.com/user/repo 

# Custom/corporate domains only
./mailansh -cd  https://github.com/user/repo -

# Extract Popular email domains only and ave them
./mailansh -p -o popular.csv https://github.com/user/repo 

# Save to file
./mailansh  -o contributors.csv https://github.com/user/repo 

# Help and usage information
./mailansh -h


```

> [!WARNING]
> #### 📝 Known Behavior & FAQs
> ### Progress Display Not Reaching 100%
> You might notice the progress indicator stopping at something like "189/221 commits (85.5%)" instead of 100%. This is **normal behavior** and doesn't affect the completeness of results.
> ### Why this happens:
> * The tool uses a **dual extraction strategy** with concurrent workers
> * **Method 1**: `git log` extraction (always completes successfully)
> * **Method 2**: Patch content analysis with 8 concurrent workers
> * Some git commits may fail to process (merge commits, corrupted patches, etc.)
> * Workers may terminate early on errors, but **all emails are still found**
> ### What this means:
> * ✅ **Your results are complete** - all contributor emails are extracted
> * ✅ The primary `git log` method captures all standard contributors
> * ✅ Patch analysis adds any additional emails found in commit content
> * ✅ Failed commits don't contain unique contributor information
> **Note: Even if patch processing shows <100%, the contributor list remains accurate and complete.**
> ### Example output:
```bash
[INFO] Extracting from commit history...
[INFO] Extracting from patch content...
[PROGRESS] 189/221 commits (85.5%)
[INFO] Found 45 unique contributors (after filtering):
```



### 📋 Command Line Options

| Flag | Description | Example Usage |
|------|-------------|---------------|
| `-h` | Show help menu | `./mailansh -h` |
| `--quiet`, `-q` | CSV output suitable for redirection | `./mailansh repo -q > output.csv` |
| `-g` | Show only GitHub noreply emails | `./mailansh -g repo ` |
| `-p` | Show only popular domain emails | `./mailansh -p repo` |
| `-cd` | Show only custom/corporate domains | `./mailansh =cd repo ` |
| `-o <file>` | Save output to file (.csv or .txt) | `./mailansh  -o results.csv repo` |

### 📝 Roadmap

- [ ] Develop Web UI
- [x] Release cross platform executables / binaries

### 🔧 Technical Details

### Architecture
- **Concurrent Design** - Parallel extraction from commit log and patches
- **Worker Pool Pattern** - 8 goroutines process commits simultaneously  
- **Thread-Safe Operations** - Mutex-protected contributor management
- **Memory Efficient** - Streaming processing with temporary cleanup
- **Context-Aware** - Proper cancellation and timeout handling



### Performance Characteristics
- **Large Repositories** - Efficiently handles 10,000+ commits
- **Memory Usage** - ~50MB average for most repositories
- **Processing Speed** - ~100-500 commits/second depending on system specs
- **Concurrency** - 8 worker goroutines + 2 extraction goroutines

### 🤔 Why This Name?

Can't disclose due to personal reasons ╰(*°▽°*)╯

### ⌚ Development Time

From initial concept to feature-complete implementation, including testing, optimization, and documentation, the development took approximately **4 hours**.

### 🙃 Why I Created This

Actually, I wanted to develop something else but this got developed by mistake💀 and because of this mess, that project was left incomplete 

### 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. **Fork the repository**
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)
3. **Commit your changes** (`git commit -m 'Add amazing feature'`)
4. **Push to the branch** (`git push origin feature/amazing-feature`)
5. **Open a Pull Request**

### Development Guidelines
- Follow Go best practices and conventions
- Add tests for new features
- Update documentation for any changes
- Ensure backward compatibility

### 📞 Contact

📧 Email: **pookielinuxuser@tutamail.com**

---

### 📄 License

Licensed under the **RPL 1.5** and a **Custom License**.  
Check here: [`CREDITS.md`](https://github.com/gigachad80/Mailansh/blob/main/CREDITS.md) (Important)  
Also see: [`LICENSE.md`](https://github.com/gigachad80/Mailansh/blob/main/LICENCE.md)


---


**Made with ❤️ in Go** - Fast, concurrent, and reliable contributor analysis.
