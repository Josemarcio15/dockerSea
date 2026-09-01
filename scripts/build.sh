#!/usr/bin/env bash

set -e

# Garante que o script sempre opere na raiz do projeto
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$ROOT_DIR"

CONFIG_FILE="build/config.yml"

if [ ! -f "$CONFIG_FILE" ]; then
  echo "❌ Arquivo $CONFIG_FILE não encontrado em $ROOT_DIR."
  exit 1
fi

# Extrai a versão completa atual (ex: 0.0.3-alpha)
FULL_VERSION=$(awk '/^  version:/ { gsub(/["'\'']/, "", $2); print $2; exit }' "$CONFIG_FILE")
if [ -z "$FULL_VERSION" ]; then
  FULL_VERSION="0.0.3-alpha"
fi

# Separa a versão numérica e o sufixo/estágio (ex: 0.0.3 e alpha)
CURRENT_NUM="${FULL_VERSION%%-*}"
if [[ "$FULL_VERSION" == *"-"* ]]; then
  CURRENT_STAGE="${FULL_VERSION#*-}"
else
  CURRENT_STAGE="alpha"
fi

echo "=================================================="
echo "          🚀 DockSea - Assistente de Build        "
echo "=================================================="
echo ""

# 1. Pergunta a Versão (autopreenchida com a versão atual)
read -e -p "Versão [$CURRENT_NUM]: " INPUT_NUM
INPUT_NUM="${INPUT_NUM:-$CURRENT_NUM}"

# 2. Pergunta o Estágio (alpha, beta, rc, stable/nenhum)
read -e -p "Estágio (alpha/beta/rc/vazio) [$CURRENT_STAGE]: " INPUT_STAGE
INPUT_STAGE="${INPUT_STAGE:-$CURRENT_STAGE}"

# Monta a nova versão final
if [ -n "$INPUT_STAGE" ] && [ "$INPUT_STAGE" != "none" ] && [ "$INPUT_STAGE" != "stable" ]; then
  NEW_VERSION="${INPUT_NUM}-${INPUT_STAGE}"
else
  NEW_VERSION="${INPUT_NUM}"
fi

echo ""
echo "📌 Versão selecionada: $NEW_VERSION"

# Atualiza no build/config.yml e no App.svelte caso tenha mudado
if [ "$NEW_VERSION" != "$FULL_VERSION" ]; then
  echo "📝 Atualizando $CONFIG_FILE para $NEW_VERSION..."
  sed -i "s/version: .*/version: \"$NEW_VERSION\"/" "$CONFIG_FILE"
  
  if [ -f "frontend/src/App.svelte" ]; then
    echo "📝 Atualizando versão no frontend (App.svelte)..."
    sed -i "s/let appVersion = .*/let appVersion = \"$NEW_VERSION\";/" "frontend/src/App.svelte"
  fi

  echo "🔄 Atualizando build assets com wails3..."
  wails3 task common:update:build-assets || true
fi

echo ""
echo "Selecione as plataformas para compilação:"
echo "  1 - Linux (.deb)"
echo "  2 - Linux (.rpm)"
echo "  3 - Linux (.AppImage)"
echo "  4 - Standalone (Linux binário)"
echo "  5 - Windows (.exe)"
echo "  6 - Todos os formatos (DEB, RPM, AppImage, Standalone e Windows .exe)"
echo "  7 - Cancelar"
echo ""

read -p "Opção [1-7]: " PLAT_CHOICE

echo ""
case "$PLAT_CHOICE" in
  1)
    echo "📦 Compilando Linux DEB..."
    wails3 task linux:create:deb
    mv -f bin/DockSea.deb "bin/DockSea_${NEW_VERSION}_linux_amd64.deb" 2>/dev/null || true
    echo "✅ Concluído! Pacote gerado: bin/DockSea_${NEW_VERSION}_linux_amd64.deb"
    ;;
  2)
    echo "📦 Compilando Linux RPM..."
    wails3 task linux:create:rpm
    mv -f bin/DockSea.rpm "bin/DockSea-${NEW_VERSION}.x86_64.rpm" 2>/dev/null || true
    echo "✅ Concluído! Pacote gerado: bin/DockSea-${NEW_VERSION}.x86_64.rpm"
    ;;
  3)
    echo "📦 Compilando Linux AppImage..."
    wails3 task linux:create:appimage
    mv -f bin/docksea-x86_64.AppImage "bin/DockSea-${NEW_VERSION}-x86_64.AppImage" 2>/dev/null || true
    echo "✅ Concluído! Pacote gerado: bin/DockSea-${NEW_VERSION}-x86_64.AppImage"
    ;;
  4)
    echo "📦 Compilando Linux Standalone..."
    wails3 task linux:build
    mv -f bin/DockSea "bin/DockSea_${NEW_VERSION}_linux_amd64" 2>/dev/null || true
    echo "✅ Concluído! Executável gerado: bin/DockSea_${NEW_VERSION}_linux_amd64"
    ;;
  5)
    echo "📦 Compilando Windows (.exe)..."
    wails3 task windows:build
    mv -f bin/DockSea.exe "bin/DockSea_${NEW_VERSION}_windows_amd64.exe" 2>/dev/null || true
    echo "✅ Concluído! Executável gerado: bin/DockSea_${NEW_VERSION}_windows_amd64.exe"
    ;;
  6)
    echo "📦 Compilando Todos os formatos..."
    echo "-> 1/4 Linux Standalone..."
    wails3 task linux:build
    mv -f bin/DockSea "bin/DockSea_${NEW_VERSION}_linux_amd64" 2>/dev/null || true

    echo "-> 2/4 Linux Pacotes (.deb, .rpm, .AppImage)..."
    wails3 task linux:create:deb
    mv -f bin/DockSea.deb "bin/DockSea_${NEW_VERSION}_linux_amd64.deb" 2>/dev/null || true

    wails3 task linux:create:rpm || true
    mv -f bin/DockSea.rpm "bin/DockSea-${NEW_VERSION}.x86_64.rpm" 2>/dev/null || true

    wails3 task linux:create:appimage || true
    mv -f bin/docksea-x86_64.AppImage "bin/DockSea-${NEW_VERSION}-x86_64.AppImage" 2>/dev/null || true

    echo "-> 3/4 Windows (.exe)..."
    wails3 task windows:build
    mv -f bin/DockSea.exe "bin/DockSea_${NEW_VERSION}_windows_amd64.exe" 2>/dev/null || true

    # Limpa eventuais arquivos temporários/residuais sem versão
    rm -f bin/DockSea bin/DockSea.exe bin/DockSea.deb bin/DockSea.rpm bin/docksea-x86_64.AppImage 2>/dev/null || true

    echo "✅ Todos os builds foram gerados na pasta ./bin exclusivamente com as versões no nome!"
    ;;
  7)
    echo "Operação cancelada."
    exit 0
    ;;
  *)
    echo "Opção inválida."
    exit 1
    ;;
esac

echo ""
echo "📂 Arquivos em ./bin:"
ls -lh bin/
