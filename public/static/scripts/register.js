const registerUserForm = document.querySelector(`[data-form="register-user-form"]`);

registerUserForm.addEventListener("submit", async e => {
  e.preventDefault();

  const data = new URLSearchParams(new FormData(e.target))

  if (data.get("password") != data.get("confirm")) {
    console.log("do something")
  }

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
  console.log(userPayload)

  localStorage.setItem("access-token", userPayload.accessToken)

})

function loadNavBar() {
  const navBar = document.querySelector("nav")

}

loadNavBar()

