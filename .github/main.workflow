workflow "build" {
  on = "push"
  resolves = ["docker://circleci/golang:1.12-1"]
}

action "docker://circleci/golang:1.12" {
  uses = "docker://circleci/golang:1.12"
  args = "setup"
  runs = "make"
}

action "docker://circleci/golang:1.12-1" {
  uses = "docker://circleci/golang:1.12"
  needs = ["docker://circleci/golang:1.12"]
  runs = "make"
  args = "test"
}
