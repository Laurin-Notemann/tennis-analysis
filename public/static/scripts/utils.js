export function displayErrorMessage(userPayload) {
  const htmlBody = document.querySelector("body")
  let errorMessage = document.querySelector("#register-error-message")
  if (!errorMessage) {
    errorMessage = document.createElement("p")
    errorMessage.id = "register-error-message"
    htmlBody.appendChild(errorMessage)
  }

  if (userPayload.message.includes("users_username_unique")) {
    errorMessage.innerHTML = "Username already exists"
  } else if (userPayload.message.includes("users_email_unique")) {
    errorMessage.innerHTML = "Email already exists"
  } else if (userPayload.message == "password") {
    errorMessage.innerHTML = "Wrong Password combination"
  } else if (userPayload.message.includes("sql: no rows in result set")) {
    errorMessage.innerHTML = "Wrong Username"
  } else if (userPayload.message.includes("hashedPassword")) {
    errorMessage.innerHTML = "Wrong Password"
  } else {
    errorMessage.innerHTML = "Error with request"
  }

}

export function getHeaders() {
  const token = localStorage.getItem("access-token")
  const headers = {
    Authorization: "Bearer " + token,
    "Content-Type": "application/json"
  }

  return headers
}

export function addMessage(messageEl, message) {
  const errorMessage = document.querySelector("#register-error-message")
  if (errorMessage) {
    errorMessage.innerHTML = ""
  }

  const successMessage = document.createElement("p")
  successMessage.innerHTML = message

  const buttonToPlayers = document.createElement("button")
  buttonToPlayers.innerHTML = "View all Players"
  buttonToPlayers.addEventListener("click", e => {
    e.preventDefault()

    window.location.href = "/players"
  })

  messageEl.append(successMessage)
  messageEl.append(buttonToPlayers)
}
