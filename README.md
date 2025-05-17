# 🖥️ CoreStat

**CoreStat** is a lightweight, cross-platform system monitoring application built with **Golang**. It provides real-time insights into running processes, CPU usage, memory consumption, and more, with an intuitive user interface and modern features.

## 🚀 Features

- 🔍 View detailed system process information (PID, memory, CPU, threads, etc.)
- 📊 Real-time memory and CPU usage stats
- 🧬 Displays parent-child process hierarchy
- 📂 Collapse/expand child processes
- 🧠 Built with performance in mind using native Go
- 🌐 Web-based or desktop UI (based on your tech stack like Wails/Svelte)

## 🛠️ Technologies Used

- **Go** – Backend logic and system resource monitoring
- **[Svelte](https://svelte.dev/)** – Reactive front-end (if using Wails or web UI)
- **[gopsutil](https://github.com/shirou/gopsutil)** – System information gathering
- **[Wails](https://wails.io/)** (if used) – To build cross-platform desktop apps with Go + Web

## 🧑‍💻 Installation

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.18+)
- Node.js & npm (if using Svelte or Wails frontend)

### Clone the repository

```bash
git clone https://github.com/yourusername/corestat.git
cd corestat
