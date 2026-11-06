variable "REGISTRY" { default = "registry.localhost:5000/simc" }
variable "TAG" { default = "dev" }

group "default" {
  targets = ["app"]
}

target "app" {
  name = "${SERVICE}"
  matrix = { SERVICE = ["api", "write-result", "run-sim"] }
  args = { SERVICE = "${SERVICE}" }
  target = "prod"
  tags = ["${REGISTRY}/${SERVICE}:${TAG}"]
  labels = {
    "org.opencontainers.image.source" = "https://github.com/RobertsMJ/simc-cloud"
  }

}

target "run-sim" {
  target = "prod-simc"
}
