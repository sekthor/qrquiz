// the width & height of the qr without the quiet zone (padding)
const QR_SIZE = 25;
// the size of a "position" square (corner square)
const POSITION_SIZE = 7

// returns the coordinates of all black pixels of the qr default layout
function layout(size) {
   let coordinates = []; 
   coordinates.push(...position(0,0))          // top left
   coordinates.push(...position(QR_SIZE-7, 0)) // top right
   coordinates.push(...position(0,QR_SIZE-7))  // bottom right
   return coordinates
}

// returns coordinates of black pixels for "position" square (see qr spec)
// from a given starting coordinate
// 7x7 pixels; starting coordinate = top left corner
function position(x,y){
    let coordinates = []

    for (let i = 0; i < POSITION_SIZE; i++) {
        coordinates.push([x+i,y])
        coordinates.push([x+i, y+POSITION_SIZE-1])
    }

    for (let i = 1; i < POSITION_SIZE-1; i++) {
        coordinates.push([x,y+i])
        coordinates.push([x+POSITION_SIZE-1,y+i])
    }

    for (let i = 2; i < 5; i++) {
        coordinates.push([x+2, y+i])
        coordinates.push([x+3, y+i])
        coordinates.push([x+4, y+i])
    }

    return coordinates
}

// creates the QR code grid and fills it with the initial layout
function initQR() {
    let qrcontent = document.getElementById("qr")
    for (let i = 0; i < QR_SIZE; i++) {
        let next = "<tr class='row"+i+"'>"
        for (let j = 0; j < QR_SIZE; j++) {
            next += "<td class='col"+j+" pixel' id='"+j+"-"+i+"' onclick='handleChange(this.id)'></td>"
        }
        next += "</tr>"
        qrcontent.innerHTML += next
    }
    const defaultLayout = layout(QR_SIZE)
    defaultLayout.forEach(pixel => {
        toggle(...pixel)
    })
}

// toggles the color of the pixel at the give coordinates
function toggle(x,y) {
    let pixel = document.getElementById(x+"-"+y)
    if (pixel.classList.contains("dark")) {
        pixel.classList.remove("dark")
    } else {
        pixel.classList.add("dark")
    }
}

function initQuestions() {
    let questions = [
        {q: "Placehoder question?", answers: {"12-6": "answer1", "7-8": "answer2"}},
        {q: "Placeholder question?", answers: {"11-10": "answer1", "11-12": "answer2"}}
    ]

    let questionsContent = ""
    for (let i = 0; i < questions.length; i++) {
        let q = questions[i]
        questionsContent += "<div><h2>" + (i+1) + ") " + q.q + "</h2>";
        for (const [key, value] of Object.entries(q.answers)) {
            questionsContent += "<input type='checkbox' name='answer' value='"+ key +"' onchange=\"handleChange(this.value);\" /><label>" + value + " (" + key + ")</label><br />"
        }
        questionsContent += "</div>"
    }

    document.getElementById("questions").innerHTML = questionsContent
}

function handleChange(coordinateString) {
    let coordinates = coordinateString.split("-")
    toggle(coordinates[0], coordinates[1])
}

initQR();
initQuestions();