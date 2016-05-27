#!/usr/bin/env bash
# we want to allow devs to use their own hooks, so if they add one, use that as link
# get all existing non symlinked hooks and append .local
# symlink

# our project level hooks

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GIT="$(git rev-parse --show-toplevel)"

echo "script running from: $DIR"
echo "in git repo at: $GIT"

if [ ! -e "${GIT}/.git/hooks/pre-push" ]; then
    ln -s "$DIR/pre-push" "$GIT/.git/hooks/pre-push"
else
    echo "cannot add pre-push githook - already exists"
fi

