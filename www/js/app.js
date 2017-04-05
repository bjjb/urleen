self.addEventListener('load', event => {
  var main = document.querySelector("main")
  var input = main.querySelector("#input")
  var output = main.querySelector("#output")
  var result = output.querySelector("#result")
  var urlInput = input.querySelector('input[name="url"]')
  input.querySelector("form").addEventListener('submit', event => {
    event.preventDefault()
    var form = event.target
    fetch(form.action, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(form.url.value)
    }).then(response => response.json()).then(id => {
      result.href = result.innerHTML = form.action + id
      input.hidden = true
      output.hidden = false
      document.getSelection().selectAllChildren(result)
    })
  })
  output.querySelector('#return').addEventListener('click', event => {
    output.hidden = true
    input.hidden = false
    urlInput.value = ''
    urlInput.focus()
  })
})
