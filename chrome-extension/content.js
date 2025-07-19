console.log("YT Music Notifier: Content script injected. Starting polling loop.");

// --- State Tracking ---
// Track the last song title we successfully sent to avoid sending duplicates.
let lastKnownTitle = '';

// --- Data Extraction Methods (no changes from before) ---

function getSongInfoFromMediaSession() {
    const metadata = navigator.mediaSession?.metadata;
    if (!metadata || !metadata.title) return null;
    return {
        source: "MediaSessionAPI",
        title: metadata.title,
        artist: metadata.artist || '',
        album: metadata.album || '',
        artworkUrl: (() => {
            if (!metadata.artwork || metadata.artwork.length === 0) return '';
            let bestArtwork = metadata.artwork[0];
            let maxArea = 0;

            for (const artwork of metadata.artwork) {
                if (artwork.sizes) {
                    const parts = artwork.sizes.split('x');
                    if (parts.length === 2) {
                        const width = parseInt(parts[0], 10);
                        const height = parseInt(parts[1], 10);
                        const area = width * height;
                        if (area > maxArea) {
                            maxArea = area;
                            bestArtwork = artwork;
                        }
                    }
                }
            }
            return bestArtwork.src;
        })(),
    };
}

function getSongInfoFromDOM() {
    const titleEl = document.querySelector('ytmusic-player-bar .title');
    const bylineEl = document.querySelector('ytmusic-player-bar .byline');
    const artworkEl = document.querySelector('ytmusic-player-bar img.image');
    if (!titleEl || !bylineEl || !artworkEl) return null;
    return {
        source: "DOM",
        title: titleEl.textContent.trim(),
        byline: bylineEl.textContent.trim(),
        artworkUrl: artworkEl.src,
    };
}

// --- Main Polling Loop ---

setInterval(() => {
    // Pre-condition: The player bar must exist on the page.
    const playerBar = document.querySelector('ytmusic-player-bar');
    if (!playerBar) {
        // If player bar is not found, do nothing and wait for the next interval.
        return;
    }

    // 1. Get current song info (Media Session preferred)
    let songInfo = getSongInfoFromMediaSession();
    if (!songInfo || !songInfo.title) {
        songInfo = getSongInfoFromDOM();
    }

    // 2. Check if the song is new and valid
    if (songInfo && songInfo.title && songInfo.title !== lastKnownTitle) {
        console.log(`New song detected: "${songInfo.title}". Previous was: "${lastKnownTitle}"`);
        
        // 3. Send the new info to the background script
        chrome.runtime.sendMessage(songInfo);
        
        // 4. Update the last known title to prevent re-sending
        lastKnownTitle = songInfo.title;
    }
}, 2000); // Check every 2 seconds