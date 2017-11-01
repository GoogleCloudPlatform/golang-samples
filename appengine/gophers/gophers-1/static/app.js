'use strict';
document.addEventListener('DOMContentLoaded', function() {

  let form = document.querySelector('form');
  form.addEventListener('submit', async function(e) {
    e.preventDefault();

    let resp = await fetch(form.action, {
      credentials: 'include',
      method: form.method,
      body: new FormData(form),
    })
    let body = await resp.json();
    if (body.Message) {
      showNotice(body.Message);
      if (resp.ok) {
        document.querySelector('input[name=message]').value = '';
      }
    }
  });

  let notice = document.querySelector('#notice');
  function showNotice(text) {
    notice.innerText = text;
    notice.display = '';
    setTimeout(function() {
      notice.display = 'none';
    }, 2000);
  }
});
