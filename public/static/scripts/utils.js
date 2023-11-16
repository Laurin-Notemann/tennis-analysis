export function displayErrorMessage(userPayload) {
  const htmlBody = document.querySelector(".form-wrapper")
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
  } else if (userPayload.message.includes("team")) {
    errorMessage.innerHTML = "Please enter two different Players"
  } else if (userPayload.message.includes("no player")) {
    errorMessage.innerHTML = "Please enter two Players"
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

export function addMessage(messageEl, message, resource, url) {
  const errorMessage = document.querySelector("#register-error-message")
  if (errorMessage) {
    errorMessage.innerHTML = ""
  }

  const successMessage = document.createElement("p")
  successMessage.innerHTML = message
  successMessage.classList.add("success-message")

  const buttonToPlayers = document.createElement("button")
  buttonToPlayers.innerHTML = "View all " + resource
  buttonToPlayers.addEventListener("click", e => {
    e.preventDefault()

    window.location.href = url
  })

  messageEl.append(successMessage)
  messageEl.append(buttonToPlayers)
}

export async function fetchAllPlayers() {
  const headers = getHeaders()
  const userId = localStorage.getItem("userId")
  const res = await fetch("/api/players/" + userId, {
    headers: {
      Authorization: headers.Authorization
    }
  })
  if (res.status == 200) {
    const allPlayers = await res.json()
    return allPlayers
  }
  return []
}

