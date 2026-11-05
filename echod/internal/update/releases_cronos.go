//go:build !dot && !spot && !panel

package update

// releases is where the Echo Show 5's releases are published.
// Integration builds must not accidentally follow upstream firmware updates
// and overwrite the ported Muse features. There are no published signed releases
// for this fork yet; OTA stays unavailable until an owner-controlled trust root
// and rollback plan are verified before device installation.
const releases = "https://github.com/dainbe/hey-muse/releases"
