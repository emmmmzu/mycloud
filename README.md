# MyCloud

## *⚠️ Learning Project*

This project is primarily for learning purposes. It is a self-hosted cloud file server written in Go.  
The main goal is to practice building a REST API, file handling, and secure server logic.  

**Future plans / roadmap:**

- Implement a **Rust-based engine** for the client side
- Build a **cross-platform desktop application** (Windows/Linux) to mount the server as a virtual drive  
- Develop an **Android application** for mobile access  
- Enhance security with **authentication and tokens**  
- Add **file downloads, syncing, and versioning** to mimic Dropbox functionality  

---

Access and manage your files via a REST API similar to Dropbox.  

---

## Features

- **File browsing:** List folders and file details (`/list`, `/stat`)
- **File upload:** Upload files to any folder (`/upload`)
- **Secure paths:** Prevents directory traversal attacks (`../`)
- **Automatic folder creation:** Uploading to a new folder creates it automatically
- **JSON responses:** Consistent API responses with error handling

---

## Installation

### 1. **Clone the repository**

```bash
git clone https://github.com/emmmmzu/mycloud.git
cd mycloud
```

### 2. **Initialize Go modules**

```bash
go mod tidy
```

### 3. **Run the server**

```bash
go run main.go
```

By default, the server runs on port `8080`. You can configure the root folder by setting the `rootFolder` variable in `main.go`.

---

## API Endpoints

### 1. GET `/list?path=...`

List files and folders in a directory.

- **Query parameter:** `path` — target folder relative to the root
- **Success response:** `200 OK`

```json
[
  { "name": "file1.txt", "type": "file", "size": 1234, "modified": "2025-10-08T12:34:56Z" },
  { "name": "docs", "type": "folder", "size": 0, "modified": "2025-10-01T08:00:00Z" }
]
```

### 2. GET `/stat?path=...`

Get file/folder details.

- **Query parameter:** `path` — target file or folder
- **Success response:** `200 OK`

```json
{
  "path": "/some/path/file1.txt",
  "type": "file",
  "size": 1234,
  "modified": "2025-10-08T12:34:56Z"
}
```

### 3. POST `/upload?path=...`

Upload a file to the server.

- **Query parameter:** `path` — target folder relative to the root
- **Form-data field:** `file` — file to upload
- **Success response:** `200 OK`

```json
{
  "message": "file uploaded successfully",
  "filename": "file.txt",
  "size": 12345
}
```

---

## Error Handling

All errors are returned in JSON format:

```json
{
  "error": "Error message",
  "status": "HTTP Status Description"
}
```

Examples:

- Missing path: `400 Bad Request`
- Invalid path (outside root): `403 Forbidden`
- File system errors: `500 Internal Server Error`

---

## Testing Uploads with `curl`

**Windows (PowerShell):**

```powershell
curl -X POST "http://localhost:8080/upload" `
     -F "path=/" `
     -F "file=@C:\path\to\testfile.txt"
```

**Command Prompt:**

```cmd
curl -X POST "http://localhost:8080/upload" -F "path=/" -F "file=@C:\path\to\testfile.txt"
```

---

## Notes

- Currently supports **Windows and Linux**.
- Designed for **local testing**; authentication will be added in future steps.
- Filenames and folder names are sanitized to prevent path traversal.
- Directories are created automatically when uploading to a non-existing folder.

---

## License

MIT License
