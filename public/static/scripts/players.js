import { loadNavBar } from "./navbar.js";
import { getHeaders } from "./utils.js";

loadNavBar()

async function renderAllPlayer() {
  const htmlBody = document.querySelector("body")
  const players = document.querySelector(`[data-all-players="all-players-div"]`)
  players.innerHTML = ""

  const createPlayer = document.createElement("button")
  createPlayer.classList.add("create-player-team-button")
  createPlayer.innerHTML = "Create New Player"
  createPlayer.addEventListener("click", e => {
    e.preventDefault()

    window.location.href = "/create-player"
  })
  players.append(createPlayer)
  const headers = getHeaders()
  const userId = localStorage.getItem("userId")
  const res = await fetch("/api/players/" + userId, {
    headers: {
      Authorization: headers.Authorization
    }
  })
  if (res.status == 500) {
    const errorEl = document.createElement("p")
    errorEl.innerHTML = "Couldn't fetch players"
  } else if (res.status == 200) {
    const playersObj = await res.json()
    if (playersObj == null) {
      const noPlayerMessage = document.createElement("p")
      noPlayerMessage.innerHTML = "No players created yet"
      htmlBody.append(noPlayerMessage)
    } else {
      playersObj.map(player => {
        const playerId = player.ID
        const playerEl = document.createElement("div")
        playerEl.classList.add("player-team-obj")

        const playerEditButton = document.createElement("button")
        playerEditButton.innerHTML = "Edit"
        playerEditButton.addEventListener("click", e => {
          e.preventDefault()
          localStorage.setItem("player-first-name", player.FirstName)
          localStorage.setItem("player-last-name", player.LastName)
          window.location.href = "/edit-player/" + playerId
        })

        const playerDeleteButton = document.createElement("button")
        playerDeleteButton.innerHTML = "x"
        playerDeleteButton.addEventListener("click", async e => {
          e.preventDefault()

          const headers = getHeaders()
          const res = await fetch("/api/players/" + playerId, {
            method: "DELETE",
            headers: headers
          })
          if (res.status == 200) {
            await renderAllPlayer()
          } else if (res.status == 500) {

          }
        })

        const playerName = document.createElement("p")
        playerName.innerHTML = player.FirstName + " " + player.LastName

        playerEl.appendChild(playerName)
        const editDelete = document.createElement("div")
        editDelete.classList.add("edit-delete-button")
        editDelete.appendChild(playerEditButton)
        editDelete.appendChild(playerDeleteButton)


        playerEl.appendChild(editDelete)
        players.appendChild(playerEl)
      })
    }
  }
}


renderAllPlayer()

