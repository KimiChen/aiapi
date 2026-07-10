aipick 是 systemd 二进制部署形态：应用监听 0.0.0.0:8080，公网入口由多台反向代理服务器转发到该端口。
aipick 服务名: sub2api
aipick 安装目录: /opt/sub2api
aipick 当前二进制: /opt/sub2api/sub2api -> /opt/sub2api/releases/20260710-161211/sub2api
aipick 资源目录: /opt/sub2api/resources
aipick 配置目录: /etc/sub2api
aipick 查看服务状态: systemctl status sub2api --no-pager -l
aipick 查看应用日志: journalctl -u sub2api -f
aipick 本机健康检查: curl -fsS http://127.0.0.1:8080/status

aipick 部署方式:
 - 本地先执行 `pnpm --dir frontend build`
 - 在 `backend/` 目录执行 linux/amd64 嵌入前端资源构建：
   - `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags embed -ldflags="-s -w -X main.Version=$(tr -d '\r\n' < cmd/server/VERSION)$(tr -d '\r\n' < cmd/server/SUB_VERSION) -X main.BuildType=release" -trimpath -o bin/server ./cmd/server`
 - 上传 `backend/bin/server` 到 `/opt/sub2api/releases/<timestamp>/sub2api`
 - 同步 `backend/resources/` 到 `/opt/sub2api/resources/`
 - 切换 `/opt/sub2api/sub2api` symlink 后执行 `systemctl restart sub2api`

aipick 8080 访问控制:
 - systemd 使用 `ExecStartPre=+/usr/local/sbin/sub2api-8080-firewall` 启动前刷新规则
 - 允许本机 loopback 访问 8080
 - 允许 WireGuard 内网 `10.111.0.0/24` 访问 8080
 - 允许反向代理公网 IP 访问 8080: `155.94.192.174`, `64.186.228.84`, `69.63.200.39`, `154.26.183.236`, `23.148.204.176`, `103.117.102.195`
 - 其它 IPv4 来源访问 8080 会被 DROP，IPv6 访问 8080 会被 DROP

aipick 部署记录:
 - 部署时间: 2026-07-10 16:12 Asia/Shanghai
 - Git HEAD: 0a64147b7b83e8bfc0a129bea200c39a2e7e012a
 - 版本号: 0.1.150.kim
 - 发布目录: /opt/sub2api/releases/20260710-161211
 - 当前二进制 SHA256: 8a6db6452f44bb4c19ebac17d9b7d3c571435980d061000e207f71cac210b006
 - 上一版发布目录: /opt/sub2api/releases/20260710-002349
 - 上一版二进制 SHA256: 82316ad4e709bff55408068dd9d4288119e9686e3c036c9e17471bae1cc413d9
 - 部署动作: 确认本地 HEAD 同时包含 `origin/main` 与 `upstream/main` 最新提交，构建前端，编译 linux/amd64 `-tags embed` 后端，压缩部署包并按 2 MiB 分块上传，远端合并校验归档 SHA256 `6dd9701a430074eee0bfc7a33ea484c35face5d6c9416eaedf238e8bde231491` 和二进制 SHA256，同步 resources，切换 symlink，重启 `sub2api`
 - 验证结果: 本机 `http://127.0.0.1:8080/status` 返回 `{"status":"perfectly nice"}`；所有反向代理 HTTPS 入口 `/status` 与 `/` 均返回 200；`systemctl show sub2api` 显示 ActiveState=active、SubState=running、NRestarts=0
 - 备注: 部署前远端运行版本为 `0.1.149.kim`，指向 `/opt/sub2api/releases/20260710-002349/sub2api`。重启后首次 8080 探测命中启动窗口，随后本机与所有反向代理入口验证正常。

 - 部署时间: 2026-07-10 00:23 Asia/Shanghai
 - Git HEAD: 4b8c9e3b983a24653fa94c27885162ea734502e5
 - 版本号: 0.1.149.kim
 - 发布目录: /opt/sub2api/releases/20260710-002349
 - 当前二进制 SHA256: 82316ad4e709bff55408068dd9d4288119e9686e3c036c9e17471bae1cc413d9
 - 上一版发布目录: /opt/sub2api/releases/20260708-073913
 - 上一版二进制 SHA256: e8065212c3601344dadbd02ed3de36c7383e43e254521d1c024a115af1f23ed2
 - 部署动作: 确认本地 HEAD 同时包含 `origin/main` 与 `upstream/main` 最新提交，构建前端，编译 linux/amd64 `-tags embed` 后端，压缩部署包并按 2 MiB 分块上传，远端合并校验归档 SHA256 `6bd60160ddab31cb4092c6ded68eae15ce87106d455470f302222b4e39e2eb34` 和二进制 SHA256，同步 resources，切换 symlink，重启 `sub2api`
 - 验证结果: 本机 `http://127.0.0.1:8080/status` 返回 `{"status":"perfectly nice"}`；所有反向代理 HTTPS 入口 `/status` 与 `/` 均返回 200；`systemctl show sub2api` 显示 ActiveState=active、SubState=running、NRestarts=0
 - 备注: 部署前远端实际运行版本为 `0.1.146.kim`，指向 `/opt/sub2api/releases/20260708-073913/sub2api`；该次发布未在本地 `aipick.md` 顶部记录，本次按远端实际状态作为上一版。重启后首次 8080 探测命中启动窗口，随后本机与所有反向代理入口验证正常。

 - 部署时间: 2026-07-07 01:14 Asia/Shanghai
 - Git HEAD: 9bb775198b01565e6f7a8419938f5b0bb4853386
 - 版本号: 0.1.145.kim
 - 发布目录: /opt/sub2api/releases/20260707-011442
 - 当前二进制 SHA256: e7fc8345faf1ecea9081dcaf1e92b43ff9480c0b3f259c0304edd365a58299bb
 - 上一版发布目录: /opt/sub2api/releases/20260706-003502
 - 上一版二进制 SHA256: 3cc44cfa54b6334650d0341c2876f0470879f3b817e6a1889adc71248fe367f3
 - 部署动作: 确认本地 HEAD 同时包含 `origin/main` 与 `upstream/main` 最新提交，构建前端，编译 linux/amd64 `-tags embed` 后端，压缩部署包并按 2 MiB 分块上传，远端合并校验归档 SHA256 `589600c868a62d5d5bb41fea8ff5157e5c49206a559178520c900451f88e3d0b` 和二进制 SHA256，同步 resources，切换 symlink，重启 `sub2api`
 - 验证结果: 本机 `http://127.0.0.1:8080/status` 返回 `{"status":"perfectly nice"}`；所有反向代理 HTTPS 入口 `/status` 与 `/` 均返回 200；`systemctl show sub2api` 显示 ActiveState=active、SubState=running、NRestarts=0
 - 备注: 部署前远端运行版本为 `0.1.144.kim`，指向 `/opt/sub2api/releases/20260706-003502/sub2api`。重启后首次 8080 探测命中启动窗口，随后本机与所有反向代理入口验证正常。

 - 部署时间: 2026-07-06 00:39 Asia/Shanghai
 - Git HEAD: edde82f159441e1cfc6ec5482b41a7f26275a4fb
 - 版本号: 0.1.144.kim
 - 发布目录: /opt/sub2api/releases/20260706-003502
 - 当前二进制 SHA256: 3cc44cfa54b6334650d0341c2876f0470879f3b817e6a1889adc71248fe367f3
 - 上一版发布目录: /opt/sub2api/releases/20260703-003814
 - 上一版二进制 SHA256: 0406af918828b2f364a2e8d05b1a8c9657745f257c73b7e2decc9fa8a0de0237
 - 部署动作: 构建前端，编译 linux/amd64 `-tags embed` 后端，压缩部署包并按 2 MiB 分块上传，远端合并校验归档 SHA256 `7e0c3d5bb041d6128a1b81ad8565082fbd9bc5c4147d83847e84607595b97f85` 和二进制 SHA256，同步 resources，切换 symlink，重启 `sub2api`
 - 验证结果: 本机 `http://127.0.0.1:8080/status` 返回 `{"status":"perfectly nice"}`；所有反向代理 HTTPS 入口 `/status` 与 `/` 均返回 200；`systemctl show sub2api` 显示 ActiveState=active、SubState=running、NRestarts=0
 - 备注: 部署前远端实际运行版本为 `0.1.143.kim`，指向 `/opt/sub2api/releases/20260703-003814/sub2api`；该次发布未在本地 `aipick.md` 顶部记录，本次按远端实际状态作为上一版。公网验证过程中个别入口出现一次瞬时 TLS/SSL 连接错误，重试后均返回 200。

 - 部署时间: 2026-07-03 00:39 Asia/Shanghai
 - Git HEAD: c23964ab8475bda2be6f52e4f99a23a6a149497e
 - 版本号: 0.1.143.kim
 - 发布目录: /opt/sub2api/releases/20260703-003814
 - 当前二进制 SHA256: 0406af918828b2f364a2e8d05b1a8c9657745f257c73b7e2decc9fa8a0de0237
 - 上一版发布目录: /opt/sub2api/releases/20260702-012210
 - 上一版二进制 SHA256: 7e0fc08e4be761c5643607399817108072af11a597c2eca360daab39c592f823
 - 部署动作: 构建前端，编译 linux/amd64 `-tags embed` 后端，压缩部署包并按 2 MiB 分块上传，远端合并校验归档 SHA256 和二进制 SHA256，同步 resources，切换 symlink，重启 `sub2api`
 - 验证结果: 本机 `http://127.0.0.1:8080/status` 返回 `{"status":"perfectly nice"}`；所有反向代理 HTTPS 入口 `/status` 均返回 200；线上静态资源已包含用户 `/usage` 的 `IP Address`、`ip_address` 与 `user-usage-hidden-columns` 标记；`systemctl show sub2api` 显示 ActiveState=active、SubState=running、NRestarts=0
 - 备注: 本次部署当前 HEAD，包含用户使用记录 IP 地址列和 IP 地区查询能力；源码中该功能已在上一部署之后合入，因此本轮主要完成线上发布。

 - 部署时间: 2026-07-02 01:22 Asia/Shanghai
 - Git HEAD: 612359b6685dca071de2c68ce9c733a6a7ced90a
 - 版本号: 0.1.142.kim
 - 发布目录: /opt/sub2api/releases/20260702-012210
 - 当前二进制 SHA256: 7e0fc08e4be761c5643607399817108072af11a597c2eca360daab39c592f823
 - 上一版发布目录: /opt/sub2api/releases/20260701-001157
 - 上一版二进制 SHA256: fcd7a99c964549a7de8864a8a8d6f1c008f212de9285db7d3cce07a00c0b5809
 - 部署动作: 构建前端，编译 linux/amd64 `-tags embed` 后端，压缩部署包并按 2 MiB 分块上传，远端合并校验归档 SHA256 和二进制 SHA256，同步 resources，切换 symlink，重启 `sub2api`
 - 验证结果: 本机 `http://127.0.0.1:8080/status` 返回 `{"status":"perfectly nice"}`；所有反向代理 HTTPS 入口 `/status` 与 `/` 均返回 200；`systemctl show sub2api` 显示 ActiveState=active、SubState=running、NRestarts=0
 - 备注: 公网 SSH 直接上传大文件和压缩流上传中途被连接重置，期间 symlink 未切换、旧版服务保持运行；最终使用分块上传完成发布。重启后立即探测 8080 曾命中启动窗口，随后服务正常监听并通过本机和所有反向代理入口验证。

 - 部署时间: 2026-07-01 00:11 Asia/Shanghai
 - Git HEAD: 0358b7f15a6ce75a7077609e7c98cc20d00e04aa
 - 版本号: 0.1.141.kim
 - 发布目录: /opt/sub2api/releases/20260701-001157
 - 当前二进制 SHA256: fcd7a99c964549a7de8864a8a8d6f1c008f212de9285db7d3cce07a00c0b5809
 - 上一版发布目录: /opt/sub2api/releases/20260628-173811
 - 上一版二进制 SHA256: 3c39c242a3922bf58a7bdbe8b72b64104bca7a814347bd6706bc8b6a6e48f8ea
 - 部署动作: 构建前端，编译 linux/amd64 `-tags embed` 后端，上传新二进制和 resources，切换 symlink，更新 8080 防火墙白名单，重启 `sub2api`
 - 验证结果: 本机 `http://127.0.0.1:8080/status` 返回 `{"status":"perfectly nice"}`；所有反向代理 HTTPS 入口 `/status` 与 `/` 均返回 200；`systemctl show sub2api` 显示 ActiveState=active、SubState=running、NRestarts=0
 - 备注: 尝试 SSH 到部分反向代理服务器做源端直连探测时，22 端口会话被对端关闭；本次以 aipick 端 iptables 白名单和各反代入口 200 作为验证依据

aipick 反向代理真实客户端 IP 修复记录:
 - 修复时间: 2026-07-01 10:38 Asia/Shanghai
 - 问题现象: aipick/sub2api 使用记录里部分请求的 `client_ip` 记录为 `127.0.0.1`
 - 原因: 非主站反向代理服务器的 nginx stream 层把公网 443 按 SNI 转发到本机 `127.0.0.1:8443`，内层 HTTPS vhost 看到的 `$remote_addr` 是本机 loopback，并继续把 `X-Real-IP: 127.0.0.1` 传给 aipick；sub2api 的客户端 IP 解析顺序会优先采用 `X-Real-IP`
 - 修复方式: 在非主站反向代理服务器上用 nginx PROXY protocol 把 stream 层看到的真实公网客户端地址传给内层 HTTPS vhost，再由 vhost 继续传给 aipick/sub2api
 - 修改的远端 stream 配置路径:
   - `64.186.228.84:/www/server/panel/vhost/nginx/tcp/remnawave-sni.conf`
   - `69.63.200.39:/www/server/panel/vhost/nginx/tcp/remnawave-sni.conf`
   - `154.26.183.236:/www/server/panel/vhost/nginx/tcp/remnawave-sni.conf`
   - `23.148.204.176:/www/server/panel/vhost/nginx/tcp/remnawave-sni.conf`
   - `103.117.102.195:/www/server/panel/vhost/nginx/tcp/remnawave-sni.conf`
 - stream 配置修改内容:
   - 新增 `upstream aihub_proxy_backend { server 127.0.0.1:9443; }`
   - 保留 `aws.amazon.com` 到 `reality_backend` 的既有分流
   - 将 aipick 相关 SNI 单独分流到 `aihub_proxy_backend`
   - 在 aipick 相关 SNI 的 stream `server` 块里增加 `proxy_protocol on;`
   - 保留默认流量到 `web_backend`
 - 修改的远端 HTTPS vhost 配置路径:
   - `64.186.228.84:/www/server/panel/vhost/nginx/www.aihub.pick.art.conf`
   - `69.63.200.39:/www/server/panel/vhost/nginx/uscu.aihub.pick.art.conf`
   - `154.26.183.236:/www/server/panel/vhost/nginx/usct.aihub.pick.art.conf`
   - `23.148.204.176:/www/server/panel/vhost/nginx/usall.aihub.pick.art.conf`
   - `103.117.102.195:/www/server/panel/vhost/nginx/jpct.aihub.pick.art.conf`
 - HTTPS vhost 配置修改内容:
   - 将 `listen 127.0.0.1:8443 ssl;` 改为 `listen 127.0.0.1:9443 ssl proxy_protocol;`
   - 新增 `real_ip_header proxy_protocol;`
   - 新增 `set_real_ip_from 127.0.0.1;`
   - 其它代理到 aipick 8080 的规则保持不变
 - 远端备份文件:
   - `64.186.228.84`: `remnawave-sni.conf.bak.codex-proxyproto-20260701023831`, `www.aihub.pick.art.conf.bak.codex-proxyproto-20260701023831`
   - `69.63.200.39`: `remnawave-sni.conf.bak.codex-proxyproto-20260701023833`, `uscu.aihub.pick.art.conf.bak.codex-proxyproto-20260701023833`
   - `154.26.183.236`: `remnawave-sni.conf.bak.codex-proxyproto-20260701023835`, `usct.aihub.pick.art.conf.bak.codex-proxyproto-20260701023835`
   - `23.148.204.176`: `remnawave-sni.conf.bak.codex-proxyproto-20260701060837`, `usall.aihub.pick.art.conf.bak.codex-proxyproto-20260701060837`
   - `103.117.102.195`: `remnawave-sni.conf.bak.codex-proxyproto-20260701023838`, `jpct.aihub.pick.art.conf.bak.codex-proxyproto-20260701023838`
 - 验证结果:
   - 五台非主站反向代理服务器执行 `nginx -t -q -c /www/server/nginx/conf/nginx.conf` 均通过
   - 五台非主站反向代理服务器均已 reload nginx
   - 所有非主站 HTTPS 入口 `/status` 返回 200
   - 通过各非主站 HTTPS 入口访问探测路径后，aipick 日志中的 `client_ip` 已记录为真实公网出口 IP，不再是 `127.0.0.1`

aipick 回滚方式:
 - 确认要回滚的 release 目录，例如 `/opt/sub2api/releases/20260628-173811`
 - 执行 `ln -sfn /opt/sub2api/releases/20260628-173811/sub2api /opt/sub2api/sub2api`
 - 执行 `systemctl restart sub2api`
 - 验证 `curl -fsS http://127.0.0.1:8080/status` 和公网反向代理入口
