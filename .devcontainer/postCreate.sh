#!/bin/bash

# enable zsh completion
echo 'deb http://download.opensuse.org/repositories/shells:/zsh-users:/zsh-completions/Debian_11/ /' | sudo tee /etc/apt/sources.list.d/shells:zsh-users:zsh-completions.list
curl -fsSL https://download.opensuse.org/repositories/shells:zsh-users:zsh-completions/Debian_11/Release.key | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/shells_zsh-users_zsh-completions.gpg >/dev/null
sudo apt update
sudo apt install zsh-completions
