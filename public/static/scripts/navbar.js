import { isLoggedIn, setAccessTokenAndUser } from "./auth.js"

export async function loadNavBar() {
  const navbar = document.querySelector("nav")
  const res = await isLoggedIn()
  const isAuthenticated = setAccessTokenAndUser(res)
  if (isAuthenticated) {
    const userInfo = document.createElement("div")
    userInfo.classList.add("user-information")
    const logoutButton = document.createElement("button")
    logoutButton.addEventListener("click", e => {
      e.preventDefault()
      localStorage.clear("access-token")
      localStorage.clear("userId")
      localStorage.clear("username")
      localStorage.clear("first-name")
      localStorage.clear("last-name")
      window.location.href = "/"
    })
    logoutButton.innerHTML = "Logout"
    logoutButton.classList.add("logout-button")
    const pImage = new Image();
    pImage.src = "/static/assets/user-4250.svg"
    pImage.classList.add("profile-icon")
    userInfo.appendChild(pImage)
    userInfo.appendChild(logoutButton)
    if (window.location.pathname === "/") {
      const emptydiv = document.createElement("div")
      emptydiv.innerHTML = ""
      navbar.appendChild(emptydiv)
      navbar.appendChild(userInfo)
    } else {
      console.log("test")
      const homeLink = document.createElement("h3")
      homeLink.innerHTML = "Tennis Analysis"
      homeLink.classList.add("nav-home-link")
      homeLink.addEventListener("click", e => {
        e.preventDefault()

        window.location.href = "/"
      })
      navbar.append(homeLink)
      navbar.appendChild(userInfo)
    }
  } else {
    if (window.location.pathname !== "/") {
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
}
