resource "openfaas_function" "function_test" {
  name            = "test-function"
  image           = "functions/alpine:latest"
  f_process       = "env"
  labels = {
    Group       = "London"
    Environment = "Test"
  }

  limits {
    memory = "20m"
    cpu    = "100m"
  }

  env_vars = {
    database_name = "${postgresql_database.function_db.name}"
  }

  annotations =  {
    CreatedDate = "Mon 24 Feb 21:32:02 GMT 2020"
  }
}