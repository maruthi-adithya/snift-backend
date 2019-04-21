package utils

// Holds the list of Badges and Messages
const (
	HTTPSBadge               = "HTTP_SECURE"
	HTTPSBadgeMessage        = "This website is encrypted and hence protected from Man-in-the-Middle attacks(MITM) and Eavesdropping Attacks."
	XSSBadge                 = "XSS_PROTECT"
	XSSBadgeMessage          = "This site is prevented from reflected cross-site scripting (XSS) attacks"
	XFrameBadge              = "CLICKJACKING_PROTECT"
	XFrameBadgeMessage       = "The content from this site cannot be embedded into other sites and is protected cross-site from Clickjacking"
	HSTSBadge                = "HTTPS_ONLY"
	HSTSBadgeMessage         = "This site can only be accessed via HTTPS"
	CSPBadge                 = "CSP_ENABLED"
	CSPBadgeMessage          = "This site has an additional security layer to proctect from Cross Site Scripting (XSS), data injection and packet sniffing attacks"
	HPKPBadge                = "PUBLIC_KEY_PINNING_ENABLED"
	HPKPBadgeMessage         = "This site has a security feature that decreases the risk of Man-in-the-Middle attacks(MITM) attacks with forged certificates."
	RPBadge                  = "ENSURE_PRIVACY"
	RPBadgeMessage           = "This site has a Referrer Policy that protects user's privacy"
	XContentTypeBadge        = "NO_SNIFF"
	XContentTypeBadgeMessage = "This site prevents the browser from doing MIME-type sniffing"
	HTTPVersionBadge         = "LATEST_HTTP"
	HTTPVersionBadgeMessage  = "This site uses the latest HyperText Transfer Protocol(HTTP) supporting better performance and security standards"
	TLSVersionBadge          = "LATEST_TLS"
	TLSVersionBadgeMessage   = "This site uses the latest Transport Layer Security(TLS) supporting better performance and security standards"
	SPFBadge                 = "EMAIL_SPOOFING_PROTECT"
	SPFBadgeMessage          = "This site has a valid Sender Policy Framework(SPF) record to improve email deliverability and protects this domain against malicious emails sent on behalf of this domain"
)
