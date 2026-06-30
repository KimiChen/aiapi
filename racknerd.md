# Racknerd Change Log

## 2026-06-30

- Refreshed racknerd PostgreSQL from an aihub snapshot taken after stopping the aihub `sub2api` application container.
  - Source dump: `/opt/compose/sub2api-deploy/backups/sub2api-stopped-20260630-004051.dump` on aihub.
  - Target dump: `/opt/sub2api/backups/sub2api-stopped-20260630-004051.dump` on racknerd.
  - Dump SHA256: `16e47e05a0b543ab25cf080f0309bfafc3900f3feae8af0923bcb9a2ea31e7b6`.
  - racknerd pulled the dump directly from aihub with the existing dedicated SSH key; no local-machine transfer was used.
  - Source observation after stopping aihub app: `schema_migrations=190`, `public_tables=74`, `accounts=1037`, `settings=225`, `usage_logs=663743`, `usage_billing_dedup=663802`, `ops_system_logs=3066807`, `source_db_size=4615 MB`.
- Restored the stopped-source dump into racknerd local BaoTa PostgreSQL.
  - Stopped `sub2api.service` during restore, dropped/recreated database `sub2api`, and restored with `/www/server/pgsql/bin/pg_restore`.
  - Fixed restored object ownership and privileges for role `sub2api`.
  - `pg_restore` log: `/opt/sub2api/backups/pg_restore-20260630-004051.log`; the log was empty after successful restore.
  - Verification before restarting racknerd app matched source counts exactly: `schema_migrations=190`, `public_tables=74`, `accounts=1037`, `settings=225`, `usage_logs=663743`, `usage_billing_dedup=663802`, `ops_system_logs=3066807`, `target_db_size=3522 MB`, `non_sub2api_owned_tables=0`.
  - Post-start verification: `schema_migrations=190`, `public_tables=74`, `accounts=1037`, `settings=226`, `usage_logs=663743`, `usage_billing_dedup=663802`, `ops_system_logs=3066816`, `target_db_size=3522 MB`, `non_sub2api_owned_tables=0`.
  - aihub `sub2api` application container remained stopped after the migration; `sub2api-postgres` and `sub2api-redis` remained running.
- Verification after stopped-source refresh:
  - `systemctl is-active sub2api` => `active`.
  - `systemctl is-enabled sub2api` => `enabled`.
  - `systemctl --failed` reports zero failed units.
  - `http://127.0.0.1:8080/status` returns the expected status JSON.
  - `https://wu.ci/status` returns the expected status JSON.
  - `https://wu.ci/login` returns 200.

- Refreshed racknerd PostgreSQL from a new aihub online snapshot.
  - Source dump: `/opt/compose/sub2api-deploy/backups/sub2api-online-20260629-163442.dump` on aihub.
  - Target dump: `/opt/sub2api/backups/sub2api-online-20260629-163442.dump` on racknerd.
  - Dump SHA256: `0d90eec791a31abd32894c8495fe16756381a23386e411d35967645c33357d3e`.
  - racknerd pulled the dump directly from aihub with the existing dedicated SSH key; no local-machine transfer was used.
  - Source observation before dump: `schema_migrations=190`, `public_tables=74`, `accounts=1037`, `settings=225`, `usage_logs=662608`, `usage_billing_dedup=662667`, `ops_system_logs=3050859`, `source_db_size=4605 MB`.
- Restored the new dump into racknerd local BaoTa PostgreSQL.
  - Stopped `sub2api.service` during restore, dropped/recreated database `sub2api`, and restored with `/www/server/pgsql/bin/pg_restore`.
  - Fixed restored object ownership and privileges for role `sub2api`.
  - `pg_restore` log: `/opt/sub2api/backups/pg_restore-20260629-163442.log`; the log was empty after successful restore.
  - Post-start verification: `schema_migrations=190`, `public_tables=74`, `accounts=1037`, `settings=225`, `usage_logs=662610`, `usage_billing_dedup=662669`, `ops_system_logs=3050874`, `target_db_size=3508 MB`, `non_sub2api_owned_tables=0`.
  - The small post-start count increase came from racknerd application activity after service restart.
- Verification after refresh:
  - `systemctl is-active sub2api` => `active`.
  - `systemctl is-enabled sub2api` => `enabled`.
  - `systemctl --failed` reports zero failed units.
  - `http://127.0.0.1:8080/status` returns `{"status":"perfectly nice"}`.
  - `https://wu.ci/status` returns `{"status":"perfectly nice"}`.
  - `https://wu.ci/login` returns 200.

## 2026-06-29

- Migrated aihub PostgreSQL data to racknerd with an online snapshot.
  - Source dump: `/opt/compose/sub2api-deploy/backups/sub2api-online-20260629-003850.dump` on aihub.
  - Target dump: `/opt/sub2api/backups/sub2api-online-20260629-003850.dump` on racknerd.
  - Dump SHA256: `ce761c48b1ea989e8708222bf9da0d841f8df93539cb8d546caa996482cc8ac6`.
  - racknerd now has a dedicated SSH key at `/root/.ssh/aihub_sub2api_pull_ed25519` for pulling backups directly from aihub.
  - Verified direct racknerd -> aihub rsync pull and matching SHA256.
- Restored the dump into racknerd local BaoTa PostgreSQL.
  - Created local PostgreSQL role/database `sub2api`.
  - Restored with `/www/server/pgsql/bin/pg_restore`.
  - Fixed restored object ownership so all public tables are owned by `sub2api`.
  - Verification: `schema_migrations=190`, `public_tables=74`, `target_db_size=3419 MB`, `non_sub2api_owned_tables=0`.
- Deployed this project on racknerd without Docker.
  - Frontend build command: `pnpm --dir frontend build`.
  - Backend binary: linux/amd64, embedded frontend, version `0.1.139.kim`.
  - Binary SHA256: `3c39c242a3922bf58a7bdbe8b72b64104bca7a814347bd6706bc8b6a6e48f8ea`.
  - Release path: `/opt/sub2api/releases/20260628-173811/sub2api`.
  - Current symlink: `/opt/sub2api/sub2api`.
  - Config path: `/etc/sub2api/config.yaml`.
  - Runtime data path: `/opt/sub2api/data`.
  - Synced runtime `public-route-blocklist.yaml` from aihub; startup log shows `source=file`, `effective_rules=20`.
  - Created and enabled `sub2api.service`; service listens on `127.0.0.1:8080`.
- Switched BaoTa Nginx for `wu.ci` to reverse proxy Sub2API.
  - Nginx backup: `/www/server/panel/vhost/nginx/wu.ci.conf.bak-sub2api-20260629-0100`.
  - Current site proxies to `http://127.0.0.1:8080`.
  - Added `underscores_in_headers on` and standard forwarded headers.
  - Kept existing certificate, HTTP->HTTPS redirect, `.well-known` include, HTTP/2 and HTTP/3 settings.
- Verification after migration/deploy:
  - `systemctl is-active sub2api` => `active`.
  - `systemctl is-enabled sub2api` => `enabled`.
  - `systemctl --failed` reports zero failed units.
  - `nginx -t` passes.
  - `https://wu.ci/status` returns `{"status":"perfectly nice"}`.
  - `https://wu.ci/login` returns 200.
  - `https://wu.ci/static/app/logo.png` returns 200.

- Disabled and stopped `caddy.service` on the racknerd server because it conflicted with the BaoTa-managed Nginx listener on port 80.
- Ran `systemctl reset-failed caddy` after disabling the unit.
- Verified `caddy.service` is `disabled` and `inactive`.
- Verified `systemctl --failed` reports zero failed units.
- Verified BaoTa Nginx config with `nginx -t`.
- Verified the public HTTPS site still returns `HTTP/2 200` from Nginx.
