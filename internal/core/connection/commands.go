package connection

// Constantes de comandos padronizados e seguros para execução em VPS / Local
const (
	// Inspeção básica do sistema
	CmdUnameKernel = "uname -srm"
	CmdCheckRoot   = "id -u"

	// 1. Detecção do binário do Docker (systemd service unit + PATH + /proc)
	CmdDiscoverDockerBin = `
		if command -v systemctl >/dev/null 2>&1; then
			p=$(systemctl show -p ExecStart docker.service 2>/dev/null | grep -o 'path=[^; ]*' | cut -d= -f2)
			[ -n "$p" ] && [ -x "$p" ] && echo "$p" && exit 0
		fi
		p=$(command -v docker 2>/dev/null || which docker 2>/dev/null)
		[ -n "$p" ] && [ -x "$p" ] && echo "$p" && exit 0
		for c in /usr/bin/docker /usr/local/bin/docker /snap/bin/docker /usr/libexec/docker/cli-plugins/docker; do
			[ -x "$c" ] && echo "$c" && exit 0
		done
		echo "/usr/bin/docker"
	`

	// 2. Detecção do Docker Compose (Caminho absoluto do arquivo binário standalone ou do plugin CLI V2)
	CmdDiscoverCompose = `
		# 1. Procura binário standalone docker-compose no PATH
		p=$(command -v docker-compose 2>/dev/null || which docker-compose 2>/dev/null)
		[ -n "$p" ] && [ -x "$p" ] && echo "$p" && exit 0

		# 2. Procura caminhos físicos conhecidos de standalone e plugins CLI V2
		for c in \
			/usr/local/bin/docker-compose \
			/usr/bin/docker-compose \
			/usr/libexec/docker/cli-plugins/docker-compose \
			/usr/local/lib/docker/cli-plugins/docker-compose \
			/usr/lib/docker/cli-plugins/docker-compose \
			"$HOME/.docker/cli-plugins/docker-compose" \
			/snap/bin/docker-compose; do
			[ -n "$c" ] && [ -x "$c" ] && echo "$c" && exit 0
		done

		# 3. Fallback: se 'docker compose' V2 estiver ativo, retorna o caminho do próprio docker
		if docker compose version >/dev/null 2>&1; then
			p=$(command -v docker 2>/dev/null || echo "/usr/bin/docker")
			echo "$p compose"
			exit 0
		fi

		echo "/usr/bin/docker-compose"
	`

	// 3. Detecção Infalível do Socket UNIX do Docker (Root + Rootless + Snap + Docker Context)
	CmdDiscoverSocket = `
		# 1. Checa Docker Context ativo do usuário (pega o endpoint rootless oficial se configurado)
		if command -v docker >/dev/null 2>&1; then
			ctx=$(docker context inspect --format '{{.Endpoints.docker.Host}}' 2>/dev/null)
			if [[ "$ctx" == unix://* ]]; then
				s="${ctx#unix://}"
				[ -S "$s" ] && echo "$s" && exit 0
			fi
		fi

		# 2. Checa Systemd do Usuário (Rootless Docker: dockerd-rootless.sh)
		if command -v systemctl >/dev/null 2>&1; then
			s=$(systemctl --user show -p ListenStream docker.socket 2>/dev/null | grep -o '/.*')
			[ -n "$s" ] && [ -S "$s" ] && echo "$s" && exit 0
		fi

		# 3. Checa Systemd Global do Sistema (Docker tradicional como root)
		if command -v systemctl >/dev/null 2>&1; then
			s=$(systemctl show -p ListenStream docker.socket 2>/dev/null | grep -o '/.*')
			[ -n "$s" ] && [ -S "$s" ] && echo "$s" && exit 0
		fi

		# 4. Checa Variável DOCKER_HOST
		if [ -n "$DOCKER_HOST" ] && [[ "$DOCKER_HOST" == unix://* ]]; then
			s="${DOCKER_HOST#unix://}"
			[ -S "$s" ] && echo "$s" && exit 0
		fi

		# 5. Varredura nos caminhos conhecidos (Rootless / User runtime + Root + Snap)
		uid=$(id -u 2>/dev/null || echo 1000)
		for s in \
			"/run/user/$uid/docker.sock" \
			"$XDG_RUNTIME_DIR/docker.sock" \
			"$HOME/.docker/run/docker.sock" \
			"$HOME/.docker/desktop/docker.sock" \
			"/var/run/docker.sock" \
			"/run/docker.sock" \
			"/var/snap/docker/current/run/docker.sock"; do
			[ -n "$s" ] && [ -S "$s" ] && echo "$s" && exit 0
		done

		echo "/var/run/docker.sock"
	`

	// 4. Verificação de status do daemon
	CmdCheckDockerVersion = "docker info --format '{{.ServerVersion}}' 2>/dev/null || echo ''"
)
