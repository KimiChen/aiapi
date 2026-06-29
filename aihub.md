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
 - 轮转策略: roll_size 1GiB, roll_keep 9；当前文件加 9 个历史文件总占用约 10GiB
 - 配置备份: /etc/caddy/Caddyfile.bak-accesslog-20260622-044859, /etc/caddy/Caddyfile.bak-logrotate-10g-20260622-045436
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
 - 部署时间: 2026-06-23
 - Git HEAD: 186f333e
 - 版本号: 0.1.138.kim
 - 发布目录: /opt/compose/sub2api-deploy/releases/20260623-032802
 - 当前容器二进制 SHA256: 06cea45fd0c3109f80acda5f411d49497510649d49b8f3d40c6f02736a7085e6
 - 当前镜像 ID: sha256:a5ae7a4988154184ade27ccddd350e01ba214094ebf45789664bd7067c5eb2bb
 - 备份镜像: weishaw/sub2api:backup-20260623-032802
 - 备份镜像 ID: sha256:adf33a3f89e7086b10c956d7b6ff0208e98c570683cbf5710b497ed24649c786
 - Compose 健康检查保持 /health，部署期间未修改 Caddy、Docker Compose、Postgres、Redis
 - 验证结果: 网址/health 返回 {"status":"ok"}，/login、/register、/email-verify 返回 200，/api/v1/auth/login 返回 JSON 校验错误，sub2api/sub2api-postgres/sub2api-redis 均 healthy
 - 回滚时间: 2026-06-23
 - 回滚到备份镜像: weishaw/sub2api:backup-20260623-024554
 - 回滚后版本号: 0.1.137-kim
 - 回滚后容器二进制 SHA256: 03394869e5200338afe508313eb095222c6ad1fa24d214fb15136c7053de2545
 - 回滚后镜像 ID: sha256:adf33a3f89e7086b10c956d7b6ff0208e98c570683cbf5710b497ed24649c786
 - Compose 健康检查已恢复为 /health，回滚前 Compose 备份文件: /opt/compose/sub2api-deploy/docker-compose.yml.bak-before-rollback-20260622-191847
 - 验证结果: 网址/health 返回 {"status":"ok"}，/login、/register、/email-verify 返回 200，/api/v1/auth/login 返回 JSON 校验错误，sub2api/sub2api-postgres/sub2api-redis 均 healthy
 - 部署时间: 2026-06-23
 - Git HEAD: b43d271e
 - 版本号: 0.1.138.kim
 - 发布目录: /opt/compose/sub2api-deploy/releases/20260623-024554
 - 当前容器二进制 SHA256: ce38cc2e85a1f626d326ce0d8be9c3061a03d674dfc8c5fbb8a8594bb00a9271
 - 当前镜像 ID: sha256:5ddde443feaa9864de7a5f981af9a8046868e7bd4dddc33817abdabb8966ea29
 - 备份镜像: weishaw/sub2api:backup-20260623-024554
 - 运行时 blocklist: /opt/compose/sub2api-deploy/data/public-route-blocklist.yaml，启动日志显示 source=file、enabled=true、effective_rules=20
 - Compose 健康检查已从 /health 改为 /status，备份文件: /opt/compose/sub2api-deploy/docker-compose.yml.bak-status-20260623-024554
 - 验证结果: 网址/status 返回 {"status":"perfectly nice"}，/login 返回 200，/user/login 返回 JSON 校验错误而非 SPA HTML，旧 /api/v1/auth/login、/register、/setup/status 均返回 404，/v1/usage 未带 key 返回 401，/static/app/logo.png 返回 200，sub2api/sub2api-postgres/sub2api-redis 均 healthy
 - 状态: 已于 2026-06-23 回滚，不再是当前线上版本
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
 - `docker tag weishaw/sub2api:backup-20260623-032802 weishaw/sub2api:latest`
 - 确认健康检查使用当前版本支持的路径；本次部署和备份版本均支持 `/health`
 - `docker compose up -d --no-deps --force-recreate sub2api`
aihub 数据结构同步记录:
 - 2026-06-30 已从 aihub PostgreSQL 容器重新在线导出逻辑备份给 racknerd 刷新使用；未停止 aihub 应用，因此 dump 后的新写入不会包含在本次 racknerd 快照中
 - 源备份文件: /opt/compose/sub2api-deploy/backups/sub2api-online-20260629-163442.dump
 - 备份 SHA256: 0d90eec791a31abd32894c8495fe16756381a23386e411d35967645c33357d3e
 - 导出源: sub2api-postgres 容器内 PostgreSQL，数据库 sub2api
 - 导出前观测: 数据库约 4605 MB，schema_migrations=190，public_tables=74，accounts=1037，settings=225，usage_logs=662608，usage_billing_dedup=662667，ops_system_logs=3050859
 - racknerd 使用已有专用 SSH key 从 aihub 直接拉取该备份，未从本机中转
 - 2026-06-29 已从 aihub PostgreSQL 容器在线导出逻辑备份给 racknerd 迁移使用；未停止 aihub 应用，因此 dump 后的新写入不会包含在本次 racknerd 快照中
 - 源备份文件: /opt/compose/sub2api-deploy/backups/sub2api-online-20260629-003850.dump
 - 备份 SHA256: ce761c48b1ea989e8708222bf9da0d841f8df93539cb8d546caa996482cc8ac6
 - 导出源: sub2api-postgres 容器内 PostgreSQL 18.4，数据库 sub2api
 - 导出前观测: 数据库约 4508 MB，schema_migrations=190，public_tables=74
 - 已给 racknerd 增加专用 SSH 公钥用于从 aihub 直接拉取该备份；authorized_keys 使用来源限制，不记录私钥
 - 同步给 racknerd 的运行数据: /opt/compose/sub2api-deploy/data/public-route-blocklist.yaml 和 data/pages/
 - 2026-06-22 已确认本地与 aihub 的 `schema_migrations` 均为 190 条，缺失/额外/checksum mismatch 均为 0
 - 最新迁移: `155_add_usage_log_upstream_traffic_bytes.sql`
 - 155 迁移 checksum: 65ae40e173eb1ce15142666c3e1672d9c1341ecea62425c998f27dee7a819941
 - 2026-06-21 已同步 `154_add_usage_log_traffic_bytes.sql`
 - 本地与 aihub 的 `schema_migrations` 均为 189 条，缺失/额外/ checksum mismatch 均为 0
 - 迁移 checksum: 4c5330329bf53b585f45797015ecbdc4ea5cabc703dcb919cecbc60a72300ccf
