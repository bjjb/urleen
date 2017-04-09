var CACHE = 'urleen-v1'
var CACHED = [
  '/',
  '/index.html',
  '/css/style.css',
  '/css/urleen.woff',
  '/js/app.js'
]

self.addEventListener('install', event => {
  event.waitUntil(caches.open(CACHE).then(cache => cache.addAll(CACHED)))
})

self.addEventListener('activate', event => {
  event.waitUntil(caches.keys().then(keys => {
    return Promise.all(keys.map(key => {
      if (key !== CACHE) {
        return caches.delete(key)
      }
    }))
  }))
})

self.addEventListener('fetch', event => {
  if (event.request.method !== 'GET') {
    return fetch(event.request);
  }
  event.respondWith(caches.match(event.request).then(response => {
    return response || fetch(event.request).then(response => {
      return caches.open(CACHE).then(cache => {
        cache.put(event.request, response.clone())
        return response
      })
    })
  }))
})
