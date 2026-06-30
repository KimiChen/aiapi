aihub 运维记录:

- 时间: 2026-07-01 02:05 Asia/Shanghai
- 类型: 反向代理 TLS/SNI 配置修复
- 问题: 某公网代理入口的 HTTPS 证书校验间歇性失败，同一 SNI 偶发返回其它站点证书。
- 原因: nginx 同时存在 stream SNI 分流入口和 HTTP 站点公网 TLS 监听，公网 443 连接会被随机分配到不同监听路径；其中一条路径转到本机 TLS 默认站点后返回了不匹配证书。
- 处理: 将相关 HTTP 站点的公网 TLS 监听改为本机 TLS 后端监听，公网 443 统一交给 stream SNI 分流入口；同时将 HTTPS 跳转判断改为基于 scheme，避免经本机后端端口访问时循环跳转；禁用该站点的公网 QUIC/H3 广告。
- 验证: nginx 配置测试通过并 reload；连续多次公网 TLS 主机名校验均 OK；公网首页返回 HTTP/2 200。

