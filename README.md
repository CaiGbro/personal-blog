

---

# Personal Blog - 个人博客系统

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue.svg)](https://golang.org/)
[![Vue.js Version](https://img.shields.io/badge/Vue.js-3.x-brightgreen.svg)](https://vuejs.org/)

**一个现代化、功能丰富的个人博客系统，采用 Go + Vue.js 3 构建，专为高性能、强交互和便捷部署而设计。**

这不仅仅是一个博客，更是一个集文章发布、视频画廊、互动社区于一体的综合性内容平台。

![项目截图](./backend/static/readme/项目截图.png)
![项目截图](./backend/static/readme/项目截图2.png)
![项目截图](./backend/static/readme/项目截图3.png)
![项目截图](./backend/static/readme/项目截图4.png)
![项目截图](./backend/static/readme/项目截图5.png)
![项目截图](./backend/static/readme/项目截图6.png)
![项目截图](./backend/static/readme/项目截图7.jpg)
![项目截图](./backend/static/readme/项目截图8.jpg)
![项目截图](./backend/static/readme/项目截图9.jpg)


---

## ✨ 功能亮点

本项目深度整合了多种媒体格式和互动功能，旨在提供卓越的内容消费和社区交流体验。

*   **📝 多格式文章系统**:
    *   支持 `Markdown`, `DOCX` 文件的自动转换为 `PDF` 格式进行在线预览。
    *   文章列表支持瀑布流无限滚动加载和实时文件名搜索。
    *   为每篇文章提供独立的表情回应和树状评论区。

*   **🎬 沉浸式视频画廊**:
    *   支持视频批量上传，并通过 FFmpeg 自动转换为Web优化格式 (`H.264` + `faststart`)。
    *   自动为视频生成 `JPG` 格式的封面图。
    *   提供类似 TikTok/Shorts 的**全屏沉浸式**上下滑动切换的视频播放器。
    *   同样支持视频的独立表情回应和树状评论区。

*   **💬 全功能的互动社区**:
    *   留言板、文章、视频均支持**树状评论**，可无限层级回复。
    *   评论支持**图片附件**上传（后端进行安全验证）。
    *   支持对内容（文章/视频）和评论进行**表情符号回应** (Reactions)。

*   **🔒 强大的后台管理**:
    *   基于 `JWT` 的管理员登录认证体系。
    *   提供**访客统计**面板，通过图表展示近期独立访客趋势。
    *   **游客模式**一键切换：开启后全站内容可访问，关闭后仅管理员可见，轻松实现站点维护。

*   **🚀 高性能与安全性**:
    *   后端采用高性能的 Go + Gin 框架。
    *   API 接口内置**速率限制**，有效防止恶意请求。
    *   所有用户输入内容（如评论）均经过严格的 **XSS 清理**。
    *   通过 IP 哈希实现对独立访客的**无隐私侵犯追踪**。

*   **🤖 赛博朋克广告系统 (可选开启)**:
    *   一个有趣的可选功能，开启后页面会切换为赛博朋克主题，并展示多种动态广告（GIF动图、浮动文字、背景标语等）。
    *   提供一键关闭功能，为用户提供选择权。

*   **🛠️ 便捷的工具脚本**:
    *   提供一系列 Python 脚本，用于自动化处理媒体资源和项目备份。

## 🛠️ 技术栈

本项目采用了前后端分离的架构，并整合了多种现代化技术。

#### **后端 (Backend)**
*   **语言**: Go 1.22+
*   **Web 框架**: Gin
*   **数据库**: PostgreSQL
*   **数据库驱动**: pgx
*   **认证**: JSON Web Tokens (JWT)
*   **环境变量**: godotenv
*   **邮件服务**: SMTP (用于发送登录验证码)

#### **前端 (Frontend)**
*   **框架**: Vue 3 (Composition API & `<script setup>`)
*   **构建工具**: Vite
*   **路由**: Vue Router
*   **HTTP客户端**: Axios
*   **UI & 动效**:
    *   `Chart.js` 用于后台统计图表
    *   `vue-pdf-embed` 用于 PDF 预览
    *   `docx-preview` 用于 DOCX 预览
    *   `marked` 用于 Markdown 渲染

#### **运维与部署 (DevOps & Infrastructure)**
*   **反向代理**: Nginx (推荐) 或 Caddy
*   **安全隧道**: Cloudflare Tunnel (用于将内网服务安全暴露到公网)
*   **媒体处理**: FFmpeg (视频转码和缩略图生成)
*   **文档转换**: WeasyPrint, Microsoft Word (通过 Python 的 `comtypes` 库调用)
*   **进程管理**: Systemd (推荐用于后台服务的持久化运行)

## 🚀 快速开始 (本地开发)

请确保你的开发环境已安装以下软件: **Go (1.22+), Node.js (18+), pnpm (或 npm/yarn), PostgreSQL, Python 3, FFmpeg**。

#### 1. 克隆项目
```bash
git clone https://github.com/CaiGbro/personal-blog.git
cd personal-blog
```

#### 2. 配置并运行后端
```bash
# 进入后端目录
cd backend

# 1. 初始化数据库
#    - 登录 PostgreSQL 并创建一个新的数据库，例如 'personal_blog'
#    - 执行 /database/init.sql 脚本来创建所需的表结构

# 2. 创建并配置 .env 文件
#    - 复制 .env.example (如果提供) 或手动创建 .env 文件
#    - 填入以下内容，并替换为你自己的配置
DATABASE_URL="postgresql://postgres:your_password@localhost:5432/personal_blog?sslmode=disable"
JWT_SECRET_KEY="use_a_very_long_and_random_string_here"
SMTP_PASSWORD="your_smtp_password_for_163_mail"

# 3. 安装 Go 依赖
go mod tidy

# 4. 运行后端服务
go run main.go
```
> 后端服务将运行在 `http://localhost:8081`。

#### 3. 配置并运行前端
```bash
# 回到项目根目录，然后进入前端目录
cd ../frontend

# 1. 安装 Node.js 依赖
pnpm install

# 2. 运行前端开发服务器
pnpm dev
```
> 前端开发服务器将运行在 `http://localhost:5173`，并已配置好代理，会自动将 `/api` 和 `/static` 请求转发到 `8081` 端口。

现在，你可以通过浏览器访问 `http://localhost:5173` 来查看你的博客了！

## 🚢 部署指南 (生产环境)

以下是一个基于 **Ubuntu/Debian + Nginx + Systemd + Cloudflare Tunnel** 的典型部署流程。

#### 1. 服务器准备
*   确保服务器上已安装 Go, Nginx, PostgreSQL, FFmpeg, Python 及其依赖 (`weasyprint` 等)。
*   将整个项目文件上传到服务器，例如 `/prodct/personal-blog/personal-blog_prduct`。

#### 2. 部署后端
*   **构建 Go 应用**:
    ```bash
    cd /prodct/personal-blog/personal-blog_prduct/backend
    go build -o blog-backend .
    ```
    这会生成一个名为 `blog-backend` 的二进制可执行文件。
*   **配置 `.env`**: 确保 `backend/.env` 文件中的数据库连接信息对于生产环境是正确的。
*   **使用 Systemd 持久化运行**:
    *   创建一个 `systemd` 服务文件: `sudo nano /etc/systemd/system/blog-backend.service`
    *   填入以下内容 (请根据你的实际路径修改 `WorkingDirectory` 和 `ExecStart`):
        ```ini
        [Unit]
        Description=Personal Blog Backend Service
        After=network.target postgresql.service

        [Service]
        User=your_user # 推荐使用非 root 用户
        Group=your_group
        WorkingDirectory=/prodct/personal-blog/personal-blog_prduct/backend
        ExecStart=/prodct/personal-blog/personal-blog_prduct/backend/blog-backend
        Restart=on-failure
        RestartSec=5s
        EnvironmentFile=/prodct/personal-blog/personal-blog_prduct/backend/.env

        [Install]
        WantedBy=multi-user.target
        ```
    *   启动并设置开机自启:
        ```bash
        sudo systemctl daemon-reload
        sudo systemctl start blog-backend
        sudo systemctl enable blog-backend
        sudo systemctl status blog-backend # 检查服务状态
        ```

#### 3. 部署前端
*   **构建前端静态文件**:
    ```bash
    cd /prodct/personal-blog/personal-blog_prduct/frontend
    pnpm install
    pnpm build
    ```
    这会在 `frontend/dist` 目录下生成所有需要的 HTML, CSS, JS 文件。

#### 4. 配置 Nginx
*   创建一个新的 Nginx 配置文件: `sudo nano /etc/nginx/sites-available/caigbro.top`
*   填入以下配置 (请将域名和文件路径替换为你自己的):
    ```nginx
    server {
        listen 8080; # Nginx 监听一个本地端口，稍后由 Cloudflare Tunnel 转发
        server_name caigbro.top www.caigbro.top;

        # API 和后端静态资源请求全部转发给 Go 服务
        location /api/ {
            proxy_pass http://127.0.0.1:8081; # 注意：末尾没有斜杠
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }

        location /static/ {
            proxy_pass http://127.0.0.1:8081; # 注意：末尾没有斜杠
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }

        # 前端静态文件由 Nginx 直接提供
        location / {
            root /prodct/personal-blog/personal-blog_prduct/frontend/dist;
            try_files $uri $uri/ /index.html;
        }

        # Gzip, 日志等配置...
        gzip on;
        gzip_types text/plain text/css application/javascript application/json;
        access_log /var/log/nginx/caigbro.access.log;
        error_log /var/log/nginx/caigbro.error.log;
    }
    ```
*   启用该配置并重启 Nginx:
    ```bash
    sudo ln -s /etc/nginx/sites-available/caigbro.top /etc/nginx/sites-enabled/
    sudo nginx -t # 测试配置是否正确
    sudo systemctl restart nginx
    ```

#### 5. 配置并运行 Cloudflare Tunnel
*   按照 [Cloudflare官方文档](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/install-and-setup/tunnel-guide/) 安装 `cloudflared`。
*   登录并创建一个隧道。
*   配置 `config.yml` 文件 (通常在 `/root/.cloudflared/` 或 `~/.cloudflared/`):
    ```yaml
    tunnel: <Your-Tunnel-UUID>
    credentials-file: /root/.cloudflared/<Your-Tunnel-UUID>.json
    ingress:
      - hostname: caigbro.top
        service: http://localhost:8080 # 转发到 Nginx 监听的端口
      - hostname: www.caigbro.top
        service: http://localhost:8080
      - service: http_status:404
    ```
*   将 `cloudflared` 作为服务运行，以便它能持久工作。

## ⚙️ 工具脚本说明

项目 `backend` 目录中包含多个实用的 Python 脚本，用于自动化维护工作。

*   `convert_videos.py`: 扫描 `static/video` 目录，将非 H.264 编码或未优化的 MP4 文件转换为 web 友好格式，并创建原始文件备份。
*   `convert_writings_to_pdf.py`: 扫描 `static/writings` 目录，将新增或更新的 `.md` 和 `.docx` 文件转换为 `.pdf` 文件。
*   `generate_thumbnails.py`: 扫描 `static/video` 目录，为所有没有封面图的 `.mp4` 视频生成 `.jpg` 格式的缩略图。
*   `rename_videos.py`: 批量将 `static/video` 目录下的视频文件重命名为唯一的 UUID 名称，便于管理和隐藏原始文件名。
*   `backup.py`: (位于项目根目录) 创建整个项目的备份副本，并自动排除 `node_modules`, `env_web` 等无需备份的文件夹。

## 🗺️ 未来计划

*   [ ] 实现文章标签 (Tags) 分类功能。
*   [ ] 增加全局站内搜索。
*   [ ] 后台管理面板增加更多数据统计维度。
*   [ ] 优化移动端视频播放器的手势操作。

## 🤝 贡献

欢迎任何形式的贡献！如果你有好的想法或发现了 Bug，请随时提交 Pull Request 或创建 Issue。

1.  Fork 本仓库
2.  创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3.  提交你的更改 (`git commit -m 'Add some AmazingFeature'`)
4.  推送到分支 (`git push origin feature/AmazingFeature`)
5.  打开一个 Pull Request

## 📄 许可证

本项目采用 [MIT License](https://opensource.org/licenses/MIT) 授权。