#!/usr/bin/env bash
# fork 与上游同步脚本
#
# 用法: scripts/sync_upstream.sh [--skip-build]
#
# 流程:
#   1. 配置 merge.ours 合并驱动(fork 接管文件,见 .gitattributes)
#   2. fetch upstream 并合并 upstream/main
#   3. 提示 fork 接管文件在上游的更新(merge=ours 会保留 fork 版本,需人工评审移植)
#   4. ent schema 有变化时重新生成代码
#   5. 后端/前端构建验证
#
# 冲突处理约定:
#   - ent 生成物(ent/*.go,schema 除外)冲突: 任选一方接受后重新生成
#       git checkout --theirs backend/ent && git add backend/ent
#   - backend/cmd/server/VERSION: 始终取上游
#       git checkout --theirs backend/cmd/server/VERSION && git add backend/cmd/server/VERSION
#   - 其余冲突: 手工解决,fork 定制应优先保留(原则上 fork 定制都在 fork 自有文件中)
set -euo pipefail
cd "$(dirname "$0")/../.."

SKIP_BUILD=false
for arg in "$@"; do
  case "$arg" in
    --skip-build) SKIP_BUILD=true ;;
    *) echo "未知参数: $arg"; exit 2 ;;
  esac
done

# fork 接管文件列表(与 .gitattributes 中 merge=ours 条目保持一致)
FORK_OWNED_FILES=(
  backend/internal/web/embed_on.go
  backend/internal/web/embed_test.go
  backend/internal/web/static_cache.go
  backend/internal/web/static_cache_test.go
)

echo "==> 配置 merge.ours 合并驱动"
git config merge.ours.driver true

echo "==> fetch upstream"
git fetch upstream

BASE="$(git merge-base HEAD upstream/main)"
echo "==> 合并基点: $BASE"

echo "==> 检查 fork 接管文件在上游的更新"
NEED_REVIEW=false
for f in "${FORK_OWNED_FILES[@]}"; do
  if ! git diff --quiet "$BASE" upstream/main -- "$f" 2>/dev/null; then
    echo "    [需评审] $f 上游有更新;合并将保留 fork 版本,请人工对比移植"
    NEED_REVIEW=true
  fi
done
$NEED_REVIEW || echo "    (无)"

echo "==> 合并 upstream/main"
if ! git merge --no-edit upstream/main; then
  cat << 'MSG'
!! 合并出现冲突,按约定处理:
   - ent 生成物: git checkout --theirs backend/ent && git add backend/ent
   - VERSION:    git checkout --theirs backend/cmd/server/VERSION && git add backend/cmd/server/VERSION
   - 其余文件:   手工解决(优先保留 fork 定制)
   解决后运行: git commit --no-edit
   然后重新运行本脚本完成生成与验证。
MSG
  exit 1
fi

if ! git diff --quiet "$BASE" HEAD -- backend/ent/schema/ 2>/dev/null; then
  echo "==> ent schema 有变化,重新生成 ent 代码"
  (cd backend && go generate ./ent)
  git add backend/ent
  git diff --cached --quiet || git commit -m "chore: regenerate ent code after upstream sync"
fi

if [ "$SKIP_BUILD" = false ]; then
  echo "==> 后端构建验证"
  (cd backend && go build ./...)

  echo "==> 前端构建验证"
  pnpm --dir frontend build
fi

echo "==> 同步完成"
