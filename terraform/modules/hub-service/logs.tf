######################################
# LOG GROUPS
######################################

# 1️⃣ ECS Log Group
resource "aws_cloudwatch_log_group" "ecs" {
  name              = "/ecs/zynchub-${var.environment}"
  retention_in_days = 14
  tags              = var.tags
}
