var CACHE = 'urleen-v1'
var URLS = [
  '/',
  '/css/style.css',
  '/js/app.js'
]

self.addEventListener('install', event => {
  event.waitUntil(caches.open(CACHE).then(cache => cache.addAll(URLS)))
})

self.addEventListener('fetch', event => {
  event.respondWith(caches.match(event.request).then(response => {
    return response ? response : fetch(event.request)
  }))
})
