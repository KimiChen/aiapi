aipick 是 systemd 二进制部署形态：应用监听 0.0.0.0:8080，公网入口由多台反向代理服务器转发到该端口。
aipick 服务名: sub2api
aipick 安装目录: /opt/sub2api
aipick 当前二进制: /opt/sub2api/sub2api -> /opt/sub2api/releases/20260701-001157/sub2api
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

aipick 回滚方式:
 - 确认要回滚的 release 目录，例如 `/opt/sub2api/releases/20260628-173811`
 - 执行 `ln -sfn /opt/sub2api/releases/20260628-173811/sub2api /opt/sub2api/sub2api`
 - 执行 `systemctl restart sub2api`
 - 验证 `curl -fsS http://127.0.0.1:8080/status` 和公网反向代理入口
