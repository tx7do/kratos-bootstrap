#!/usr/bin/env bash
set -euo pipefail

# ============================================================================
# upgrade.sh — Recursively upgrade all Go module dependencies in the monorepo.
# Traverses every directory containing a go.mod file and runs:
#   go get -u ./...
#   go mod tidy
# ============================================================================

# Color codes (disabled if stdout is not a terminal).
if [ -t 1 ]; then
    C_GREEN='\033[0;32m'
    C_RED='\033[0;31m'
    C_YELLOW='\033[0;33m'
    C_BLUE='\033[0;34m'
    C_RESET='\033[0m'
else
    C_GREEN=''
    C_RED=''
    C_YELLOW=''
    C_BLUE=''
    C_RESET=''
fi

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
FAILED_DIRS=()
TOTAL=0
SUCCESS=0

upgrade_module() {
    local mod_dir="$1"
    local rel_path="${mod_dir#$ROOT_DIR/}"
    [ "$rel_path" = "$ROOT_DIR" ] && rel_path="(root)"

    printf "${C_BLUE}[upgrade]${C_RESET} %s\n" "$rel_path"

    local output

    # go get -u ./...
    output=$(cd "$mod_dir" && go get -u ./... 2>&1) || {
        printf "${C_RED}  ✗ go get -u failed:${C_RESET}\n%s\n" "$output"
        FAILED_DIRS+=("$rel_path (go get)")
        return 1
    }

    # go mod tidy
    output=$(cd "$mod_dir" && go mod tidy 2>&1) || {
        printf "${C_RED}  ✗ go mod tidy failed:${C_RESET}\n%s\n" "$output"
        FAILED_DIRS+=("$rel_path (go mod tidy)")
        return 1
    }

    printf "${C_GREEN}  ✓ done${C_RESET}\n"
    return 0
}

# Find all directories containing go.mod (including root).
while IFS= read -r -d '' mod_file; do
    mod_dir="$(dirname "$mod_file")"
    TOTAL=$((TOTAL + 1))

    if upgrade_module "$mod_dir"; then
        SUCCESS=$((SUCCESS + 1))
    fi
done < <(find "$ROOT_DIR" -name go.mod -not -path '*/.git/*' -not -path '*/.idea/*' -print0)

# Summary.
echo ""
printf "${C_BLUE}=== Summary ===${C_RESET}\n"
printf "Total modules: %d\n" "$TOTAL"
printf "${C_GREEN}Succeeded:     %d${C_RESET}\n" "$SUCCESS"

if [ "${#FAILED_DIRS[@]}" -gt 0 ]; then
    printf "${C_RED}Failed:        %d${C_RESET}\n" "${#FAILED_DIRS[@]}"
    echo ""
    printf "${C_RED}Failed modules:${C_RESET}\n"
    for dir in "${FAILED_DIRS[@]}"; do
        printf "  ✗ %s\n" "$dir"
    done
    exit 1
else
    printf "${C_GREEN}All modules upgraded successfully!${C_RESET}\n"
fi