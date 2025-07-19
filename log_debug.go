//go:build debuglog

package main

import "log"

// logSongInfo prints detailed information about the received song data.
// This version is compiled only when the 'debuglog' build tag is specified.
func logSongInfo(songInfo SongInfo) {
	log.Printf("Received song data (Source: %s)", songInfo.Source)
	if songInfo.Source == "MediaSessionAPI" {
		log.Printf("  Title: %s, Artist: %s, Album: %s", songInfo.Title, songInfo.Artist, songInfo.Album)
	} else {
		log.Printf("  Title: %s, Byline: %s", songInfo.Title, songInfo.Byline)
	}
	log.Printf("  Artwork: %s", songInfo.ArtworkURL)
}
