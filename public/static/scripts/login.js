import { setAccessTokenAndUser } from "./auth.js";
import { loadNavBar } from "./navbar.js";
import { displayErrorMessage } from "./utils.js";

const loginUserForm = document.querySelector(`[data-form="login-user-form"]`);

loginUserForm.addEventListener("submit", async e => {
  e.preventDefault();

  const data = new URLSearchParams(new FormData(e.target))

  const body = {
    usernameOrEmail: data.get("username"),
    password: data.get("password"),
  }

  const res = await fetch("/api/login", {
    method: "POST",
    body: JSON.stringify(body),
    headers: {
      "Content-Type": "application/json"
    }
  })
  if (res.status == 200) {
    const userPayload = await res.json()
    setAccessTokenAndUser({ success: true, payload: userPayload })
    window.location.href = "/"
  }
  else {
    displayErrorMessage(await res.json())
  }
})

loadNavBar()

