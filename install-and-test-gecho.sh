#!/bin/bash

set -e

echo "🧪 Building and installing Gecho CLI"
go install .

echo "🔍 Verifying installed CLI"
which gecho || { echo "❌ gecho not found in PATH"; exit 1; }

echo "🧪 Running version check (should show help)"
gecho --help

echo "✅ Gecho CLI installed and ready."
