#!/usr/bin/env bash
# update-tari-grpc.sh — regenerate go-tari-grpc-lib's protos/Go bindings
# from a given upstream Tari release tag.
#
# WHY THIS EXISTS: go-tari-grpc-lib wraps Tari's application-facing GRPC
# interface (applications/minotari_app_grpc/proto in the tari-project/tari
# repo). That interface changes with every Tari release and the wrapper repo
# has historically drifted badly out of sync (confirmed 2026-08-16: base_node
# and wallet protos were 800+ line diffs behind current source). This script
# makes "update to release vX.Y.Z" a repeatable one-command operation instead
# of a manual proto-by-proto diff exercise.
#
# WHAT IT DOES:
#   1. Shallow-clones tari-project/tari at the given tag into a scratch dir.
#   2. Copies the 8 proto files go-tari-grpc-lib actually consumes from
#      applications/minotari_app_grpc/proto/ into a work copy of the repo's
#      tari_protos/ dir, preserving the existing go_package override (the
#      upstream protos don't set it — go-tari-grpc-lib's local copies do,
#      pointing at its own tari_generated/ package. This script re-applies
#      that override after copying, since a raw `cp` would drop it).
#   3. Runs protoc (via mise-managed protoc/protoc-gen-go/protoc-gen-go-grpc)
#      to regenerate tari_generated/*.pb.go and *_grpc.pb.go for the 3 protos
#      that declare a `service` (base_node, p2pool, wallet) plus plain
#      message-only .pb.go for the rest.
#   4. Diffs the regenerated output against what's currently committed and
#      reports what changed — does NOT auto-commit. Review the diff, run
#      `go build ./...` in the repo, then commit yourself.
#
# USAGE:
#   ./update-tari-grpc.sh <tari-release-tag> [path-to-go-tari-grpc-lib-clone]
#
# EXAMPLE:
#   ./update-tari-grpc.sh v5.6.0 /workspace/go-tari-grpc-lib
#
# REQUIRES: git, protoc, protoc-gen-go, protoc-gen-go-grpc on PATH (all
# mise-managed on the hermes-devbench-pattern image — see that repo's
# .mise.toml). Verified against protoc 29.1 / protoc-gen-go 1.36.x /
# protoc-gen-go-grpc 1.6.x.
#
# KNOWN LIMITATION: the set of 8 proto files and the 3 with `service`
# definitions is HARDCODED below, discovered by inspecting the repo as of
# 2026-08-16. If Tari adds/removes a proto file from
# applications/minotari_app_grpc/proto/ in a future release, this script
# will not pick it up automatically — it will silently skip a new file or
# fail cleanly on a removed one (missing import error from protoc). Check
# the PROTO_FILES/SERVICE_FILES arrays below against the upstream directory
# listing if a run behaves unexpectedly after a big Tari version jump.

set -euo pipefail

TAG="${1:?usage: update-tari-grpc.sh <tari-release-tag> [path-to-go-tari-grpc-lib-clone]}"
REPO_PATH="${2:-$(pwd)}"

if [ ! -f "$REPO_PATH/go.mod" ] || ! grep -q "go-tari-grpc-lib" "$REPO_PATH/go.mod"; then
  echo "error: $REPO_PATH doesn't look like a go-tari-grpc-lib checkout (no go.mod matching module name)" >&2
  exit 1
fi

for tool in git protoc protoc-gen-go protoc-gen-go-grpc; do
  command -v "$tool" >/dev/null 2>&1 || { echo "error: $tool not found on PATH" >&2; exit 1; }
done

SCRATCH_DIR="$(mktemp -d)"
trap 'rm -rf "$SCRATCH_DIR"' EXIT

echo "==> Cloning tari-project/tari @ $TAG (shallow)..."
git clone --depth 1 --branch "$TAG" https://github.com/tari-project/tari.git "$SCRATCH_DIR/tari" 2>&1 | tail -3

UPSTREAM_PROTO_DIR="$SCRATCH_DIR/tari/applications/minotari_app_grpc/proto"
if [ ! -d "$UPSTREAM_PROTO_DIR" ]; then
  echo "error: $UPSTREAM_PROTO_DIR doesn't exist in this tag — Tari may have moved the proto dir. Update this script." >&2
  exit 1
fi

# The 8 protos go-tari-grpc-lib consumes (hardcoded — see KNOWN LIMITATION above).
PROTO_FILES=(base_node block network p2pool sidechain_types transaction types wallet)
# Of those, the 3 that declare a `service` (need _grpc.pb.go too).
SERVICE_FILES=(base_node p2pool wallet)

WORK_PROTO_DIR="$SCRATCH_DIR/work_protos"
mkdir -p "$WORK_PROTO_DIR"

echo "==> Copying protos and re-applying go_package override..."
for f in "${PROTO_FILES[@]}"; do
  src="$UPSTREAM_PROTO_DIR/$f.proto"
  if [ ! -f "$src" ]; then
    echo "error: upstream proto $f.proto not found at $src — tag $TAG may have renamed/removed it. Update PROTO_FILES." >&2
    exit 1
  fi
  # Strip any existing go_package option (upstream Tari protos don't set
  # one for this dir, but be defensive), then append go-tari-grpc-lib's own.
  grep -v '^option go_package' "$src" > "$WORK_PROTO_DIR/$f.proto"
  # Insert the go_package option right after the `package` line so it's
  # valid proto syntax (option statements must follow package/syntax).
  awk -v pkg='option go_package = "github.com/Snipa22/go-tari-grpc-lib/tari_generated";' \
    '{print} /^package /{print pkg}' "$WORK_PROTO_DIR/$f.proto" > "$WORK_PROTO_DIR/$f.proto.tmp"
  mv "$WORK_PROTO_DIR/$f.proto.tmp" "$WORK_PROTO_DIR/$f.proto"
done

echo "==> Regenerating Go bindings..."
GEN_DIR="$SCRATCH_DIR/generated"
mkdir -p "$GEN_DIR"

for f in "${PROTO_FILES[@]}"; do
  IS_SERVICE=0
  for s in "${SERVICE_FILES[@]}"; do
    [ "$s" = "$f" ] && IS_SERVICE=1
  done

  if [ "$IS_SERVICE" = "1" ]; then
    protoc \
      --proto_path="$WORK_PROTO_DIR" \
      --go_out="$GEN_DIR" --go_opt=paths=source_relative \
      --go-grpc_out="$GEN_DIR" --go-grpc_opt=paths=source_relative \
      "$WORK_PROTO_DIR/$f.proto"
  else
    protoc \
      --proto_path="$WORK_PROTO_DIR" \
      --go_out="$GEN_DIR" --go_opt=paths=source_relative \
      "$WORK_PROTO_DIR/$f.proto"
  fi
done

echo "==> Diffing against current repo state..."
DIFF_FOUND=0

echo ""
echo "--- tari_protos/ ---"
for f in "${PROTO_FILES[@]}"; do
  if ! diff -q "$REPO_PATH/tari_protos/$f.proto" "$WORK_PROTO_DIR/$f.proto" >/dev/null 2>&1; then
    echo "CHANGED: tari_protos/$f.proto"
    DIFF_FOUND=1
  fi
done

echo ""
echo "--- tari_generated/ ---"
for f in "${PROTO_FILES[@]}"; do
  if ! diff -q "$REPO_PATH/tari_generated/$f.pb.go" "$GEN_DIR/$f.pb.go" >/dev/null 2>&1; then
    echo "CHANGED: tari_generated/$f.pb.go"
    DIFF_FOUND=1
  fi
  for s in "${SERVICE_FILES[@]}"; do
    if [ "$s" = "$f" ]; then
      if ! diff -q "$REPO_PATH/tari_generated/${f}_grpc.pb.go" "$GEN_DIR/${f}_grpc.pb.go" >/dev/null 2>&1; then
        echo "CHANGED: tari_generated/${f}_grpc.pb.go"
        DIFF_FOUND=1
      fi
    fi
  done
done

if [ "$DIFF_FOUND" = "0" ]; then
  echo ""
  echo "==> No changes detected — go-tari-grpc-lib is already current with $TAG."
  exit 0
fi

echo ""
echo "==> Applying changes to $REPO_PATH..."
cp "$WORK_PROTO_DIR"/*.proto "$REPO_PATH/tari_protos/"
cp "$GEN_DIR"/*.pb.go "$REPO_PATH/tari_generated/"

echo ""
echo "==> Done. Review the diff, then in $REPO_PATH:"
echo "      git diff"
echo "      go build ./..."
echo "      go vet ./..."
echo "    Check nodeGRPC/ and walletGRPC/ (the hand-written wrapper layer) —"
echo "    if any message/service shape changed, those wrappers may need"
echo "    matching updates. protoc regeneration does NOT touch them."
