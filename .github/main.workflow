workflow "Build and Publish" {
  on = "push"
  resolves = "Publish"
}

action "Format" {
  uses = "./.github/actions/go"
  args = "format"
}

action "Lint" {
  uses = "./.github/actions/go"
  args = "lint"
}

action "Test" {
  uses = "./.github/actions/go"
  args = "test"
}

action "Build" {
  needs = ["Format", "Lint", "Test"]
  uses = "./.github/actions/go"
  secrets = ["DOCKER_IMAGE"]
  args = "build"
}

action "If repo was tagged" {
  needs = ["Build"]
  uses = "actions/bin/filter@master"
  args = "tag v*"
}

action "If branch is master" {
  needs = ["Build"]
  uses = "actions/bin/filter@master"
  args = "branch master"
}

action "Publish" {
  needs = ["If repo was tagged"]
  secrets = ["GITHUB_TOKEN"]
  uses = "docker://goreleaser/goreleaser:v0.97"
  args = ["release", "--debug"] 
  }

action "Docker Hub Login" {
  needs = ["If repo was tagged"]
  uses = "actions/docker/login@master"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "Docker Hub Publish" {
  needs = ["Docker Login"]
  uses = "./.github/actions/docker"
  secrets = ["DOCKER_IMAGE"]
  args = ["publish", "Dockerfile"]
}

