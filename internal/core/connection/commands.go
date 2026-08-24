package connection

// Constantes de comandos padronizados e seguros para execução em VPS / Local
const (
	// Inspeção básica do sistema
	CmdUnameKernel = "uname -srm"
	CmdCheckRoot   = "id -u"

	// Detecção de múltiplos binários do Docker
	CmdDiscoverAllDockerBins = `
		for c in \
			"$HOME/.local/bin/docker" \
			"$HOME/bin/docker" \
			/usr/bin/docker \
			/usr/local/bin/docker \
			/snap/bin/docker \
			/bin/docker; do
			if [ -n "$c" ] && [ -x "$c" ] && [ ! -d "$c" ]; then
				readlink -f "$c" 2>/dev/null || echo "$c"
			fi
		done
		p=$(PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:$HOME/.local/bin:$HOME/bin" command -v docker 2>/dev/null)
		if [ -n "$p" ] && [ -x "$p" ] && [ ! -d "$p" ]; then
			readlink -f "$p" 2>/dev/null || echo "$p"
		fi
	`

	// Detecção de múltiplos sockets do Docker
	CmdDiscoverAllSockets = `
		uid=$(id -u 2>/dev/null || echo 1000)
		for s in \
			"/run/user/$uid/docker.sock" \
			"$XDG_RUNTIME_DIR/docker.sock" \
			"$HOME/.docker/run/docker.sock" \
			"$HOME/.docker/desktop/docker.sock" \
			"/var/run/docker.sock" \
			"/run/docker.sock" \
			/var/snap/docker/current/run/docker.sock; do
			if [ -n "$s" ] && [ -S "$s" ]; then
				readlink -f "$s" 2>/dev/null || echo "$s"
			fi
		done
		if command -v systemctl >/dev/null 2>&1; then
			s=$(systemctl --user show -p ListenStream docker.socket 2>/dev/null | grep -o '/.*')
			[ -n "$s" ] && [ -S "$s" ] && { readlink -f "$s" 2>/dev/null || echo "$s"; }

			s=$(systemctl show -p ListenStream docker.socket 2>/dev/null | grep -o '/.*')
			[ -n "$s" ] && [ -S "$s" ] && { readlink -f "$s" 2>/dev/null || echo "$s"; }
		fi
	`

	// Detecção de múltiplos Docker Composes
	CmdDiscoverAllComposes = `
		for c in \
			"$HOME/.local/bin/docker-compose" \
			"$HOME/bin/docker-compose" \
			"$HOME/.docker/cli-plugins/docker-compose" \
			/usr/local/lib/docker/cli-plugins/docker-compose \
			/usr/lib/docker/cli-plugins/docker-compose \
			/usr/libexec/docker/cli-plugins/docker-compose \
			/usr/local/bin/docker-compose \
			/usr/bin/docker-compose \
			/snap/bin/docker-compose; do
			if [ -n "$c" ] && [ -x "$c" ] && [ ! -d "$c" ]; then
				readlink -f "$c" 2>/dev/null || echo "$c"
			fi
		done
		p=$(PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:$HOME/.local/bin:$HOME/bin" command -v docker-compose 2>/dev/null)
		if [ -n "$p" ] && [ -x "$p" ] && [ ! -d "$p" ]; then
			readlink -f "$p" 2>/dev/null || echo "$p"
		fi
		docker_bin=$(PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:$HOME/.local/bin:$HOME/bin" command -v docker 2>/dev/null || echo "/usr/bin/docker")
		if "$docker_bin" compose version >/dev/null 2>&1; then
			full_docker=$(readlink -f "$docker_bin" 2>/dev/null || echo "$docker_bin")
			echo "$full_docker compose"
		fi
	`

	// 4. Verificação de status do daemon
	CmdCheckDockerVersion = "docker info --format '{{.ServerVersion}}' 2>/dev/null || echo ''"
)
