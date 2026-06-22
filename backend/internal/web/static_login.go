//go:build embed

package web

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/gin-gonic/gin"
)

const staticLoginSettingsTimeout = 2 * time.Second

type staticLoginConfig struct {
	SiteName     string `json:"site_name"`
	SiteSubtitle string `json:"site_subtitle"`
}

var staticLoginTemplate = template.Must(template.New("static-login").Parse(staticLoginHTMLTemplate))

func isStaticLoginPath(path string) bool {
	trimmed := strings.TrimSpace(path)
	return trimmed == "/login" || trimmed == "/login/"
}

func defaultStaticLoginConfig() staticLoginConfig {
	return staticLoginConfig{
		SiteName:     "企业数据中台",
		SiteSubtitle: "统一数据目录、治理与服务编排入口",
	}
}

func staticLoginConfigFromSettings(settings any) staticLoginConfig {
	cfg := defaultStaticLoginConfig()
	if settings == nil {
		return cfg
	}

	raw, err := json.Marshal(settings)
	if err != nil {
		return cfg
	}

	var injected staticLoginConfig
	if err := json.Unmarshal(raw, &injected); err != nil {
		return cfg
	}

	if siteName := strings.TrimSpace(injected.SiteName); siteName != "" {
		cfg.SiteName = siteName
	}
	if subtitle := strings.TrimSpace(injected.SiteSubtitle); subtitle != "" {
		cfg.SiteSubtitle = subtitle
	}
	return cfg
}

func (s *FrontendServer) serveStaticLoginHTML(c *gin.Context) {
	cfg := defaultStaticLoginConfig()
	if s != nil && s.settings != nil {
		ctx, cancel := context.WithTimeout(c.Request.Context(), staticLoginSettingsTimeout)
		defer cancel()

		if settings, err := s.settings.GetPublicSettingsForInjection(ctx); err == nil {
			cfg = staticLoginConfigFromSettings(settings)
		}
	}
	serveStaticLoginHTML(c, cfg)
}

func serveStaticLoginHTML(c *gin.Context, cfg staticLoginConfig) {
	content, err := renderStaticLoginHTML(cfg, middleware.GetNonceFromContext(c))
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to render login page")
		c.Abort()
		return
	}

	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
	c.Abort()
}

func renderStaticLoginHTML(cfg staticLoginConfig, nonce string) ([]byte, error) {
	data := struct {
		Nonce        string
		SiteName     string
		SiteSubtitle string
	}{
		Nonce:        nonce,
		SiteName:     strings.TrimSpace(cfg.SiteName),
		SiteSubtitle: strings.TrimSpace(cfg.SiteSubtitle),
	}
	if data.SiteName == "" {
		data.SiteName = defaultStaticLoginConfig().SiteName
	}
	if data.SiteSubtitle == "" {
		data.SiteSubtitle = defaultStaticLoginConfig().SiteSubtitle
	}

	var buf bytes.Buffer
	if err := staticLoginTemplate.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

const staticLoginHTMLTemplate = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <link rel="icon" type="image/png" sizes="64x64" href="/static/app/favicon.ico?v=cellular-network-20260621">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{.SiteName}} - 安全登录</title>
  <style nonce="{{.Nonce}}">
    :root {
      color-scheme: light;
      --bg: #f4f7fb;
      --panel: #ffffff;
      --panel-soft: #eef5f4;
      --ink: #172033;
      --muted: #637083;
      --line: #d9e2ea;
      --accent: #0f766e;
      --accent-strong: #0b5f5a;
      --accent-soft: #dff4f1;
      --warn: #b42318;
      --focus: #2563eb;
      font-family: Inter, "Segoe UI", "PingFang SC", "Microsoft YaHei", Arial, sans-serif;
    }

    * {
      box-sizing: border-box;
    }

    body {
      min-height: 100vh;
      margin: 0;
      color: var(--ink);
      background:
        linear-gradient(120deg, rgba(15, 118, 110, 0.10), transparent 32%),
        linear-gradient(245deg, rgba(37, 99, 235, 0.10), transparent 35%),
        var(--bg);
    }

    body::before {
      position: fixed;
      inset: 0;
      z-index: -1;
      pointer-events: none;
      content: "";
      background-image:
        linear-gradient(rgba(23, 32, 51, 0.045) 1px, transparent 1px),
        linear-gradient(90deg, rgba(23, 32, 51, 0.045) 1px, transparent 1px);
      background-size: 36px 36px;
      mask-image: linear-gradient(to bottom, rgba(0, 0, 0, 0.75), transparent);
    }

    .shell {
      display: grid;
      grid-template-columns: minmax(0, 1fr) 432px;
      gap: 48px;
      align-items: center;
      width: min(1120px, calc(100% - 48px));
      min-height: 100vh;
      margin: 0 auto;
      padding: 40px 0;
    }

    .brand {
      display: grid;
      gap: 28px;
    }

    .brand-mark {
      display: inline-flex;
      align-items: center;
      gap: 12px;
      width: fit-content;
      padding: 8px 12px 8px 8px;
      border: 1px solid rgba(15, 118, 110, 0.20);
      border-radius: 8px;
      background: rgba(255, 255, 255, 0.74);
      box-shadow: 0 16px 40px rgba(23, 32, 51, 0.08);
      backdrop-filter: blur(16px);
    }

    .brand-icon {
      display: grid;
      width: 38px;
      height: 38px;
      place-items: center;
      border-radius: 8px;
      color: #ffffff;
      font-weight: 800;
      background: linear-gradient(135deg, var(--accent), #2563eb);
    }

    .brand-name {
      font-size: 15px;
      font-weight: 700;
    }

    h1 {
      max-width: 680px;
      margin: 0;
      font-size: clamp(38px, 5vw, 58px);
      line-height: 1.05;
      letter-spacing: 0;
    }

    .lead {
      max-width: 660px;
      margin: 0;
      color: var(--muted);
      font-size: 18px;
      line-height: 1.75;
    }

    .capabilities {
      display: grid;
      grid-template-columns: repeat(2, minmax(0, 1fr));
      gap: 14px;
      max-width: 680px;
    }

    .capability {
      min-height: 116px;
      padding: 18px;
      border: 1px solid rgba(121, 138, 158, 0.28);
      border-radius: 8px;
      background: rgba(255, 255, 255, 0.72);
      box-shadow: 0 18px 44px rgba(23, 32, 51, 0.08);
      backdrop-filter: blur(18px);
    }

    .capability span {
      display: inline-block;
      margin-bottom: 10px;
      color: var(--accent-strong);
      font-size: 13px;
      font-weight: 800;
    }

    .capability strong {
      display: block;
      margin-bottom: 6px;
      font-size: 16px;
    }

    .capability p {
      margin: 0;
      color: var(--muted);
      font-size: 13px;
      line-height: 1.6;
    }

    .login-panel {
      width: 100%;
      padding: 32px;
      border: 1px solid rgba(121, 138, 158, 0.30);
      border-radius: 8px;
      background: rgba(255, 255, 255, 0.92);
      box-shadow: 0 26px 70px rgba(23, 32, 51, 0.14);
      backdrop-filter: blur(20px);
    }

    .panel-head {
      display: flex;
      justify-content: space-between;
      gap: 16px;
      align-items: flex-start;
      margin-bottom: 28px;
    }

    .panel-head h2 {
      margin: 0 0 8px;
      font-size: 24px;
      line-height: 1.2;
      letter-spacing: 0;
    }

    .panel-head p {
      margin: 0;
      color: var(--muted);
      font-size: 14px;
      line-height: 1.6;
    }

    .status-pill {
      flex: 0 0 auto;
      padding: 6px 10px;
      border-radius: 999px;
      color: var(--accent-strong);
      background: var(--accent-soft);
      font-size: 12px;
      font-weight: 700;
      white-space: nowrap;
    }

    form {
      display: grid;
      gap: 18px;
    }

    label {
      display: grid;
      gap: 8px;
      color: #2f3a4d;
      font-size: 13px;
      font-weight: 700;
    }

    .field {
      display: flex;
      align-items: center;
      height: 46px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: #ffffff;
      transition: border-color 0.16s ease, box-shadow 0.16s ease;
    }

    .field:focus-within {
      border-color: var(--focus);
      box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.14);
    }

    input {
      width: 100%;
      min-width: 0;
      border: 0;
      outline: 0;
      color: var(--ink);
      background: transparent;
      font: inherit;
      font-size: 15px;
      padding: 0 14px;
    }

    input::placeholder {
      color: #9aa6b6;
    }

    .toggle {
      width: 46px;
      height: 44px;
      border: 0;
      border-left: 1px solid var(--line);
      border-radius: 0 8px 8px 0;
      color: var(--muted);
      background: transparent;
      cursor: pointer;
      font-size: 13px;
      font-weight: 700;
    }

    .toggle:hover {
      color: var(--accent-strong);
      background: #f7fafc;
    }

    .submit {
      display: inline-flex;
      justify-content: center;
      align-items: center;
      width: 100%;
      height: 46px;
      border: 0;
      border-radius: 8px;
      color: #ffffff;
      background: linear-gradient(135deg, var(--accent), #1d4ed8);
      box-shadow: 0 16px 30px rgba(15, 118, 110, 0.24);
      cursor: pointer;
      font-size: 15px;
      font-weight: 800;
    }

    .submit:hover {
      filter: brightness(0.98);
    }

    .submit:disabled {
      cursor: wait;
      opacity: 0.72;
    }

    .error {
      min-height: 20px;
      margin: -4px 0 0;
      color: var(--warn);
      font-size: 13px;
      line-height: 1.5;
    }

    .totp {
      display: none;
      margin-top: 18px;
      padding-top: 18px;
      border-top: 1px solid var(--line);
    }

    .totp.active {
      display: grid;
    }

    .totp-note {
      margin: -4px 0 0;
      color: var(--muted);
      font-size: 13px;
      line-height: 1.6;
    }

    .panel-foot {
      display: flex;
      justify-content: space-between;
      gap: 16px;
      align-items: center;
      margin-top: 24px;
      color: var(--muted);
      font-size: 13px;
      line-height: 1.6;
    }

    .panel-foot a {
      color: var(--accent-strong);
      font-weight: 700;
      text-decoration: none;
    }

    .panel-foot a:hover {
      text-decoration: underline;
    }

    @media (max-width: 900px) {
      .shell {
        grid-template-columns: 1fr;
        gap: 28px;
        width: min(680px, calc(100% - 32px));
        align-items: start;
        padding: 28px 0;
      }

      .brand {
        gap: 18px;
      }

      h1 {
        font-size: 36px;
      }

      .lead {
        font-size: 16px;
      }
    }

    @media (max-width: 560px) {
      .shell {
        width: min(100% - 24px, 680px);
        padding: 16px 0 24px;
      }

      .capabilities {
        grid-template-columns: 1fr;
      }

      .login-panel {
        padding: 24px 18px;
      }

      .panel-head,
      .panel-foot {
        display: grid;
      }

      .status-pill {
        width: fit-content;
      }
    }
  </style>
</head>
<body>
  <main class="shell">
    <section class="brand" aria-label="数据中台能力">
      <div class="brand-mark">
        <div class="brand-icon" aria-hidden="true">D</div>
        <div class="brand-name">{{.SiteName}}</div>
      </div>
      <h1>连接、治理并服务企业核心数据资产</h1>
      <p class="lead">{{.SiteSubtitle}}，面向跨系统资料整合、数据血缘、质量规则、指标共享与服务发布提供统一工作入口。</p>
      <div class="capabilities">
        <article class="capability">
          <span>资产目录</span>
          <strong>统一数据目录</strong>
          <p>汇集结构化与半结构化资料，建立可检索、可追踪的数据资产地图。</p>
        </article>
        <article class="capability">
          <span>数据治理</span>
          <strong>治理与质量规则</strong>
          <p>沉淀标准、分级、血缘和校验策略，让资料状态持续透明。</p>
        </article>
        <article class="capability">
          <span>服务发布</span>
          <strong>数据服务发布</strong>
          <p>将可信数据封装为 API、报表和业务指标，支撑多团队协作。</p>
        </article>
        <article class="capability">
          <span>流程编排</span>
          <strong>流程编排监控</strong>
          <p>统一调度同步、转换、稽核与告警任务，提升数据链路稳定性。</p>
        </article>
      </div>
    </section>

    <section class="login-panel" aria-label="登录">
      <div class="panel-head">
        <div>
          <h2>安全登录</h2>
          <p>使用企业账号进入数据资产协作空间。</p>
        </div>
        <div class="status-pill">安全通道</div>
      </div>

      <form id="login-form" autocomplete="on" novalidate>
        <label>
          工作邮箱
          <span class="field">
            <input id="email" name="email" type="email" inputmode="email" autocomplete="username" placeholder="name@company.com" required>
          </span>
        </label>
        <label>
          访问密码
          <span class="field">
            <input id="password" name="password" type="password" autocomplete="current-password" placeholder="输入访问密码" required>
            <button class="toggle" id="toggle-password" type="button" aria-label="显示或隐藏密码">显示</button>
          </span>
        </label>
        <p class="error" id="login-error" role="alert" aria-live="polite"></p>
        <button class="submit" id="login-submit" type="submit">登录数据中台</button>
      </form>

      <form class="totp" id="totp-form" autocomplete="one-time-code" novalidate>
        <p class="totp-note" id="totp-note">请输入动态验证码完成二次校验。</p>
        <label>
          动态验证码
          <span class="field">
            <input id="totp-code" name="totp_code" type="text" inputmode="numeric" autocomplete="one-time-code" placeholder="6 位数字" maxlength="6" pattern="[0-9]{6}">
          </span>
        </label>
        <p class="error" id="totp-error" role="alert" aria-live="polite"></p>
        <button class="submit" id="totp-submit" type="submit">完成校验</button>
      </form>

      <div class="panel-foot">
        <span>账号由管理员统一开通。</span>
        <a href="/home">返回门户首页</a>
      </div>
    </section>
  </main>

  <script nonce="{{.Nonce}}">
    (function () {
      var loginForm = document.getElementById("login-form");
      var totpForm = document.getElementById("totp-form");
      var emailInput = document.getElementById("email");
      var passwordInput = document.getElementById("password");
      var togglePassword = document.getElementById("toggle-password");
      var loginSubmit = document.getElementById("login-submit");
      var totpSubmit = document.getElementById("totp-submit");
      var loginError = document.getElementById("login-error");
      var totpError = document.getElementById("totp-error");
      var totpCode = document.getElementById("totp-code");
      var totpNote = document.getElementById("totp-note");
      var tempToken = "";

      function redirectTarget() {
        var params = new URLSearchParams(window.location.search);
        var redirect = params.get("redirect") || "/dashboard";
        if (redirect.charAt(0) === "/" && redirect.charAt(1) !== "/" && redirect.charAt(1) !== "\\") {
          return redirect;
        }
        return "/dashboard";
      }

      function setBusy(button, busy, text) {
        button.disabled = busy;
        button.textContent = busy ? "处理中..." : text;
      }

      function clearAuthStorage() {
        try {
          localStorage.removeItem("auth_token");
          localStorage.removeItem("refresh_token");
          localStorage.removeItem("auth_user");
          localStorage.removeItem("token_expires_at");
        } catch (error) {
          return;
        }
      }

      function storeAuth(payload) {
        if (!payload || !payload.access_token) {
          throw new Error("AUTH_PAYLOAD_INVALID");
        }

        clearAuthStorage();
        localStorage.setItem("auth_token", payload.access_token);
        if (payload.refresh_token) {
          localStorage.setItem("refresh_token", payload.refresh_token);
        }
        if (payload.expires_in) {
          localStorage.setItem("token_expires_at", String(Date.now() + payload.expires_in * 1000));
        }
        if (payload.user) {
          var user = Object.assign({}, payload.user);
          delete user.run_mode;
          localStorage.setItem("auth_user", JSON.stringify(user));
        }

        window.location.assign(redirectTarget());
      }

      function friendlyMessage(error) {
        if (!error || error.status === 0) {
          return "网络连接异常，请稍后重试。";
        }
        if (error.status === 400) {
          return "请填写有效的账号信息。";
        }
        if (error.status === 401 || error.status === 403) {
          return "认证失败，请检查账号、密码或访问权限。";
        }
        if (error.status === 429) {
          return "尝试次数过多，请稍后再试。";
        }
        return "服务暂时不可用，请稍后重试。";
      }

      async function postJSON(url, payload) {
        var response = await fetch(url, {
          method: "POST",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
            "Accept-Language": navigator.language || "zh-CN"
          },
          body: JSON.stringify(payload)
        });

        var body = null;
        try {
          body = await response.json();
        } catch (error) {
          body = null;
        }

        if (!response.ok || (body && typeof body.code !== "undefined" && body.code !== 0)) {
          var err = new Error("REQUEST_FAILED");
          err.status = response.status;
          throw err;
        }

        if (body && Object.prototype.hasOwnProperty.call(body, "data")) {
          return body.data;
        }
        return body;
      }

      togglePassword.addEventListener("click", function () {
        var visible = passwordInput.type === "text";
        passwordInput.type = visible ? "password" : "text";
        togglePassword.textContent = visible ? "显示" : "隐藏";
      });

      loginForm.addEventListener("submit", async function (event) {
        event.preventDefault();
        loginError.textContent = "";
        totpError.textContent = "";

        var email = emailInput.value.trim();
        var password = passwordInput.value;
        if (!email || !password) {
          loginError.textContent = "请输入工作邮箱和访问密码。";
          return;
        }

        setBusy(loginSubmit, true, "登录数据中台");
        try {
          var data = await postJSON("/user/login", {
            email: email,
            password: password
          });

          if (data && data.requires_2fa) {
            tempToken = data.temp_token || "";
            totpForm.classList.add("active");
            totpNote.textContent = data.user_email_masked
              ? "已向 " + data.user_email_masked + " 对应账号发起二次校验。"
              : "请输入动态验证码完成二次校验。";
            totpCode.focus();
            return;
          }

          storeAuth(data);
        } catch (error) {
          loginError.textContent = friendlyMessage(error);
        } finally {
          setBusy(loginSubmit, false, "登录数据中台");
        }
      });

      totpCode.addEventListener("input", function () {
        totpCode.value = totpCode.value.replace(/\D/g, "").slice(0, 6);
      });

      totpForm.addEventListener("submit", async function (event) {
        event.preventDefault();
        totpError.textContent = "";

        var code = totpCode.value.trim();
        if (!tempToken || code.length !== 6) {
          totpError.textContent = "请输入 6 位动态验证码。";
          return;
        }

        setBusy(totpSubmit, true, "完成校验");
        try {
          var data = await postJSON("/user/login/2fa", {
            temp_token: tempToken,
            totp_code: code
          });
          storeAuth(data);
        } catch (error) {
          totpError.textContent = friendlyMessage(error);
        } finally {
          setBusy(totpSubmit, false, "完成校验");
        }
      });
    })();
  </script>
</body>
</html>
`
