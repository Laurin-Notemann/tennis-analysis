import { isLoggedIn, setAccessTokenAndUser } from "./auth.js"

export async function loadNavBar() {
  const navbar = document.querySelector("nav")
  const res = await isLoggedIn()
  const isAuthenticated = setAccessTokenAndUser(res)
  if (isAuthenticated) {
    const pUsername = document.createElement("p")
    const username = localStorage.getItem("username")
    pUsername.innerHTML = username
    navbar.appendChild(pUsername)

    const logoutButton = document.createElement("button")
    logoutButton.addEventListener("click", e => {
      e.preventDefault()
      localStorage.clear("access-token")
      localStorage.clear("userId")
      localStorage.clear("username")
      window.location.href = "/"
    })
    logoutButton.innerHTML = "Logout"
    navbar.appendChild(logoutButton)
  } else {
    const loginLink = document.createElement("button")
    loginLink.addEventListener("click", e => {
      e.preventDefault()
      window.location = "/login"
    })
    loginLink.innerHTML = "Login"

    const registerLink = document.createElement("button")
    registerLink.addEventListener("click", e => {
      e.preventDefault()
      window.location = "/register"
    })
    registerLink.innerHTML = "Register"

    navbar.appendChild(loginLink)
    navbar.appendChild(registerLink)
  }
}
