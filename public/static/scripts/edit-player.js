import { loadNavBar } from "./navbar.js";
import { getHeaders, addMessage, displayErrorMessage } from "./utils.js";

loadNavBar()

const editPlayerForm = document.querySelector(`[data-form="edit-player-form"]`)
const firstNameInput = document.querySelector(`[data-input="edit-first-name"]`)
firstNameInput.value = localStorage.getItem("player-first-name")
const lastNameInput = document.querySelector(`[data-input="edit-last-name"]`)
lastNameInput.value = localStorage.getItem("player-last-name")

editPlayerForm.addEventListener("submit", async e => {
  e.preventDefault()

  const message = document.querySelector(`[data-message="edit-player-message"]`)
  message.innerHTML = ""
  setTimeout(() => {
  }, 100)

  const currentPath = window.location.pathname
  const splitPath = currentPath.split("/")
  const playerId = splitPath[2]

  const headers = getHeaders()

  const data = new URLSearchParams(new FormData(e.target))
  const body = {
    FirstName: data.get("first-name"),
    LastName: data.get("last-name"),
    ID: playerId
  }

  const res = await fetch("/api/players", {
    method: "PUT",
    body: JSON.stringify(body),
    headers
  })

  if (res.status == 200) {
    addMessage(message, "Player was successfully updated", "Players", "/players")
    localStorage.setItem("player-first-name", data.get("first-name"))
    localStorage.setItem("player-last-name", data.get("last-name"))
  } else {
    message.innerHTML = ""
    displayErrorMessage(await res.json())
  }
})
