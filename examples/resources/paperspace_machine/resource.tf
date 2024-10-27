
# Manage example machine
resource "paperspace_machine" "example" {
  name         = "Example Name"
  machine_type = "C1"
  template_id  = "tkni3aa4"
  disk_size    = 50
  region       = "ny2"
}
