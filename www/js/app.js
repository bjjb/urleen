// register service worker, if supported
if ('serviceWorker' in navigator) {
  self.addEventListener('load', event => {
    navigator.serviceWorker.register('/sw.js').then(registration => {
      console.log("Service worker registered; scope: %s", registration.scope)
      console.debug(registration)
    }).catch(err => {
      console.error("Service worker registration failed: %s", err)
    })
  })
}

// listen for paste events
self.addEventListener('paste', event => console.log(event))
