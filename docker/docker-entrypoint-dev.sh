#!/bin/bash
set -e
export LANG=C.UTF-8
export LC_ALL=C.UTF-8

export MISE_HIDE_UPDATE_WARNING=1

COLOR_PREFIX="\033[1;36m[Entrypoint]\033[0m"
log() {
    printf "${COLOR_PREFIX} %s\n" "$1"
}

log "Initializing Development Environment..."

DEV_UID=${DEV_UID:-1000}
DEV_GID=${DEV_GID:-1000}

if ! id -u devuser >/dev/null 2>&1; then
    log "Creating user 'devuser' (UID: $DEV_UID, GID: $DEV_GID)..."
    groupadd -o -g "$DEV_GID" devgroup
    useradd -o -u "$DEV_UID" -g "$DEV_GID" -m -s /bin/bash devuser

    echo "devuser ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/devuser
    chmod 0440 /etc/sudoers.d/devuser
fi

chown devuser:devgroup /app/envs /app/web/node_modules /go/pkg/mod 2>/dev/null || true

MISE_DIR="/app/envs/mise"
mkdir -p "$MISE_DIR"
log "Syncing mise environment from base..."
rsync -a --ignore-existing /opt/mise-dev/ "$MISE_DIR" || true
log "Mise environment synced"

chown -R ${DEV_UID}:${DEV_GID} /opt/mise-dev /go /root/.npm /app

export MISE_DATA_DIR="$MISE_DIR"
export MISE_CONFIG_DIR="$MISE_DIR"
export PATH="$MISE_DIR/shims:$MISE_DIR/bin:$PATH"

export PIP_INDEX_URL=https://pypi.tuna.tsinghua.edu.cn/simple

log "mise version: $(mise --version 2>/dev/null | head -n 1)"
log "python: $(python --version 2>&1 | head -n 1) at $(which python)"
log "node: $(node --version 2>&1 | head -n 1) at $(which node)"
log "npm: $(npm --version 2>&1 | head -n 1) at $(which npm)"

log "Environment ready! Starting make dev..."
exec gosu devuser "$@"
