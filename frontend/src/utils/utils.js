const setCookie = (value, expiration) => {
    const date = new Date()
    date.setTime(date.getTime() + expiration * 24 * 60 * 60 * 1000)
    document.cookie = `session=${value};expires=${date.toUTCString()};path=/`
}

const getCookie = () => {
    const cookies = document.cookie.split(' ')
    for (let i = 0; i < cookies.length; i++) {
        const cookieParts = cookies[i].split('=')
        const name = cookieParts[0]
        const value = cookieParts[1]

        if (name === "session") {
            return decodeURIComponent(value)
        }
    }
    return null
}

const checkCookie = () => {
    return document.cookie.includes("session");
}

const formatDate = date => {
    const originalDate = new Date(date);

    const year = originalDate.getFullYear();
    const month = originalDate.getMonth() + 1;
    const day = originalDate.getDate();

    const hours = originalDate.getHours();
    const minutes = originalDate.getMinutes();

    const dateFormatee = `${year}/${month < 10 ? '0' : ''}${month}/${day < 10 ? '0' : ''}${day} - ${hours}h${minutes < 10 ? '0' : ''}${minutes}`;

    return dateFormatee
}




export {
    setCookie,
    getCookie,
    checkCookie,
    formatDate
}