provider "openfaas" {
  version = "~> 0.2.0"
  uri       = "http://localhost:8080"
  user_name = "admin"
  password  = var.openfaas_provider_password
}

provider "postgresql" {
  host      = "localhost"
  username  =  var.database_username
  password  =  var.database_password
  sslmode   = "disable"
}
