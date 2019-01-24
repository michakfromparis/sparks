workflow "Build, test and Publish" {
  on = "push"
  resolves = ["Publish to Github Release", "Publish to Docker Hub"]
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

action "Publish to Github Release" {
  needs = ["If repo was tagged"]
  secrets = ["GITHUB_TOKEN"]
  uses = "docker://goreleaser/goreleaser:v0.97"
  args = ["release", "--debug"] 
  }

action "Login to Docker Hub" {
  needs = ["If repo was tagged"]
  uses = "actions/docker/login@master"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "Publish to Docker Hub" {
  needs = ["Docker Hub Login"]
  uses = "./.github/actions/docker"
  secrets = ["DOCKER_IMAGE"]
  args = ["publish", "Dockerfile"]
}

