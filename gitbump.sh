set -e

oldv=$(git tag --sort=-refname --list "v[0-9]*" | head -n 1)
echo "oldv: $oldv"

newv=$(docker run --rm -v "$PWD":/app treeder/bump --input "$oldv" patch)
echo "newv: $newv"

git tag -a "v$newv" -m "version $newv"
git push --follow-tags
