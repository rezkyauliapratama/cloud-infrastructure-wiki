#!/bin/sh

set -e

# Install kustomization
if ! command -v kustomize &> /dev/null
then
  echo "Installing kustomize ..."
  brew install kustomize
fi

# Install kubernetes-split-yaml
if ! command -v kubernetes-split-yaml &> /dev/null
then
  echo "Installing kubernetes-split-yaml ..."
  go install -v github.com/mogensen/kubernetes-split-yaml@v0.4.0
fi

find platform/ -maxdepth 1 -mindepth 1 -type d -exec sh -c '
  # only generate the kustomization that contains #kustomization:render at the top
  header=$(cat "{}/kustomization.yaml" | head -n1)
  if [[ -f "{}/kustomization.yaml" && $header =~ ^\#\s?kustomization\:render ]]; then
    # generate remote kubernetes manifest locally
    kustomize build --enable-helm {} > {}/00_install.yaml

    # # split larger kubernetes manifest file to separated files
    rm -rf {}/default
    filenames=$(kubernetes-split-yaml --outdir {}/default {}/00_install.yaml 2>&1 | grep -o "[A-z0-9\.\:\-]*\.yaml" | sed "s/^/- /g")
    # mkdir {}/default
    touch {}/default/kustomization.yaml
    echo "---\nresources:\n$filenames" > {}/default/kustomization.yaml
    rm -rf {}/00_install.yaml
  fi
' {} \;
