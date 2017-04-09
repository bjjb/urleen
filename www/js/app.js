if ('serviceWorker' in navigator) {
  navigator.serviceWorker.register('/sw.js').catch(console.error)
}

self.addEventListener('load', event => {
  var back    = document.getElementById('back')
  var copy    = document.getElementById('copy')
  var form    = document.getElementById('form')
  var input   = document.getElementById('input')
  var output  = document.getElementById('output')
  var result  = document.getElementById('result')
  var url     = document.getElementById('url')

  form.addEventListener('submit', event => {
    event.preventDefault()
    fetch(form.action, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(url.value)
    }).then(response => response.json()).then(id => {
      result.innerHTML = form.action + id
      input.hidden = true
      output.hidden = false
      document.getSelection().selectAllChildren(result)
    })
  })

  back.addEventListener('click', event => {
    output.hidden = true
    input.hidden = false
    copy.classList.remove('ok', 'error')
    url.value = ''
    url.focus()
  })

  copy.addEventListener('click', event => {
    if (!document.execCommand('copy')) {
      copy.classList.add('error')
    }
  })

  document.addEventListener('copy', event => {
    copy.classList.remove('error')
    copy.classList.add('ok')
  })
})
