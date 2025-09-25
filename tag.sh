#!/bin/bash

set -e -u -o pipefail


MAX_PRE_VERSION_COUNT=50


GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

function green() {
	echo -e "${GREEN}$*${NC}"
}

function yellow() {
	echo -e "${YELLOW}$*${NC}"
}

function blue() {
	echo -e "${BLUE}$*${NC}"
}

function die () {
	echo -e "${RED}${1}${NC}" >&2
	exit 1
}

BRANCH=${DRONE_COMMIT_BRANCH:-$(git rev-parse --abbrev-ref HEAD)}


calc_version () {
	VERSION_FILE=$1
	DIR=$(dirname "$VERSION_FILE")
	if [ "$DIR" = "." ]; then
		MODULE_PATH=
	else
		MODULE_PATH=$DIR/
	fi

	echo "${MODULE_PATH}v$(cat "$VERSION_FILE")"
}

VERSION_FILES=$(find . -name "Version" -or -name "Versions" | cut -d/ -f2-)
if [[ -z "${VERSION_FILES// /}" ]]; then
	echo "No Version of Versions file found - skipping this step..."
	exit
fi

# Would do this, except that Drone is very weird. It seems to create a master branch for some reason,
# even though there is no master branch in this repo.
# DEFAULT_BRANCH=$(git rev-parse --abbrev-ref origin/HEAD | cut -c8-)
DEFAULT_BRANCH=$( git remote show $( git config --get remote.origin.url ) | grep 'HEAD branch' | cut -d' ' -f5 )
echo "Default branch is: $DEFAULT_BRANCH"

if [ "$1" = "check_version" ]; then
    blue "\\nChecking version and verify if we need to update the version file."
    shift
    VERSION_FILES=$*
    for version_file in $VERSION_FILES
    do
        VERSION=$(calc_version "$version_file")
        VER_EXIST=$(git tag -l "$VERSION")

        if [ -n "$VER_EXIST" ]
        then
            echo "Version ${VERSION} already tagged"
            if [ "$BRANCH" != "$DEFAULT_BRANCH" ]
            then
                die "Need to update version file: $version_file - exiting"
            fi
        else
            echo "Version ${VERSION} not tagged"
        fi

    done
    green "Version file(s) look good!"
    exit 0
fi


blue "\\nChecking out and pulling the $BRANCH branch"
git checkout "$BRANCH" > /dev/null
git pull --rebase origin "$BRANCH" > /dev/null

blue "\\nLooking for the next version tag"
FOUND=no
for version_file in $VERSION_FILES
do
	VERSION=$(calc_version "$version_file")

    # First sed command attempts to extract DEV-N from the beginning of the branch name, otherwise it leaves the branch name unchanged.
    # Second sed command replaces all sequences of characters that aren't letters or digits with a single dash.
    PRERELEASE=`echo $BRANCH | sed -e 's/^.*\(DEV-[0-9]\+\).*$/\1/' -e 's/[^0-9a-zA-Z]\+/-/g'`

	if [ "$BRANCH" != "$DEFAULT_BRANCH" ]; then

		COUNTER=1
		while [  $COUNTER -lt $MAX_PRE_VERSION_COUNT ]; do
			tag="$VERSION-$PRERELEASE.$COUNTER"
            echo "Trying $tag..."

			# Check if the tag exist
			if [ -n "$(git tag -l "$tag")" ]
			then
                yellow "Tag exists!"

                # might be that another developer has pushed out the same version number in another branch that hasn't been merged in yet.
                # Let's check and make sure that is not the case
                tag_exists_in_other_branches=no
                raw_version=$(cat "$version_file")
                all_revs=$(git rev-list --all)
                # This looks for all commits, across all branches, where that particular version appears in that version file
                for commit in $(git grep "$raw_version" $all_revs -- "$version_file" | cut -d: -f1 | sort -u)
                do
                    # Let's see if any of these commits appear in other branches than our own
                    for branch in $(git --no-pager branch --contains "$commit" --all | rev | cut -d" " -f1 | rev)
                    do
                        if [ "$branch" != "$BRANCH" ]
                        then
                            tag_exists_in_other_branches=yes
                        fi
                    done
                done

                if [ "tag_exists_in_other_branches" == "yes" ]
                then
				    die "Tag does exist, but not as a part of your branch (another developer working in parallel?). Consider bumping the version number."
                fi
            else
                break
			fi

			(( COUNTER=COUNTER+1 ))

		done

		VERSION=$tag
	fi
	FOUND=yes
	green "About to tag $VERSION"
	git tag "$VERSION"
done

if [ "$FOUND" == "no" ]
then
	die "Did not find a new version tag!"
fi

green "About to push!"
git push --no-verify --tags origin "$BRANCH"

