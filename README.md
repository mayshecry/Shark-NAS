# SharkNAS 🦈

SharkNAS is a lightweight, high-performance Network Attached Storage (NAS) solution written in Go. It features a modern, purple-themed glassmorphism UI and a robust set of file management tools.

## ✨ Features

- **Modern UI**: Sleek glassmorphism design with a dark purple and black color scheme.
- **Multi-File Upload**: Upload multiple files simultaneously with ease.
- **Storage Management**: 
    - Configurable storage limits (default 1GB).
    - Real-time storage usage visualization.
- **Smart Tracking (IP-Based)**:
    - **Favorites**: Star your important files.
    - **Recent**: Quickly access files you've recently viewed or uploaded.
- **Trash System**: Safety first! Deleted files are moved to a hidden `.trash` folder before permanent removal.
- **Deep Search**: Instantly find files and folders across your entire storage directory.
- **Image Previews**: Built-in modal viewer for images (jpg, png, webp, gif).
- **Folder Creation**: Organize your NAS with custom directories.

## 🚀 Getting Started

### Prerequisites
- Go (1.16 or higher recommended)

### Installation

1. Clone the repository or navigate to your project folder.
2. Ensure your directory structure looks like this:
   ```text
   .
   ├── main.go
   ├── handlers.go
   ├── models.go
   ├── templates.go
   ├── www/          # Static website files
   └── storage/      # Your NAS files
   ```

### Running the Server

```bash
go run .
```

The server will start on `http://localhost:50`.

## ⚙️ Configuration

You can adjust constants in `main.go`:
- `serverPort`: The port the NAS runs on.
- `storageDir`: The root directory for file storage.
- `maxStorageSize`: The maximum byte limit for the entire storage directory.

## 🛠️ Tech Stack
- **Backend**: Go (Standard Library)
- **Frontend**: HTML5, CSS3 (Glassmorphism), JavaScript
- **Icons**: Lucide Icons
- **Metadata**: Persistent JSON storage for favorites and recents.