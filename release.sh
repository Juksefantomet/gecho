#!/bin/bash
set -e

VERSION=$(cat VERSION)
TAG="v$VERSION"

if [[ -z "$VERSION" ]]; then
  echo "✗ VERSION file is empty"
  exit 1
fi

# Check for clean working directory
if [[ -n $(git status --porcelain) ]]; then
  echo "✗ Working directory is not clean. Commit or stash changes first."
  exit 1
fi

# Confirm version
echo "Preparing to tag release:"
echo "  VERSION file: $VERSION"
echo "  Git tag:      $TAG"
echo

read -p "Proceed with tagging and pushing? (y/N) " confirm
if [[ "$confirm" != "y" ]]; then
  echo "Aborted."
  exit 0
fi

# Create and push the tag
git tag "$TAG"
git push origin "$TAG"

echo
echo "✓ Tagged and pushed: $TAG"
echo "→ Create GitHub release here:"
echo "  https://github.com/Juksefantomet/gecho/releases/new?tag=$TAG"
