package connection

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
)

// ExecCommand executa um comando no servidor (com ou sem elevação sudo inteligente)
func (c *Client) ExecCommand(cmd string, useSudo bool) (string, error) {
	finalCmd := c.buildFinalCmd(cmd, useSudo)

	if isLocal(c.server) {
		return execLocalCommand(finalCmd)
	}

	if c.sshClient == nil {
		return "", fmt.Errorf("cliente SSH não conectado")
	}

	session, err := c.sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("falha ao criar sessão SSH: %w", err)
	}
	defer session.Close()

	out, err := session.CombinedOutput(finalCmd)
	if err != nil {
		return string(out), fmt.Errorf("comando falhou (%w): %s", err, string(out))
	}

	return string(out), nil
}

// StartCommandOutput inicia um comando e retorna seu stdout como io.Reader para streaming contínuo
func (c *Client) StartCommandOutput(ctx context.Context, cmd string, useSudo bool) (io.Reader, func() error, error) {
	finalCmd := c.buildFinalCmd(cmd, useSudo)

	if isLocal(c.server) {
		cmdObj := exec.CommandContext(ctx, "bash", "-c", finalCmd)
		stdout, err := cmdObj.StdoutPipe()
		if err != nil {
			return nil, nil, fmt.Errorf("falha ao abrir stdout pipe local: %w", err)
		}
		stderr, err := cmdObj.StderrPipe()
		if err != nil {
			return nil, nil, fmt.Errorf("falha ao abrir stderr pipe local: %w", err)
		}
		if err := cmdObj.Start(); err != nil {
			return nil, nil, fmt.Errorf("falha ao iniciar comando local: %w", err)
		}
		waitFn := func() error {
			return cmdObj.Wait()
		}
		return io.MultiReader(stdout, stderr), waitFn, nil
	}

	if c.sshClient == nil {
		return nil, nil, fmt.Errorf("cliente SSH não conectado")
	}

	session, err := c.sshClient.NewSession()
	if err != nil {
		return nil, nil, fmt.Errorf("falha ao criar sessão SSH: %w", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		_ = session.Close()
		return nil, nil, fmt.Errorf("falha ao abrir stdout pipe SSH: %w", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		_ = session.Close()
		return nil, nil, fmt.Errorf("falha ao abrir stderr pipe SSH: %w", err)
	}

	// Cria um pipe síncrono para fazer o merge concorrente de stdout e stderr sem deadlock
	pr, pw := io.Pipe()

	if err := session.Start(finalCmd); err != nil {
		_ = session.Close()
		_ = pr.Close()
		_ = pw.Close()
		return nil, nil, fmt.Errorf("falha ao iniciar comando SSH: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, _ = io.Copy(pw, stdout)
	}()

	go func() {
		defer wg.Done()
		_, _ = io.Copy(pw, stderr)
	}()

	go func() {
		wg.Wait()
		_ = pw.Close()
	}()

	waitFn := func() error {
		defer session.Close()
		defer pr.Close()
		return session.Wait()
	}

	return pr, waitFn, nil
}

// StartCommandInput inicia um comando e retorna seu stdin como io.WriteCloser para recepção de dados via streaming
func (c *Client) StartCommandInput(ctx context.Context, cmd string, useSudo bool) (io.WriteCloser, func() error, error) {
	finalCmd := c.buildFinalCmd(cmd, useSudo)

	if isLocal(c.server) {
		cmdObj := exec.CommandContext(ctx, "bash", "-c", finalCmd)
		stdin, err := cmdObj.StdinPipe()
		if err != nil {
			return nil, nil, fmt.Errorf("falha ao abrir stdin pipe local: %w", err)
		}
		if err := cmdObj.Start(); err != nil {
			return nil, nil, fmt.Errorf("falha ao iniciar comando receptor local: %w", err)
		}
		waitFn := func() error {
			return cmdObj.Wait()
		}
		return stdin, waitFn, nil
	}

	if c.sshClient == nil {
		return nil, nil, fmt.Errorf("cliente SSH não conectado")
	}

	session, err := c.sshClient.NewSession()
	if err != nil {
		return nil, nil, fmt.Errorf("falha ao criar sessão SSH: %w", err)
	}

	var stderrBuf bytes.Buffer
	session.Stderr = &stderrBuf

	stdin, err := session.StdinPipe()
	if err != nil {
		_ = session.Close()
		return nil, nil, fmt.Errorf("falha ao abrir stdin pipe SSH: %w", err)
	}

	if err := session.Start(finalCmd); err != nil {
		_ = session.Close()
		return nil, nil, fmt.Errorf("falha ao iniciar comando receptor SSH: %w", err)
	}

	waitFn := func() error {
		defer session.Close()
		if err := session.Wait(); err != nil {
			if stderrBuf.Len() > 0 {
				return fmt.Errorf("%w: %s", err, stderrBuf.String())
			}
			return err
		}
		return nil
	}

	return stdin, waitFn, nil
}

func (c *Client) buildFinalCmd(cmd string, useSudo bool) string {
	finalCmd := cmd
	if useSudo {
		username := strings.TrimSpace(strings.ToLower(c.server.Username))
		if username == "root" {
			finalCmd = cmd
		} else if strings.TrimSpace(c.server.SudoPassword) != "" {
			finalCmd = fmt.Sprintf("echo %s | sudo -S -p '' bash -c %s", escapeShell(c.server.SudoPassword), escapeShell(cmd))
		} else {
			if !strings.HasPrefix(cmd, "sudo ") {
				finalCmd = fmt.Sprintf("sudo -n bash -c %s", escapeShell(cmd))
			}
		}
	}
	return finalCmd
}

func execLocalCommand(cmd string) (string, error) {
	cmdObj := exec.Command("bash", "-c", cmd)
	out, err := cmdObj.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("comando local falhou (%w): %s", err, string(out))
	}
	return string(out), nil
}
