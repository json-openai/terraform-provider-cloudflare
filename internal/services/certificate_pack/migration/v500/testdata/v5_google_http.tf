resource "cloudflare_certificate_pack" "%s" {
  zone_id               = "%s"
  type                  = "advanced"
  hosts                 = ["%s.com"]
  validation_method     = "http"
  validity_days         = 365
  certificate_authority = "google"
}
