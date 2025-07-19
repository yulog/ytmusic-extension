// Establish connection with the native application
const port = chrome.runtime.connectNative('com.google.chrome.example.echo');

// Optional: Log messages received from the native app (if any)
port.onMessage.addListener((msg) => {
  console.log('Received message from native app:', msg);
});

// Log disconnection errors
port.onDisconnect.addListener(() => {
  if (chrome.runtime.lastError) {
    console.log('Disconnected from native app:', chrome.runtime.lastError.message);
  } else {
    console.log('Disconnected from native app');
  }
});

// Listen for messages from content scripts
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  // We only care about messages from our content script, not other extensions.
  // The message should be the song info object.
  if (sender.tab && sender.tab.url.includes("music.youtube.com")) {
    console.log('Received song info from content script:', message);
    
    // The native app expects an object with a "text" property.
    // We'll serialize the song info object into a JSON string.
    const messageToSend = {
      text: JSON.stringify(message) 
    };

    port.postMessage(messageToSend);
    console.log('Sent song info to native app.');
  }
  // Return true to indicate you wish to send a response asynchronously
  // (although we don't use it here, it's good practice).
  return true; 
});