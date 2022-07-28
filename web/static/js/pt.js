function formatCreatedDate(d) {
    let date = new Date(d)
    return date.toLocaleString('en-gb', { 
        weekday:"long", 
        month:"long", 
        year:"numeric", 
        day:"numeric", 
        hour:"numeric", 
        minute:"numeric"
    })
}