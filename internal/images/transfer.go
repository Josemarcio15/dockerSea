package images

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"go-walis/internal/core/connection"
)

// TransferImages transfere imagens de um servidor origem para um destino via pipeline direto de streaming (io.Pipe)
func TransferImages(ctx context.Context, srcClient, dstClient *connection.Client, imageIDs []string, onProgress func(progress DockerTransferProgress)) error {
	if len(imageIDs) == 0 {
		return fmt.Errorf("nenhuma imagem fornecida para transferência")
	}

	// 0. Calcular o tamanho real de transferência das imagens selecionadas na origem
	var totalExpectedBytes int64
	for _, idOrTag := range imageIDs {
		trimmed := strings.TrimSpace(idOrTag)
		out, err := srcClient.ExecCommand(fmt.Sprintf("docker image inspect %s --format '{{.Size}}'", trimmed), false)
		if err == nil {
			var sz int64
			if _, scanErr := fmt.Sscanf(strings.TrimSpace(out), "%d", &sz); scanErr == nil && sz > 0 {
				totalExpectedBytes += sz
				continue
			}
		}
	}

	// Fallback para ListImages se o inspect não retornar
	if totalExpectedBytes == 0 {
		allImages, _ := ListImages(srcClient, true)
		for _, idOrTag := range imageIDs {
			trimmed := strings.TrimSpace(idOrTag)
			for _, img := range allImages {
				key := img.Repo + ":" + img.Tag
				if trimmed == key || trimmed == img.ID || trimmed == img.Repo {
					totalExpectedBytes += img.RawSizeBytes
					break
				}
			}
		}
	}

	imagesArg := strings.Join(imageIDs, " ")
	cmdSave := fmt.Sprintf("docker save %s", imagesArg)
	cmdLoad := "docker load"

	if onProgress != nil {
		onProgress(DockerTransferProgress{
			Stage:               "preparing",
			Message:             "Iniciando canais de streaming na origem e destino...",
			FormattedBytes:      "0 B",
			TotalBytes:          totalExpectedBytes,
			FormattedTotalBytes: formatBytes(totalExpectedBytes),
			Speed:               "0 B/s",
			Percent:             0,
		})
	}

	// 1. Inicia comando na origem que gera o stream de saída (docker save direto)
	srcStdout, waitSrc, err := srcClient.StartCommandOutput(ctx, cmdSave, false)
	if err != nil {
		return fmt.Errorf("falha ao iniciar stream na origem: %w", err)
	}

	// 2. Inicia comando no destino que consome o stream de entrada (docker load direto)
	dstStdin, waitDst, err := dstClient.StartCommandInput(ctx, cmdLoad, false)
	if err != nil {
		return fmt.Errorf("falha ao iniciar receptor no destino: %w", err)
	}

	// 3. Pipeline de cópia com medição de velocidade e bytes trafegados reais
	progressReader := &countingReader{
		reader:     srcStdout,
		totalLimit: totalExpectedBytes,
		onProgress: func(bytesSent int64, speed string, percent int) {
			if onProgress != nil {
				onProgress(DockerTransferProgress{
					Stage:               "transferring",
					BytesSent:           bytesSent,
					FormattedBytes:      formatBytes(bytesSent),
					TotalBytes:          totalExpectedBytes,
					FormattedTotalBytes: formatBytes(totalExpectedBytes),
					Speed:               speed,
					Percent:             percent,
					Message:             fmt.Sprintf("Transferindo %s de %s...", formatBytes(bytesSent), formatBytes(totalExpectedBytes)),
				})
			}
		},
	}

	// Executes a cópia em stream direto com buffer de alta performance de 128KB
	copyBuf := make([]byte, 128*1024)
	_, copyErr := io.CopyBuffer(dstStdin, progressReader, copyBuf)
	_ = dstStdin.Close() // Closes o stdin para sinalizar EOF ao docker load

	if copyErr != nil {
		return fmt.Errorf("erro durante streaming de dados: %w", copyErr)
	}

	if onProgress != nil {
		onProgress(DockerTransferProgress{
			Stage:               "loading",
			BytesSent:           progressReader.totalBytes,
			FormattedBytes:      formatBytes(progressReader.totalBytes),
			TotalBytes:          totalExpectedBytes,
			FormattedTotalBytes: formatBytes(totalExpectedBytes),
			Speed:               "0 B/s",
			Percent:             98,
			Message:             "Finalizando importação de metadados no destino...",
		})
	}

	// Aguarda os processos remotos/locais finalizarem
	if waitErr := waitDst(); waitErr != nil {
		return fmt.Errorf("falha na importação no destino: %w", waitErr)
	}

	_ = waitSrc()

	if onProgress != nil {
		onProgress(DockerTransferProgress{
			Stage:               "complete",
			BytesSent:           progressReader.totalBytes,
			FormattedBytes:      formatBytes(progressReader.totalBytes),
			TotalBytes:          totalExpectedBytes,
			FormattedTotalBytes: formatBytes(totalExpectedBytes),
			Percent:             100,
			Message:             "Concluído com sucesso!",
		})
	}

	return nil
}

type countingReader struct {
	reader     io.Reader
	totalBytes int64
	totalLimit int64
	lastBytes  int64
	lastTime   time.Time
	onProgress func(bytesSent int64, speed string, percent int)
}

func (cr *countingReader) Read(p []byte) (int, error) {
	n, err := cr.reader.Read(p)
	if n > 0 {
		cr.totalBytes += int64(n)
		now := time.Now()
		if cr.lastTime.IsZero() {
			cr.lastTime = now
		}

		elapsed := now.Sub(cr.lastTime)
		if elapsed >= 400*time.Millisecond {
			bytesDiff := cr.totalBytes - cr.lastBytes
			speedBytesSec := float64(bytesDiff) / elapsed.Seconds()
			speedStr := formatBytes(int64(speedBytesSec)) + "/s"

			percent := 0
			if cr.totalLimit > 0 {
				percent = int((float64(cr.totalBytes) / float64(cr.totalLimit)) * 100)
				if percent > 100 {
					percent = 100
				}
			} else {
				percent = int(10 + (float64(cr.totalBytes)/(1024*1024*300))*80)
				if percent > 95 {
					percent = 95
				}
			}

			if cr.onProgress != nil {
				cr.onProgress(cr.totalBytes, speedStr, percent)
			}

			cr.lastBytes = cr.totalBytes
			cr.lastTime = now
		}
	}
	return n, err
}

func formatBytes(bytes int64) string {
	if bytes <= 0 {
		return "0 B"
	}
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
