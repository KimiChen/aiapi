# Racknerd Change Log

## 2026-06-29

- Disabled and stopped `caddy.service` on the racknerd server because it conflicted with the BaoTa-managed Nginx listener on port 80.
- Ran `systemctl reset-failed caddy` after disabling the unit.
- Verified `caddy.service` is `disabled` and `inactive`.
- Verified `systemctl --failed` reports zero failed units.
- Verified BaoTa Nginx config with `nginx -t`.
- Verified the public HTTPS site still returns `HTTP/2 200` from Nginx.
