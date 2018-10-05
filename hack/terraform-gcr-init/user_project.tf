resource "google_project" "project" {
  name            = "cluster-api-${var.short_name}"
  project_id      = "cluster-provider-aws-${var.short_name}"
  org_id          = "${var.org_id}"
  billing_account = "${var.billing_account}"
}

resource "google_project_services" "project" {
  project  = "${google_project.project.project_id}"
  services = ["containerregistry.googleapis.com", "pubsub.googleapis.com", "storage-api.googleapis.com"]
}

resource "google_project_iam_policy" "project" {
  project     = "${google_project.project.project_id}"
  policy_data = "${data.google_iam_policy.owner.policy_data}"
}

data "google_iam_policy" "owner" {
  binding {
    role = "roles/owner"

    members = [
      "${var.owners}",
    ]
  }
}

resource "null_resource" "gcr" {
  // We have to push at least one image for the repository to be created.
  provisioner "local-exec" {
    command = "sh -c 'docker pull busybox && docker tag busybox gcr.io/${google_project.project.project_id}/busybox:init && docker push gcr.io/${google_project.project.project_id}/busybox:init'"
  }
}

// This makes the bucket publicly readable
resource "google_storage_bucket_acl" "acl" {
  bucket = "artifacts.${google_project.project.project_id}.appspot.com"

  role_entity = [
    "OWNER:${var.owners}",
    "READER:allUsers",
  ]

  provisioner "local-exec" {
    command = "gsutil -m acl ch -r -u AllUsers:READ gs://artifacts.${google_project.project.project_id}.appspot.com"
  }

  depends_on = ["null_resource.gcr"]
}
