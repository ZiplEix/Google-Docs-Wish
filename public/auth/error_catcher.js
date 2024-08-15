document.body.addEventListener('htmx:beforeSwap', function(evt) {
    if(evt.detail.xhr.status == 401 || evt.detail.xhr.status == 500) {
        evt.detail.shouldSwap = true;
        evt.detail.isError = false;
        console.log(evt.detail.xhr.responseText);
    }
});
