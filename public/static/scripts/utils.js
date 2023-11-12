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
  }

}
