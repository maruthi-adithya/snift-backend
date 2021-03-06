package utils

// Holds the list of Badges and Messages
const (
	HTTPSBadge                   = "HTTP_SECURE"
	HTTPSBadgeMessage            = "Encrypted HTTPS Connection"
	HTTPSBadgeDescription        = "This site is encrypted and is less prone to Man-in-the-Middle attacks(MITM) and Eavesdropping Attacks"
	XSSBadge                     = "XSS_PROTECT"
	XSSBadgeMessage              = "Prevention from reflected Cross-Site Scripting (XSS) Attacks"
	XSSBadgeDescription          = "This site is less prone to from reflected cross-site scripting (XSS) attacks"
	XFrameBadge                  = "CLICKJACKING_PROTECT"
	XFrameBadgeMessage           = "Protection from Cross-Site Click Jacking Attacks"
	XFrameBadgeDescription       = "The content from this site cannot be embedded into other sites and is protected from cross-site Clickjacking"
	HSTSBadge                    = "HTTPS_ONLY"
	HSTSBadgeMessage             = "Enforces HTTPS-Only Site Access"
	HSTSBadgeDescription         = "This site can only be accessed via HTTPS"
	CSPBadge                     = "CSP_ENABLED"
	CSPBadgeMessage              = "Protection against Cross Site Scripting (XSS), Data Injection and Packet Sniffing attacks"
	CSPBadgeDescription          = "This site has is relatively secure against Cross Site Scripting (XSS), Data Injection and Packet Sniffing attacks"
	HPKPBadge                    = "PUBLIC_KEY_PINNING_ENABLED"
	HPKPBadgeMessage             = "Prevention against Man-in-the-Middle attacks(MITM) using forged certificates"
	HPKPBadgeDescription         = "This site has a decreased risk of Man-in-the-Middle attacks(MITM) with forged certificates"
	RPBadge                      = "ENSURE_PRIVACY"
	RPBadgeMessage               = "Enforces a Referrer Policy to avoid leaking sensitive user information from being shared."
	RPBadgeDescription           = "This site has a Referrer Policy that may help protect user privacy"
	XContentTypeBadge            = "NO_SNIFF"
	XContentTypeBadgeMessage     = "Prevention from media-type (MIME) sniffing"
	XContentTypeBadgeDescription = "This site prevents the browser from media type (MIME) sniffing"
	HTTPVersionBadge             = "LATEST_HTTP"
	HTTPVersionBadgeMessage      = "Uses the latest version of the HTTP Protocol"
	HTTPVersionBadgeDescription  = "This site uses the latest HyperText Transfer Protocol(HTTP) supporting better performance and security standards"
	TLSVersionBadge              = "LATEST_TLS"
	TLSVersionBadgeMessage       = "Uses the latest version of the TLS Protocol"
	TLSVersionBadgeDescription   = "This site uses the latest Transport Layer Security(TLS) supporting better performance and security standards"
	SPFBadge                     = "EMAIL_SPOOFING_PROTECT"
	SPFBadgeMessage              = "Prevention from Email Spoofing by having a valid Sender Policy Framework Record"
	SPFBadgeDescription          = "This site has a valid Sender Policy Framework(SPF) record that reduces the risk of forged emails being sent on behalf of this domain"
)
