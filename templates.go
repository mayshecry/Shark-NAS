package main

const fileListTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>SharkNAS - {{.CurrentPath}}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <script src="https://unpkg.com/lucide@latest"></script>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg: #000000;
            --sidebar: #000000;
            --surface: #0a0a0a;
            --border: #1a1a1a;
            --accent: #ffffff;
            --text: #ffffff;
            --text-muted: #666666;
            --danger: #ff453a;
            --font-main: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
        }

        * { box-sizing: border-box; margin: 0; padding: 0; }

        body {
            font-family: var(--font-main);
            background: var(--bg);
            color: var(--text);
            display: flex;
            height: 100vh;
            -webkit-font-smoothing: antialiased;
        }

        ::-webkit-scrollbar { width: 3px; }
        ::-webkit-scrollbar-thumb { background: #333; border-radius: 10px; }

        .sidebar {
            width: 260px;
            background: var(--sidebar);
            border-right: 1px solid var(--border);
            display: flex;
            flex-direction: column;
            padding: 40px 20px;
        }

        .logo { margin-bottom: 48px; padding: 0 12px; display: flex; align-items: center; gap: 14px; font-weight: 600; font-size: 1rem; letter-spacing: -0.02em; }
        .logo i { color: #fff; }

        .nav-item { 
            display: flex; 
            align-items: center; 
            gap: 14px; 
            padding: 12px 14px; 
            color: var(--text-muted); 
            text-decoration: none; 
            border-radius: 10px; 
            font-size: 0.85rem;
            font-weight: 500;
            transition: 0.2s cubic-bezier(0.4, 0, 0.2, 1);
            margin-bottom: 6px;
        }
        .nav-item:hover { background: #0f0f0f; color: #fff; }
        .nav-item.active { background: #ffffff; color: #000; }
        .nav-item.active i { color: #000; }

        .main-content { flex: 1; display: flex; flex-direction: column; overflow: hidden; }

        header { 
            height: 80px; 
            padding: 0 48px; 
            display: flex; 
            align-items: center; 
            gap: 40px;
            border-bottom: 1px solid var(--border); 
            background: var(--bg);
        }
        
        .search-bar { 
            flex: 1;
            background: #0a0a0a; 
            border: 1px solid var(--border); 
            padding: 10px 20px; 
            border-radius: 12px; 
            display: flex; 
            align-items: center; 
            gap: 16px; 
            transition: 0.2s ease;
        }
        .search-bar:focus-within { border-color: #444; background: #0f0f0f; }
        .search-bar input { background: transparent; border: none; color: white; outline: none; font-size: 0.9rem; width: 100%; font-weight: 400; }
        .search-bar i { color: var(--text-muted); width: 16px; }

        .metrics-bar {
            display: flex;
            gap: 40px;
            padding: 24px 48px;
            background: var(--bg);
            border-bottom: 1px solid var(--border);
        }
        .metric { display: flex; flex-direction: column; gap: 8px; min-width: 140px; }
        .metric-label { font-size: 0.65rem; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.05em; font-weight: 700; }
        .metric-value { font-size: 0.9rem; font-weight: 600; color: var(--text); }
        .metric-progress { width: 100%; height: 2px; background: #1a1a1a; border-radius: 1px; overflow: hidden; }
        .metric-fill { height: 100%; background: #fff; transition: width 1s cubic-bezier(0.4, 0, 0.2, 1); }

        .workspace { padding: 48px; overflow-y: auto; flex: 1; }
        .workspace.dragover {
            background: rgba(255, 255, 255, 0.03);
            outline: 2px dashed var(--border);
            outline-offset: -20px;
        }
        .breadcrumb { font-size: 0.75rem; color: var(--text-muted); margin-bottom: 32px; display: flex; align-items: center; gap: 10px; font-weight: 500; }
        .breadcrumb a { color: var(--text-muted); text-decoration: none; }
        .breadcrumb a:hover { color: #fff; }

        .file-grid { 
            display: grid; 
            grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); 
            gap: 24px; 
        }
        
        .file-card {
            background: var(--surface);
            border: 1px solid var(--border);
            padding: 24px;
            border-radius: 16px;
            display: flex;
            flex-direction: column;
            gap: 16px;
            transition: 0.2s cubic-bezier(0.4, 0, 0.2, 1);
            position: relative;
        }
        .file-card:hover { border-color: #444; background: #0f0f0f; transform: translateY(-4px); }

        .icon-box { 
            height: 140px; 
            display: flex; 
            align-items: center; 
            justify-content: center; 
            background: #000; 
            border-radius: 12px; 
        }
        .preview-img { width: 100%; height: 100%; object-fit: cover; border-radius: 10px; opacity: 0.9; }
        .file-card:hover .preview-img { opacity: 1; }

        .file-info { display: flex; flex-direction: column; gap: 4px; }
        .file-name { font-size: 0.9rem; font-weight: 600; color: #fff; text-decoration: none; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; letter-spacing: -0.01em; }
        .file-meta { font-size: 0.75rem; color: var(--text-muted); font-weight: 500; }

        .btn {
            padding: 10px 20px;
            border-radius: 10px;
            font-size: 0.85rem;
            font-weight: 600;
            cursor: pointer;
            display: flex;
            align-items: center;
            gap: 10px;
            border: 1px solid var(--border);
            transition: 0.2s;
        }
        .btn-primary { background: #ffffff; color: #000; border: none; }
        .btn-ghost { background: transparent; color: #ffffff; }
        .btn:hover { opacity: 0.8; }

        .quick-delete { position: absolute; top: 12px; right: 12px; opacity: 0; transition: 0.2s; color: var(--text-muted); border: none; background: none; cursor: pointer; }
        .file-card:hover .quick-delete { opacity: 1; }
        .quick-delete:hover { color: var(--danger); }

        .modal-overlay { display: none; position: fixed; inset: 0; background: rgba(0,0,0,0.9); backdrop-filter: blur(4px); z-index: 200; align-items: center; justify-content: center; }
        .modal-content { background: #0a0a0a; border: 1px solid var(--border); padding: 40px; border-radius: 20px; width: 440px; }
        .modal-content h2 { font-size: 1.5rem; margin-bottom: 24px; font-weight: 600; letter-spacing: -0.02em; }
        input[type="text"] { width: 100%; background: #000; border: 1px solid var(--border); padding: 14px; color: white; border-radius: 10px; margin-bottom: 20px; outline: none; }
    </style>
</head>
<body>
    <aside class="sidebar">
        <div class="logo">
            <i data-lucide="shield"></i>
            <span>SHARK_NAS</span>
        </div>
        <nav>
            <a href="/storage/" class="nav-item {{if or (eq .CurrentPath "/") (eq .CurrentPath "Root")}}active{{end}}"><i data-lucide="layers"></i>Files</a>
            <a href="/storage/recent" class="nav-item {{if eq .CurrentPath "Recent"}}active{{end}}"><i data-lucide="activity"></i>Recent</a>
            <a href="/storage/favorites" class="nav-item {{if eq .CurrentPath "Favorites"}}active{{end}}"><i data-lucide="bookmark"></i>Pinned</a>
            <a href="/storage/trash/" class="nav-item {{if eq .CurrentPath ".trash"}}active{{end}}"><i data-lucide="trash"></i>Trash</a>
        </nav>
    </aside>

    <main class="main-content">
        <header>
            <div class="search-bar">
                <i data-lucide="search" style="width: 16px"></i>
                <form action="/storage/search" method="GET" style="width:100%"><input type="text" name="q" placeholder="Search files..."></form>
            </div>
            <div class="action-btns">
                <button onclick="showModal('mkdirModal')" class="btn btn-ghost"><i data-lucide="folder-plus" style="width:16px"></i> New Folder</button>
                <button onclick="showModal('uploadModal')" class="btn btn-primary"><i data-lucide="upload" style="width:16px"></i> Upload</button>
            </div>
        </header>

        <div class="metrics-bar">
            <div class="metric">
                <span class="metric-label">Storage</span>
                <span class="metric-value">{{usagePercent .UsedStorage .MaxStorage | printf "%.1f"}}%</span>
                <div class="metric-progress"><div class="metric-fill" style="width: {{usagePercent .UsedStorage .MaxStorage}}%"></div></div>
            </div>
            <div class="metric">
                <span class="metric-label">Memory</span>
                <span class="metric-value">{{.RAMPercent | printf "%.1f"}}%</span>
                <div class="metric-progress"><div class="metric-fill" style="width: {{.RAMPercent}}%"></div></div>
            </div>
            <div class="metric">
                <span class="metric-label">Threads</span>
                <span class="metric-value">{{.CPUUsage}}</span>
                <div class="metric-progress"><div class="metric-fill" style="width: 10%"></div></div>
            </div>
        </div>

        <div class="workspace">
            <div class="breadcrumb">
                <a href="/storage/">ROOT</a>
                {{if and (ne .CurrentPath "/") (ne .CurrentPath "Recent") (ne .CurrentPath "Favorites")}}
                    <i data-lucide="chevron-right" style="width:12px"></i> <span>{{.CurrentPath}}</span>
                {{end}}
            </div>

            <div class="file-grid">
                {{range .Files}}
                <div class="file-card">
                    <form action="/storage/delete" method="POST" onsubmit="return confirm('Delete {{.Name}}?')">
                        <input type="hidden" name="path" value="{{.Path}}">
                        <button type="submit" class="quick-delete"><i data-lucide="trash-2" style="width:14px"></i></button>
                    </form>
                    <div class="icon-box">
                    {{if .IsDir}}
                        <i data-lucide="folder" style="width: 32px; height: 32px; color: #71717a;"></i>
                    {{else}}
                        {{if (isImage .Name)}}<img src="/storage{{.Path}}" class="preview-img" onclick="openPreview(this.src)">
                        {{else}}<i data-lucide="file" style="width: 32px; height: 32px; color: #3f3f46;"></i>{{end}}
                    {{end}}
                    </div>
                    <div class="file-info">
                        <a href="/storage{{.Path}}" class="file-name">{{.Name}}</a>
                        <span class="file-meta">{{if .IsDir}}Folder{{else}}{{formatSize .Size}}{{end}}</span>
                    </div>
                </div>
                {{else}}
                <div style="grid-column: 1/-1; padding: 80px; text-align: center; color: #3f3f46;">No files found</div>
                {{end}}
            </div>
        </div>
    </main>

    <div id="uploadModal" class="modal-overlay" onclick="hideModal('uploadModal')">
        <div class="modal-content" onclick="event.stopPropagation()">
            <h2>Upload File</h2>
            <form action="/storage/upload?dir={{.CurrentPath}}" method="POST" enctype="multipart/form-data" onsubmit="this.querySelector('button').innerText='Uploading...'; this.querySelector('button').style.opacity='0.5'; this.querySelector('button').style.pointerEvents='none';">
                <input type="file" name="uploadFile" multiple required>
                <button type="submit" class="btn btn-primary" style="width:100%; margin-top:20px; justify-content: center">Upload</button>
            </form>
        </div>
    </div>

    <div id="mkdirModal" class="modal-overlay" onclick="hideModal('mkdirModal')">
        <div class="modal-content" onclick="event.stopPropagation()">
            <h2>New Folder</h2>
            <form action="/storage/mkdir" method="POST">
                <input type="hidden" name="dir" value="{{.CurrentPath}}">
                <input type="text" name="name" placeholder="Folder name" required autofocus>
                <button type="submit" class="btn btn-primary" style="width:100%; justify-content: center">Create</button>
            </form>
        </div>
    </div>

    <div id="previewModal" class="modal-overlay" onclick="hideModal('previewModal')">
        <img id="modalImg" src="" onclick="event.stopPropagation()">
    </div>

    <script>
        const workspace = document.querySelector('.workspace');

        function showModal(id) { document.getElementById(id).style.display = 'flex'; }
        function hideModal(id) { document.getElementById(id).style.display = 'none'; }
        function openPreview(src) { 
            document.getElementById('modalImg').src = src;
            showModal('previewModal');
        }

        // Drag and Drop Logic
        ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
            workspace.addEventListener(eventName, e => {
                e.preventDefault();
                e.stopPropagation();
            }, false);
        });

        ['dragenter', 'dragover'].forEach(eventName => {
            workspace.addEventListener(eventName, () => workspace.classList.add('dragover'), false);
        });

        ['dragleave', 'drop'].forEach(eventName => {
            workspace.addEventListener(eventName, () => workspace.classList.remove('dragover'), false);
        });

        workspace.addEventListener('drop', e => {
            const files = e.dataTransfer.files;
            if (files.length > 0) {
                const formData = new FormData();
                for (let i = 0; i < files.length; i++) {
                    formData.append('uploadFile', files[i]);
                }
                workspace.style.opacity = '0.5';
                workspace.style.pointerEvents = 'none';

                fetch('/storage/upload?dir={{.CurrentPath}}', {
                    method: 'POST',
                    body: formData
                }).then(() => window.location.reload())
                  .catch(() => {
                      alert('Upload failed');
                      workspace.style.opacity = '1';
                      workspace.style.pointerEvents = 'auto';
                  });
            }
        }, false);

        lucide.createIcons();
    </script>
</body>
</html>
`
