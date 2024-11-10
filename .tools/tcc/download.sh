#!/bin/bash

WIN_URL="https://download.savannah.gnu.org/releases/tinycc/tcc-0.9.27-win64-bin.zip"
LINUX_URL="https://download.savannah.gnu.org/releases/tinycc/tcc-0.9.27.tar.bz2"
WIN_TARGET=".tools/tcc/win"
LINUX_TARGET=".tools/tcc/linux"

if [ ! -d ".tools/tcc/win" ]; then
  mkdir -p .tools/tcc/win
  echo "Downloading and extracting Windows version..."
  curl -s -L -o /tmp/tcc-win.zip "$WIN_URL"
  unzip -q /tmp/tcc-win.zip -d "$WIN_TARGET"
  mv "$WIN_TARGET"/tcc*/* "$WIN_TARGET"
  rm -rf "$WIN_TARGET"/tcc
  rm /tmp/tcc-win.zip
else
    echo "Directory .tools/tcc/win already exists. Skipping the remaining steps."
fi

if [ ! -d ".tools/tcc/linux" ]; then
  mkdir -p .tools/tcc/linux
  echo "Downloading and extracting Linux version..."
  curl -s -L -o /tmp/tcc-linux.tar.bz2 "$LINUX_URL"
  tar -xjf /tmp/tcc-linux.tar.bz2 -C "$LINUX_TARGET"
  mv "$LINUX_TARGET"/tcc*/* "$LINUX_TARGET"
  rm -rf "$LINUX_TARGET"/tcc-*
  rm /tmp/tcc-linux.tar.bz2
else
    echo "Directory .tools/tcc/linux already exists. Skipping the remaining steps."
fi