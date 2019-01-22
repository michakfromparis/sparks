#!/bin/bash

# TODO Turn this script into an action and rpc call the github api to run
# the action with parameters from this script that will shrink to a two liner

set -e
set -x

# current git branch
branch=$(git symbolic-ref HEAD | sed -e 's,.*/\(.*\),\1,')

# current project name
projectName=$(git config --local remote.origin.url|sed -n 's#.*/\([^.]*\)\.git#\1#p')

# set master branch
masterBranch=$branch

version="$1"

newVersion="$2"

version ()
{
	# master branch validation
	if [ $branch != "master" ]; then
		echo "needs to run on master branch"
		exit -1
	fi

	if [ "$version" = "" ]; then
		echo "Enter the release version number without a leading v"
		read version
	fi

	if [ "$newVersion" = "" ]; then
		echo "Enter version number for $projectName after the release"
		read newVersion
	fi

	# v1.0.0, v1.7.8, etc..
	versionLabel="v$version"

	# establish branch and tag name variables
	releaseBranch=release-$version
	tagName=$versionLabel

	echo "Started releasing $versionLabel for $projectName ..."

	# pull the latest version of the code from master
	git pull

	# file in which to update version number
	versionFile="version.go"
 
	# find version number ("1.5.5" for example)
	# and replace it with newly specified version number
	sed -i.backup -E "s/[0-9.]+[0-9]+/$version/" $versionFile $versionFile

	# remove backup file created by sed command
	rm -f $versionFile.backup
	
	# adding the changes to git
	git add $versionFile

	# commit version number increment
	git commit -m "[Pre Release] Setting version to $versionLabel, Creating branch $releaseBranch"
}

version2 ()
{
		# create the release branch from the -master branch
	git checkout -b $releaseBranch $masterBranch

	# push local releaseBranch to remote
	git push -u origin $releaseBranch

	echo "$projectName $versionLabel successfully released and labeled as $versionLabel"
	echo "Incrementing $masterBranch version for new builds after the release"

	# checkout to master branch
	git checkout $masterBranch

	# pull the latest version of the code from master
	git pull
	
	# find version number ("1.5.5" for example)
	# and replace it with newly specified version number
	sed -i.backup -E "s/[0-9.]+[0-9]+/$newVersion/" $versionFile $versionFile

	# remove backup file created by sed command
	rm -f $versionFile.backup
	
	# adding the changes to git
	git add $versionFile

	# Commit setting new master branch version	
	git commit -m "[Post Release] Setting version to $newVersionNumer"

	# push commit to remote origin
	git push

}

tag ()
{
	# create tag for new version from -master
	git tag -f $tagName
}

push ()
{
	# push tag to remote origin. This will actually trigger CI
	git push -f --tags origin
}

release ()
{
	local version=$1
	local newVersion=$2

	echo "[$1] -> [$2]"
	version "$version"
	tag "v$version"
	push
	if [ post = true ]; then
		version2 newVersion
		push
	fi
}

release $1 $2