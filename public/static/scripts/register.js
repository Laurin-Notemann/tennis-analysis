import { setAccessTokenAndUser } from "./auth.js";
import { loadNavBar } from "./navbar.js";
import { displayErrorMessage } from "./utils.js";

const registerUserForm = document.querySelector(`[data-form="register-user-form"]`);

registerUserForm.addEventListener("submit", async e => {
  e.preventDefault();

  const data = new URLSearchParams(new FormData(e.target))

  if (data.get("password") != data.get("confirm")) {
    displayErrorMessage({ message: "password" })
  } else {

    const body = {
      username: data.get("username"),
      email: data.get("email"),
      password: data.get("password"),
      confirm: data.get("confirm")
    }

    const res = await fetch("/api/register", {
      method: "POST",
      body: JSON.stringify(body),
      headers: {
        "Content-Type": "application/json"
      }
    })
    const userPayload = await res.json()
    if (res.status == 201) {
      setAccessTokenAndUser({ success: true, payload: userPayload })
      console.log(localStorage.getItem("user"))
      window.location.href = "/"
    } else if (res.status == 409) {
      displayErrorMessage(userPayload)
    }
  }
})

loadNavBar()

