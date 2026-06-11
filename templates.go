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
        :root { --bg: #0A0A0C; --sidebar: rgba(18, 10, 25, 0.75); --card: rgba(30, 15, 45, 0.1); --border: rgba(50, 25, 75, 0.2); --accent: #A78BFA; --text: #E0E0E0; --text-muted: #A0A0A0; }
        * { box-sizing: border-box; }
        body { font-family: 'Inter', sans-serif; margin: 0; background: radial-gradient(circle at top left, #2A1A3A, #0A0A0C); color: var(--text); display: flex; height: 100vh; overflow: hidden; }

        .sidebar { width: 280px; background: var(--sidebar); border-right: 1px solid var(--border); display: flex; flex-direction: column; padding: 24px; backdrop-filter: blur(20px); -webkit-backdrop-filter: blur(20px); }
        .logo { display: flex; align-items: center; gap: 12px; font-weight: 600; font-size: 1.25rem; margin-bottom: 40px; color: var(--accent); }
        .nav-item { display: flex; align-items: center; gap: 12px; padding: 10px 14px; color: var(--text-muted); text-decoration: none; border-radius: 8px; transition: 0.2s; margin-bottom: 4px; border: 1px solid transparent; background: rgba(0,0,0,0.1); }
        .nav-item:hover, .nav-item.active { background: rgba(255, 255, 255, 0.08); color: var(--text); border-color: rgba(255,255,255,0.15); }

        .main-content { flex: 1; display: flex; flex-direction: column; overflow: hidden; }
        header { padding: 20px 40px; border-bottom: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center; background: rgba(9, 9, 11, 0.4); backdrop-filter: blur(12px); -webkit-backdrop-filter: blur(12px); z-index: 10; }
        
        .content-body { padding: 40px; overflow-y: auto; flex: 1; }

        .breadcrumb { display: flex; align-items: center; gap: 8px; color: var(--text-muted); font-size: 0.9rem; margin-bottom: 24px; }
        .breadcrumb a { color: var(--text-muted); text-decoration: none; transition: 0.2s; }
        .breadcrumb a:hover { color: var(--accent); }

        .grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 20px; }
        .item-card { background: var(--card); border: 1px solid var(--border); border-radius: 16px; padding: 16px; transition: 0.3s; position: relative; display: flex; flex-direction: column; gap: 12px; backdrop-filter: blur(10px); -webkit-backdrop-filter: blur(10px); }
        .item-card:hover { border-color: rgba(255, 255, 255, 0.2); transform: translateY(-4px); background: rgba(255, 255, 255, 0.08); box-shadow: 0 20px 40px -15px rgba(0,0,0,0.5); }
        
        .icon-wrapper { height: 120px; background: var(--bg); border-radius: 8px; display: flex; align-items: center; justify-content: center; overflow: hidden; border: 1px solid var(--border); }
        .preview-img { width: 100%; height: 100%; object-fit: cover; }
        
        .item-meta { display: flex; flex-direction: column; }
        .item-name { font-size: 0.875rem; font-weight: 500; color: var(--text); text-decoration: none; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
        .item-info { font-size: 0.75rem; color: var(--text-muted); margin-top: 2px; }

        .btn-group { display: flex; gap: 8px; }
        .btn { padding: 8px 16px; border-radius: 8px; font-size: 0.875rem; font-weight: 500; cursor: pointer; display: flex; align-items: center; gap: 8px; transition: 0.2s; border: 1px solid transparent; }
        .btn-primary { background: var(--accent); color: white; }
        .btn-primary:hover { opacity: 0.9; }
        .btn-ghost { background: transparent; color: var(--text-muted); border-color: var(--border); }
        .btn-ghost:hover { background: var(--card); color: var(--text); }
        
        .delete-btn { position: absolute; top: 8px; right: 8px; opacity: 0; background: rgba(239, 68, 68, 0.1); color: #ef4444; border: 1px solid rgba(239, 68, 68, 0.2); border-radius: 6px; width: 28px; height: 28px; display: flex; align-items: center; justify-content: center; cursor: pointer; transition: 0.2s; }
        .delete-btn:hover { background: #ef4444; color: white; }
        .item-card:hover .delete-btn { opacity: 1; }

        .storage-card { margin-top: auto; padding: 16px; background: rgba(255, 255, 255, 0.05); border-radius: 12px; border: 1px solid var(--border); backdrop-filter: blur(8px); -webkit-backdrop-filter: blur(8px); }
        .storage-header { display: flex; justify-content: space-between; font-size: 0.75rem; margin-bottom: 8px; color: var(--text-muted); }
        .progress-bg { height: 6px; background: var(--bg); border-radius: 3px; overflow: hidden; }
        .progress-fill { height: 100%; background: var(--accent); border-radius: 3px; }

        .modal-overlay { display: none; position: fixed; inset: 0; background: rgba(0,0,0,0.7); backdrop-filter: blur(10px); -webkit-backdrop-filter: blur(10px); z-index: 50; align-items: center; justify-content: center; }
        .form-card { background: rgba(18, 10, 25, 0.9); border: 1px solid var(--border); padding: 24px; border-radius: 20px; width: 400px; box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.6); backdrop-filter: blur(20px); }
        input[type="text"], input[type="file"] { width: 100%; background: var(--bg); border: 1px solid var(--border); color: white; padding: 10px; border-radius: 8px; margin: 12px 0; outline: none; }
        input[type="text"]:focus { border-color: var(--accent); }

        #previewModal { display: none; position: fixed; inset: 0; background: rgba(0,0,0,0.95); backdrop-filter: blur(20px); -webkit-backdrop-filter: blur(20px); z-index: 100; align-items: center; justify-content: center; }
        #modalImg { max-width: 90%; max-height: 90%; border-radius: 8px; box-shadow: 0 0 50px rgba(0,0,0,0.5); }
    </style>
    <script>
        function showModal(id) { document.getElementById(id).style.display = 'flex'; }
        function hideModal(id) { document.getElementById(id).style.display = 'none'; }
        function openPreview(src) { 
            document.getElementById('modalImg').src = src;
            showModal('previewModal');
        }
    </script>
</head>
<body>
    <aside class="sidebar">
        <div class="logo">
            <i data-lucide="cloud"></i>
            <span>SharkNAS</span>
        </div>
        <nav>
            <a href="/storage/" class="nav-item {{if or (eq .CurrentPath "/") (eq .CurrentPath "Root")}}active{{end}}">
                <i data-lucide="folder"></i> All Files
            </a>
            <a href="/storage/recent" class="nav-item {{if eq .CurrentPath "Recent"}}active{{end}}"><i data-lucide="clock"></i> Recent</a>
            <a href="/storage/favorites" class="nav-item {{if eq .CurrentPath "Favorites"}}active{{end}}"><i data-lucide="star"></i> Favorites</a>
            <a href="/storage/trash/" class="nav-item {{if eq .CurrentPath ".trash"}}active{{end}}"><i data-lucide="trash-2"></i> Trash</a>
        </nav>
        <div class="storage-card">
            <div class="storage-header">
                <span>Used Storage</span>
                <span>{{usagePercent .UsedStorage .MaxStorage | printf "%.1f"}}%</span>
            </div>
            <div class="progress-bg">
                <div class="progress-fill" style="width: {{usagePercent .UsedStorage .MaxStorage}}%"></div>
            </div>
            <div style="font-size: 11px; color: var(--text-muted); margin-top: 8px;">
                {{formatSize .UsedStorage}} of {{formatSize .MaxStorage}}
            </div>
        </div>
    </aside>
    <main class="main-content">
        <header>
            <form action="/storage/search" method="GET" style="display:flex; align-items:center; gap:20px; flex:1">
                 <i data-lucide="search" style="color: var(--text-muted); width: 18px;"></i>
                 <input type="text" name="q" placeholder="Search your files..." style="background:transparent; border:none; color:white; width: 300px; margin:0">
            </form>
            <div class="btn-group">
                <button onclick="showModal('mkdirModal')" class="btn btn-ghost"><i data-lucide="folder-plus" style="width:18px"></i> New Folder</button>
                <button onclick="showModal('uploadModal')" class="btn btn-primary"><i data-lucide="upload" style="width:18px"></i> Upload</button>
            </div>
        </header>
        <div class="content-body">
            <div class="breadcrumb">
                <a href="/storage/"><i data-lucide="home" style="width:16px"></i></a>
                <span>/</span>
                {{if and (ne .CurrentPath "/") (ne .CurrentPath "Recent") (ne .CurrentPath "Favorites")}}
                    <span>{{.CurrentPath}}</span>
                {{else if eq .CurrentPath "Recent"}}
                    <span>Recent</span>
                {{else if eq .CurrentPath "Favorites"}}
                    <span>Favorites</span>
                {{else}}<span>Root</span>{{end}}
            </div>
            <div class="grid">
                {{range .Files}}
                <div class="item-card">
                    <form action="/storage/delete" method="POST" onsubmit="return confirm('Delete {{.Name}}?')">
                        <input type="hidden" name="path" value="{{.Path}}">
                        <button type="submit" class="delete-btn"><i data-lucide="trash-2" style="width:16px"></i></button>
                    </form>
                    <form action="/storage/favorite/toggle" method="POST" style="position: absolute; top: 8px; left: 8px; z-index: 10;">
                        <input type="hidden" name="path" value="{{.Path}}">
                        <button type="submit" style="background:none; border:none; cursor:pointer; padding:0; color: {{if .IsFavorite}}#eab308{{else}}#71717a{{end}}">
                            <i data-lucide="star" style="width:16px; {{if .IsFavorite}}fill:#eab308{{end}}"></i>
                        </button>
                    </form>

                    <div class="icon-wrapper">
                    {{if .IsDir}}
                        <i data-lucide="folder" style="width: 48px; height: 48px; color: #eab308; fill: rgba(234, 179, 8, 0.1);"></i>
                    {{else}}
                        {{if (isImage .Name)}}
                            <img src="/storage{{.Path}}" class="preview-img" onclick="openPreview(this.src)">
                        {{else}}
                            <i data-lucide="file" style="width: 48px; height: 48px; color: #71717a;"></i>
                        {{end}}
                    {{end}}
                    </div>
                    <div class="item-meta">
                        <a href="/storage{{.Path}}" class="item-name" title="{{.Name}}">{{.Name}}</a>
                        <span class="item-info">{{if .IsDir}}Folder{{else}}{{formatSize .Size}}{{end}}</span>
                    </div>
                </div>
                {{else}}
                <div style="grid-column: 1/-1; text-align: center; padding: 100px; color: var(--text-muted);">
                    <i data-lucide="inbox" style="width: 48px; height: 48px; margin-bottom: 16px;"></i>
                    <p>No files found in this directory</p>
                </div>
                {{end}}
            </div>
        </div>
    </main>
    <div id="uploadModal" class="modal-overlay" onclick="hideModal('uploadModal')">
        <div class="form-card" onclick="event.stopPropagation()">
            <h3 style="margin-top:0">Upload Files</h3>
            <form action="/storage/upload?dir={{.CurrentPath}}" method="POST" enctype="multipart/form-data">
                <input type="file" name="uploadFile" multiple required>
                <button type="submit" class="btn btn-primary" style="width:100%; justify-content:center">Start Upload</button>
            </form>
        </div>
    </div>
    <div id="mkdirModal" class="modal-overlay" onclick="hideModal('mkdirModal')">
        <div class="form-card" onclick="event.stopPropagation()">
            <h3 style="margin-top:0">New Folder</h3>
            <form action="/storage/mkdir" method="POST">
                <input type="hidden" name="dir" value="{{.CurrentPath}}">
                <input type="text" name="name" placeholder="Enter folder name..." required autofocus>
                <button type="submit" class="btn btn-primary" style="width:100%; justify-content:center">Create Folder</button>
            </form>
        </div>
    </div>
    <div id="previewModal" class="modal-overlay" onclick="hideModal('previewModal')">
        <img id="modalImg" src="" onclick="event.stopPropagation()">
    </div>
    <script>lucide.createIcons();</script>
</body>
</html>
`
