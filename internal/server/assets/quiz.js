// toggles the color of the pixel at the give coordinates
function toggle(pixelid) {
    let pixel = document.getElementById(pixelid)
    if (pixel.classList.contains("dark")) {
        pixel.classList.remove("dark")
    } else {
        pixel.classList.add("dark")
    }
}

function batchToggle(coordinates) {

    coordinates.forEach(pixel => {
        toggle(`${pixel[0]}-${pixel[1]}`)
    });
}