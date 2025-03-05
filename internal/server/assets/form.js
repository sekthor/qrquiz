
function addQuestion() {
    let question = document.getElementById("question").value
    let answers = JSON.parse(localStorage.getItem("answers"));
    let questions = JSON.parse(localStorage.getItem("questions"));

    if (questions === null) {
        questions = []
    }

    questions.push({question: question, answers: answers})
    
    localStorage.setItem("questions", JSON.stringify(questions))
    localStorage.setItem("answers", JSON.stringify([]))
    clearQuestion();
    displayAnswers();
}

function addAnswer() {
    let answer = document.getElementById("answer").value;
    let correct = document.getElementById("correct").value;
    let answers = JSON.parse(localStorage.getItem("answers"));

    if (answers === null) {
        answers = []
    }

    answers.push({"text": answer, "correct": (correct == "true" ? true : false)})
    localStorage.setItem("answers", JSON.stringify(answers))
    displayAnswers();
    clearAnswer();
}

function clearAnswer() {
    document.getElementById("answer").value = null;
}

function clearQuestion() {
    document.getElementById("question").value = null;
}

function displayAnswers() {
    let answers = JSON.parse(localStorage.getItem("answers"));
    if (answers === null) {
        return
    }
    
    let content = [];
    answers.forEach(answer => {
        content += `<tr>
            <td>${answer.text}</td>
            <td><i>${answer.correct ? "correct" : "wrong"}</i></td>
            <td><button type='button' onclick='removeAnswer("${answer.text}")'>Remove</button></td>
        </tr>`
    });
        
    document.getElementById("savedanswers").innerHTML = content
}

function removeAnswer(answertext) {
    let answers = JSON.parse(localStorage.getItem("answers"));
    if (answers === null) {
        answers = []
    }

    answers = answers.filter(answer => answer.text !== answertext)
    localStorage.setItem("answers", JSON.stringify(answers))
    displayAnswers();
}

function addQuiznameToTitle() {

}

function saveQuiz() {
    localStorage.setItem("title", document.getElementById("title").value);
    localStorage.setItem("secret", document.getElementById("secret").value);
    localStorage.setItem("questions", JSON.stringify([]))
    localStorage.setItem("answers", JSON.stringify([]))
    window.location.href = "/new/question";
}

function review() {
    window.location.href = "/new/review";
}

function fillReviewTemplate() {
    document.getElementById("title").innerText = localStorage.getItem("title")
    document.getElementById("secret").innerText = localStorage.getItem("secret")

    let questions = JSON.parse(localStorage.getItem("questions"));

    if (questions === null) {
        questions = []
    }

    let content = ""
    questions.forEach(question => {

        let answers = question.answers.map(answer => 
            `<li>${answer.text} <i>(${answer.correct ? "correct" : "wrong"})</i></li>`
        ).join("")

        content += `<div>
        <p><strong>${question.question}</strong></p>
        <ol>
            ${answers}
        </ol>
        </div>`
    })

    document.getElementById("answers").innerHTML = content
}

async function submit() {

    const body = {
            title: localStorage.getItem("title"),
            secret: localStorage.getItem("secret"),
            questions: JSON.parse(localStorage.getItem("questions"))
        }

    const headers = new Headers();
    headers.append("Content-Type", "application/json")

    const response = await fetch("/new", {
        method: "POST",
        body: JSON.stringify(body),
        headers: headers,
    })

    if (response.ok) {
        let quiz = await response.json()
        window.location.href = `/quiz/${quiz.id}`
    }
}