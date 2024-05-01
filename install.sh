#!/bin/bash

set -euo pipefail

version=${CURL2SHELL_VERSION:-"0.0.1"}
# shellcheck disable=SC2034
url_download_Darwin_x86_64="https://github.com/curl2shell/cli/releases/download/v${version}/curl2shell_Darwin_x86_64"
# shellcheck disable=SC2034
url_download_Darwin_arm64="https://github.com/curl2shell/cli/releases/download/v${version}/curl2shell_Darwin_arm64"
# shellcheck disable=SC2034
url_download_Linux_i686="https://github.com/curl2shell/cli/releases/download/v${version}/curl2shell_Linux_i386"
# shellcheck disable=SC2034
url_download_Linux_x86_64="https://github.com/curl2shell/cli/releases/download/v${version}/curl2shell_Linux_x86_64"
# shellcheck disable=SC2034
url_download_Windows_i686="https://github.com/curl2shell/cli/releases/download/v${version}/curl2shell_Windows_i386.exe"
# shellcheck disable=SC2034
url_download_Windows_x86_64="https://github.com/curl2shell/cli/releases/download/v${version}/curl2shell_Windows_x86_64.exe"

main() {
  platform=$(uname -s)
  arch=$(uname -m)

  if [[ $platform == CYGWIN* ]] || [[ $platform == MINGW* ]] || [[ $platform == MSYS* ]]; then
    platform="Windows"
  fi

  install_dir=${INSTALL_DIR:-"/usr/local/bin"}
  install_path=${INSTALL_PATH:-${install_dir}/curl2shell}

  download_url_lookup="url_download_${platform}_${arch}"
  download_url=${!download_url_lookup:-}

  echo "This script will automatically install curl2shell ${version} for you."
  echo "Installation path: ${install_path}"
  if [ "$(id -u)" == "0" ]; then
    echo "Warning: this script is currently running as root. This is dangerous. "
    echo "         Instead run it as normal user. We will sudo as needed."
  fi

  if [ -z "$download_url" ]; then
    echo "error: your platform and architecture (${platform}-${arch}) is unsupported."
    exit 1
  fi

  if ! hash curl 2> /dev/null; then
    echo "error: you do not have 'curl' installed which is required for this script."
    exit 1
  fi

  tmp_file=$(mktemp "${TMPDIR:-/tmp}/.curl2shell.XXXXXXXX")

  curl -fsSL "$download_url" > "$tmp_file"
  chmod +x "$tmp_file"
  if ! mv "$tmp_file" "$install_path" 2> /dev/null; then
    sudo -k mv "$tmp_file" "$install_path"
  fi

  echo 'Done!'
}

main
