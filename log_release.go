//go:build !debuglog

package main

// logSongInfo is a no-op function for release builds.
// This version is compiled when the 'debuglog' build tag is NOT specified.
func logSongInfo(songInfo SongInfo) {
	// Do nothing
}
