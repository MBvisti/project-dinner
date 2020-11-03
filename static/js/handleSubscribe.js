async function Subscribe(event) {
    event.preventDefault()
    const data = new FormData(event.target)

    let object = {}
    data.forEach(function(value, key){
        object[key] = value;
    });

    try {
        const req = await fetch("/v1/api/sign-up", {
            headers: {
                'Content-Type': 'application/json'
            },
            method: "post",
            body: JSON.stringify(object),
        })

        console.log(req)
    } catch (e) {
        console.log(e)
    } finally {
        return "hello"
    }
}