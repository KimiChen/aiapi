aihub 是 Docker Compose 形态：Caddy 反代本机 127.0.0.1:8080，应用跑在容器里，Postgres/Redis 也在容器里。
aihub 部署目录: /opt/compose/sub2api-deploy
aihub Compose 服务名: sub2api
aihub 容器名: sub2api, sub2api-postgres, sub2api-redis
aihub 线上网址: 网址
aihub 当前镜像: weishaw/sub2api:latest
aihub 查看容器状态: cd /opt/compose/sub2api-deploy && docker compose ps
aihub 查看应用日志: docker logs -f --tail 200 sub2api
aihub 健康检查: curl -fsS 网址/health
aihub Caddy access log:
 - 2026-06-22 已开启完整 JSON access log
 - Caddy 配置文件: /etc/caddy/Caddyfile
 - 日志文件: /var/log/caddy/aihub_access.log
 - 轮转策略: roll_size 100MiB, roll_keep 10, roll_keep_for 720h
 - 配置备份: /etc/caddy/Caddyfile.bak-accesslog-20260622-044859
 - 查看最新日志: tail -f /var/log/caddy/aihub_access.log
 - 统计访问路径: python3 -c 'import json,collections; c=collections.Counter(); [c.update([(json.loads(l).get("request") or {}).get("uri","").split("?")[0]]) for l in open("/var/log/caddy/aihub_access.log") if l.strip()]; print(c.most_common(50))'
 - 验证结果: caddy validate 通过，systemctl reload caddy 成功，网址/health 返回 200 且 access log 记录 /health
aihub 部署方式:
 - 本地先执行 `pnpm --dir frontend build`
 - 在 `backend/` 目录执行 linux/amd64 嵌入前端资源构建：
   - `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags embed -ldflags="-s -w -X main.Version=$(tr -d '\r\n' < cmd/server/VERSION)$(tr -d '\r\n' < cmd/server/SUB_VERSION) -X main.BuildType=release" -trimpath -o bin/server ./cmd/server`
 - 上传 `backend/bin/server` 到 aihub 的 `/opt/compose/sub2api-deploy/releases/<timestamp>/sub2api`
 - 基于旧镜像备份并替换 `/app/sub2api`；如旧镜像缺 `/app/resources`，同步 `backend/resources` 到镜像内 `/app/resources`
 - 只重建应用容器：`cd /opt/compose/sub2api-deploy && docker compose up -d --no-deps --force-recreate sub2api`
 - 部署期间不要修改 Caddy、Docker Compose、Postgres、Redis，除非用户明确要求
aihub 部署记录:
 - 部署时间: 2026-06-22
 - Git HEAD: 24a0dfbb
 - 版本号: 0.1.137-kim
 - 发布目录: /opt/compose/sub2api-deploy/releases/20260622-025825
 - 当前容器二进制 SHA256: 03394869e5200338afe508313eb095222c6ad1fa24d214fb15136c7053de2545
 - 当前镜像 ID: sha256:adf33a3f89e7086b10c956d7b6ff0208e98c570683cbf5710b497ed24649c786
 - 备份镜像: weishaw/sub2api:backup-20260622-025825
 - 验证结果: 网址/health 返回 {"status":"ok"}，首页返回 200，sub2api/sub2api-postgres/sub2api-redis 均 healthy
 - 部署时间: 2026-06-21
 - 发布目录: /opt/compose/sub2api-deploy/releases/20260621-211925
 - 当前容器二进制 SHA256: 639b81fdc84b5461585e423718c96c2f66f4382e7e688a3771ecbdcfdf7786c5
 - 当前镜像 ID: sha256:d2ae8a82d9a3ce6b5d9b806f0ca8f54cbe9f6c752a8e0e7bd573964ac412856e
 - 备份镜像: weishaw/sub2api:backup-20260621-211925
 - 验证结果: 网址/health 返回 {"status":"ok"}，sub2api/sub2api-postgres/sub2api-redis 均 healthy
aihub 回滚方式:
 - `cd /opt/compose/sub2api-deploy`
 - `docker tag weishaw/sub2api:backup-20260622-025825 weishaw/sub2api:latest`
 - `docker compose up -d --no-deps --force-recreate sub2api`
aihub 数据结构同步记录:
 - 2026-06-22 已确认本地与 aihub 的 `schema_migrations` 均为 190 条，缺失/额外/checksum mismatch 均为 0
 - 最新迁移: `155_add_usage_log_upstream_traffic_bytes.sql`
 - 155 迁移 checksum: 65ae40e173eb1ce15142666c3e1672d9c1341ecea62425c998f27dee7a819941
 - 2026-06-21 已同步 `154_add_usage_log_traffic_bytes.sql`
 - 本地与 aihub 的 `schema_migrations` 均为 189 条，缺失/额外/ checksum mismatch 均为 0
 - 迁移 checksum: 4c5330329bf53b585f45797015ecbdc4ea5cabc703dcb919cecbc60a72300ccf
