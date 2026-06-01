package uasurfer

import "strings"

// IsWebView reports whether ua is an embedded in-app/webview user agent.
//
// Pass the UserAgent returned by Parse to avoid reparsing. If parsed is nil,
// IsWebView parses ua itself. The classifier is intentionally separate from
// UserAgent's struct layout so adding this signal remains source-compatible for
// callers that use unkeyed UserAgent literals.
func IsWebView(ua string, parsed *UserAgent) bool {
	normalized := normalise(ua)
	if normalized == "" {
		return false
	}
	if parsed == nil {
		parsed = Parse(ua)
	}

	if hasInAppWebViewToken(normalized) {
		return true
	}
	if strings.Contains(normalized, "; wv)") {
		return true
	}
	if !isIOSDevice(normalized) {
		return false
	}

	// uasurfer's iPhone/iPad fallback classifies unknown iOS WebKit agents as
	// Safari. Real Mobile Safari carries both Version/ and Safari/ tokens; bare
	// WKWebViews generally do not.
	return parsed.Browser.Name == BrowserSafari && !hasMobileSafariPair(normalized)
}

func hasInAppWebViewToken(ua string) bool {
	for _, token := range []string{
		"fb_iab",
		"fban",
		"fbav",
		"fbios",
		"instagram",
		"musical_ly",
		"bytedance",
		"tiktok",
		"trill",
		"snapchat",
		"pinterest",
		"line/",
		" gsa/",
	} {
		if strings.Contains(ua, token) {
			return true
		}
	}
	return false
}

func isIOSDevice(ua string) bool {
	return strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") || strings.Contains(ua, "ipod")
}

func hasMobileSafariPair(ua string) bool {
	return strings.Contains(ua, "version/") && strings.Contains(ua, "safari/")
}
