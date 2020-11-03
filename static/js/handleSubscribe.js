async function Subscribe(event) {
    event.preventDefault()
    const data = new FormData(event.target)
    let object = {}
    data.forEach(function(value, key){
        object[key] = value;
    });
    const req = await fetch("/v1/api/sign-up", {
        headers: {
            'Content-Type': 'application/json'
        },
        method: "post",
        body: JSON.stringify(object),
    })

    if (req.status === 200) {
        // TODO: will only redirect back to font page now
        return window.location.replace("/subscribe/success")
    }
    if (req.status === 400) {
        return window.location.replace("/subscribe/failure")
    }
}