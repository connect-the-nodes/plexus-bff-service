resource "aws_ecr_repository" "main" {
  name = "zynchub-repo-${var.environment}"
}

resource "aws_ecr_lifecycle_policy" "cleanup" {
  repository = aws_ecr_repository.main.name
  policy = jsonencode({
    rules = [{
      rulePriority = 1
      selection = {
        tagStatus = "any"
        countType = "imageCountMoreThan"
        countNumber = 2
      }
      action = { type = "expire" }
    }]
  })
}
