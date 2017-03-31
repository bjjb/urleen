// register service worker, if supported
if ('serviceWorker' in navigator) {
  self.addEventListener('load', event => {
    navigator.serviceWorker.register('/sw.js').then(registration => {
      console.log("Service worker registered; scope: %s", registration.scope)
    }).catch(err => {
      console.error("Service worker registration failed: %s", err)
    })
  })
}

self.addEventListener('load', event => {
  document.querySelector("main form").addEventListener('submit', event => {
    event.preventDefault()
    if (self.fetch) {
      form = new FormData(event.target)
      fetch(event.target.action, {
        method: 'POST',
        headers: { 'Content-Type': event.target.enctype },
        body: "url=" + encodeURI(event.target['url'].value)
      }).then(response => {
        response = response.clone()
        console.debug(response.clone())
        response.text().then(text => console.debug(response.url + text))
      })
    }
  })
  console.debug("Page loaded")
})
