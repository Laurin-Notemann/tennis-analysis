import { loadNavBar } from "./navbar.js"

loadNavBar()

async function loadIndexPage() {

  const userId = localStorage.getItem("userId")
  const mainHtml = document.querySelector("main")
  if (userId) {

    //const navbar = document.querySelector("nav")

    const username = localStorage.getItem("username")
    mainHtml.appendChild(mainPageButton("Players", "/players"))
    mainHtml.appendChild(mainPageButton("Teams", "/teams"))
    mainHtml.appendChild(mainPageButton("Matches", "/matches"))

    const greetingEl = document.querySelector(".user-greeting")
    const usernameEl = document.createElement("h2")
    usernameEl.innerHTML = username

    greetingEl.appendChild(usernameEl)
  } else {
    const err = document.createElement("p")
    err.innerHTML = "Pls log in to Test"
    const mainHtml = document.querySelector("main")
    mainHtml.appendChild(mainPageButton("Sign in", "/login"))
    mainHtml.appendChild(mainPageButton("Sign up", "/register"))

    const greetingEl = document.querySelector(".user-greeting")
    const usernameEl = document.createElement("h2")
    usernameEl.innerHTML = "please Sign in or Sign up"

    greetingEl.appendChild(usernameEl)
  }
}

loadIndexPage()

function mainPageButton(resource, url) {
  const button = document.createElement("button")
  button.innerHTML = resource
  button.addEventListener("click", e => {
    e.preventDefault()

    window.location.href = url
  })
  return button
}

function testButon() {
  const testButton = document.createElement("button")
  htmlBody.appendChild(testButton)

  testButton.innerHTML = "Test"
  testButton.addEventListener("click", async e => {
    e.preventDefault()

    const token = localStorage.getItem("access-token")
    const res = await fetch("/api/players/" + userId, {
      headers: {
        Authorization: "Bearer " + token
      }
    })
    try {
      console.log(await res.json())
    } catch (e) {
      console.log(e)
    }
  })
}
