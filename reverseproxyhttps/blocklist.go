package main

// blocklist returns a map of bad IP addresses which are to be blocked.
func blocklist() map[string]struct{} {
	badIPs := make(map[string]struct{})
	// Each bad IP is registered to the blocklist map with an empty struct
	// as this requires zero memory allocation. If this was map[string]bool
	// then the bool would have to be allocated memory.
	badIPs["127.0.0.2"] = struct{}{}
	badIPs["127.0.0.3"] = struct{}{}
	badIPs["127.0.0.4"] = struct{}{}
	return badIPs
}
