# aipick 运维规划

本文档只维护当前运行基线、标准操作和保留策略，不再累积逐次部署流水账。每次变更仅更新“当前基线”和末尾的最小变更记录。

## 1. 当前运行基线

| 项目 | 当前值 |
| --- | --- |
| 部署形态 | systemd 二进制部署，应用监听 `0.0.0.0:8080` |
| 服务名 | `sub2api` |
| 安装目录 | `/opt/sub2api` |
| 当前二进制 | `/opt/sub2api/sub2api -> /opt/sub2api/releases/20260719-150347/sub2api` |
| 上一版 release | `/opt/sub2api/releases/20260718-095218` |
| 资源目录 | `/opt/sub2api/resources` |
| 配置目录 | `/etc/sub2api` |
| 当前版本 | `0.1.161.kim` |
| 部署分支 | `origin/kimi` |
| Git HEAD | `682b1e13a6617398769691894b36c9ad0b736b1a` |
| 上游基线 | `upstream/main` `d4b9797ff72024960a035cf22fdd8f213e149169` |
| 数据库迁移 | `schema_migrations=230` |
| 最近验证 | 2026-07-19 15:41 Asia/Shanghai，服务健康，根分区占用 36% |

常用检查：

```bash
systemctl status sub2api --no-pager -l
systemctl show sub2api -p ActiveState -p SubState -p NRestarts -p MainPID
curl -fsS http://127.0.0.1:8080/status
journalctl -u sub2api -n 200 --no-pager
```

公网入口包括 `aihub.pick.art`、`www.aihub.pick.art`、`uscu.aihub.pick.art`、`usct.aihub.pick.art`、`usall.aihub.pick.art` 和 `jpct.aihub.pick.art`。

## 2. 运维目标与保留策略

- 发布必须可回滚：切换前完成数据库备份并保留上一版 release。
- 当前版和上一版 release 始终保留；其它超过 3 天的 release 才可清理。
- 保留最近两个经校验的数据库备份；其它超过 3 天的备份才可清理。
- PostgreSQL、Sub2API、rsyslog 和 systemd journal 的日志统一保留 3 天。
- 根分区达到 75% 时检查增长来源；达到 85% 时优先处理日志、旧 release 和旧备份。
- 清理前必须列出精确对象；不得删除当前 symlink 指向的 release、上一版 release 或最近两个有效数据库备份。

release 和数据库备份目前由人工在部署后或磁盘告警时清理；日志清理由 systemd timer 自动执行。

## 3. 标准发布流程

### 3.1 发布前

1. 确认工作区干净，`origin/kimi` 包含预期的 `upstream/main` 基线，并记录两个提交 SHA。
2. 检查迁移文件、配置变更和回滚兼容性；不可逆迁移必须单独制定数据库恢复方案。
3. 创建 PostgreSQL custom-format 备份，并用 `pg_restore -l` 校验目录；记录文件大小和 SHA256。
4. 如 resources 有变化，先备份线上 `/opt/sub2api/resources`。
5. 确认根分区空间足以同时容纳新 release、备份和构建归档。

### 3.2 构建与部署

```bash
pnpm --dir frontend build

cd backend
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
  -tags embed \
  -ldflags="-s -w -X main.Version=$(tr -d '\r\n' < cmd/server/VERSION)$(tr -d '\r\n' < cmd/server/SUB_VERSION) -X main.BuildType=release" \
  -trimpath \
  -o bin/server ./cmd/server
```

部署时执行以下动作：

1. 校验产物为 linux/amd64 ELF，版本、Git revision 正确且 `vcs.modified=false`。
2. 上传到 `/opt/sub2api/releases/<timestamp>/sub2api`；大文件优先经 `64.186.228.84` 跳板分块上传。
3. 在远端复核归档和二进制 SHA256。
4. 仅在有变化时同步 resources，并保持正确属主。
5. 原子切换 `/opt/sub2api/sub2api` symlink，执行 `systemctl restart sub2api`。
6. 验证通过后删除 staging；保留新旧 release 和数据库备份。

### 3.3 发布验证

- 本机 `/status` 返回 `{"status":"perfectly nice"}`。
- systemd 为 `active/running`，`NRestarts=0`。
- PostgreSQL 可连接，`schema_migrations` 数量和新增对象符合预期。
- 所有公网入口的 `/status` 与 `/` 返回 200。
- 未授权 `/v1/responses` 返回预期 401，避免 API 路由被 SPA 接管。
- 抽查首页及哈希静态资源，确认嵌入资源可访问。
- 检查启动后的 error/panic/fatal 日志；迁移期间允许短暂未就绪，但必须等待健康检查恢复。

## 4. 备份与回滚

数据库备份统一放在 `/opt/sub2api/backups/<timestamp>/`。当前应保留：

- `/opt/sub2api/backups/20260719-150347/sub2api.dump`
- `/opt/sub2api/backups/20260718-095022/sub2api.dump`

二进制回滚：

```bash
ln -sfn /opt/sub2api/releases/<previous>/sub2api /opt/sub2api/sub2api
systemctl restart sub2api
curl -fsS http://127.0.0.1:8080/status
```

回滚后还必须检查 systemd、全部公网入口和数据库兼容性。数据库恢复属于高风险操作，仅在确认迁移与旧二进制不兼容时执行，并在恢复前另存当前数据库。

日志策略变更的配置备份位于 `/root/backup/aipick-log-policy-20260719-153237`。

## 5. 日志与磁盘管理

PostgreSQL 当前策略：

```text
log_statement = none
log_min_duration_statement = 5s
log_min_messages = warning
log_min_error_statement = error
```

这会关闭普通 SQL 全量记录，仅保留执行时间至少 5 秒的慢查询，以及警告、错误和相关错误语句。

三天日志保留由以下配置实施：

- `aipick-log-retention.timer`：每日执行 PostgreSQL、Sub2API、rsyslog 和 journal 清理。
- journald：`MaxRetentionSec=3day`。
- rsyslog/logrotate：`daily`、`rotate 3`、`maxage 3`、`compress`。

检查命令：

```bash
systemctl status aipick-log-retention.timer --no-pager
systemctl show aipick-log-retention.service -p Result -p ExecMainStatus
journalctl --disk-usage
du -sh /www/server/pgsql/logs /var/log /opt/sub2api/data/logs
df -h /
```

首次治理结果作为容量基线：PostgreSQL 日志约 4.6 GiB、journal 约 703 MiB、根分区可用约 37 GB。若 3 天内日志仍快速增长，应先分析异常流量或重复错误，不应直接缩短保留期掩盖问题。

## 6. 网络与反向代理

8080 访问控制由 systemd 的 `ExecStartPre=+/usr/local/sbin/sub2api-8080-firewall` 刷新：

- 允许 loopback 和 WireGuard `10.111.0.0/24`。
- 允许反向代理公网 IP：`155.94.192.174`、`64.186.228.84`、`69.63.200.39`、`154.26.183.236`、`23.148.204.176`、`103.117.102.195`。
- 其它 IPv4 和全部 IPv6 来源访问 8080 时均丢弃。

五台非主站反向代理使用 nginx PROXY protocol 保留真实客户端 IP：stream 层将 aipick SNI 转到 `127.0.0.1:9443` 并启用 `proxy_protocol`；HTTPS vhost 使用 `listen ... proxy_protocol`、`real_ip_header proxy_protocol` 和 `set_real_ip_from 127.0.0.1`。

涉及的域名为 `www`、`uscu`、`usct`、`usall`、`jpct.aihub.pick.art`，配置位于各服务器的：

- `/www/server/panel/vhost/nginx/tcp/remnawave-sni.conf`
- `/www/server/panel/vhost/nginx/<domain>.conf`

修改反向代理后必须先运行 `nginx -t -q -c /www/server/nginx/conf/nginx.conf`，通过后再 reload，并验证公网 `/status` 与应用记录的 `client_ip` 均正确。

## 7. 周期性运维计划

| 周期 | 任务 | 完成标准 |
| --- | --- | --- |
| 每日（自动） | 执行三天日志保留 | timer 最近一次 `Result=success` |
| 每周 | 检查磁盘、日志目录、服务重启次数和慢查询 | 根分区低于 75%，无异常增长或重启 |
| 每次部署 | 备份数据库、发布、验证、更新当前基线 | 备份可列出，全部健康检查通过，可回滚 |
| 每次部署后 | 按策略检查旧 release 和数据库备份 | 当前/上一版及最近两个有效备份仍在 |
| 每月 | 抽查备份可恢复性、反代真实 IP、8080 白名单 | 恢复清单有效，真实 IP 正确，白名单无过期节点 |
| 磁盘达到 75% | 执行容量分析 | 找出增长目录并形成处理决定 |
| 磁盘达到 85% | 执行安全清理 | 清理后服务、数据库和公网入口正常 |

## 8. 最小变更记录规范

后续只保留最近一次有效基线，不再追加大段流水账。每次变更记录以下字段即可：

- 时间、部署分支、Git HEAD、上游基线和版本号。
- 当前 release、上一版 release、数据库备份路径及校验结果。
- 数据库迁移前后数量。
- 本机、公网、systemd 和关键日志验证结果。
- 异常、回滚条件及最终处理结果。

详细历史以 Git 提交、systemd journal、release 目录和备份清单为准。
