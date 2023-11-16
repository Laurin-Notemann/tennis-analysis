import { addMessage, displayErrorMessage, getHeaders } from "./utils.js";
import { loadNavBar } from "./navbar.js";

const createPlayer = document.querySelector(`[data-form="create-player-form"]`);

createPlayer.addEventListener("submit", async e => {
  e.preventDefault();
  const message = document.querySelector(`[data-message="create-player-message"]`)
  message.innerHTML = ""

  setTimeout(() => {
  }, 100)

  const data = new URLSearchParams(new FormData(e.target))

  const body = {
    firstName: data.get("first-name"),
    lastName: data.get("last-name"),
    userId: localStorage.getItem("userId")
  }

  const headers = getHeaders()

  const res = await fetch("/api/players", {
    method: "POST",
    body: JSON.stringify(body),
    headers
  })

  if (res.status == 201) {
    addMessage(message, "Player was successfully created", "Players", "/players")
  }
  else {
    message.innerHTML = ""
    displayErrorMessage(await res.json())
  }
})

loadNavBar()
