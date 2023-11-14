import { loadNavBar } from "./navbar.js";

loadNavBar()

async function doSomething() {

  const userId = localStorage.getItem("userId")
  const htmlBody = document.querySelector("body")
  if (userId) {
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

      console.log(res)
      console.log(res.body)
      try {
        console.log(await res.json())
      } catch (e) {
        console.log(e)
      }
    })
  } else {
    const err = document.createElement("p")
    err.innerHTML = "Pls log in to Test"
    htmlBody.appendChild(err)
  }
}

doSomething()
