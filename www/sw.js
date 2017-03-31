var CACHE = 'urleen-v1'
var CACHED = [
  '/',
  'index.html',
  '/css/style.css',
  '/js/app.js'
]

self.addEventListener('install', event => {
  console.debug('install! %o', event)
  // event.waitUntil(caches.open(CACHE).then(cache => cache.addAll(CACHED)))
})

self.addEventListener('activate', event => {
  console.debug('activate! %o', event)
})

self.addEventListener('message', event => {
  console.debug('message! %o', event)
})

self.addEventListener('fetch', event => {
  console.debug('fetch! %o', event)
  // event.respondWith(caches.match(event.request).then(response => {
  //   return response ? response : fetch(event.request)
  // }))
})

self.addEventListener('sync', event => {
  console.debug('sync! %o', event)
})

self.addEventListener('push', event => {
  console.debug('push! %o', event)
})
