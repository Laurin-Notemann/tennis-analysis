import { loadNavBar } from "./navbar.js";
import { addMessage, displayErrorMessage, getHeaders, fetchAllPlayers } from "./utils.js";

loadNavBar()

const dropdownPOne = document.querySelector(`[data-dropdown="dropdown-player-one"]`)
const dropdownPTwo = document.querySelector(`[data-dropdown="dropdown-player-two"]`)
const playerOnInput = document.querySelector(`[data-input="player-one"]`)
playerOnInput.value = ""
const playerTwInput = document.querySelector(`[data-input="player-two"]`)
playerTwInput.value = ""

async function loadPlayerDropdown(dropdownEl, playerNumber) {
  const listOfPlayerButtons = []
  const allPlayers = await fetchAllPlayers()
  allPlayers.map(player => {
    const playerButton = document.createElement("button")
    playerButton.innerHTML = player.FirstName + " " + player.LastName
    playerButton.addEventListener("click", e => {
      e.preventDefault()

      if (playerNumber == 1) {
        const playerOneInput = document.querySelector(`[data-input="player-one"]`)
        playerOneInput.value = player.ID
        const playerOneButton = document.querySelector(`[data-button="choose-player-one"]`)
        playerOneButton.innerHTML = player.FirstName + " " + player.LastName
        dropdownPOne.style.visibility = "hidden"
        dropdownPOne.visible = false
      } else if (playerNumber == 2) {
        const playerTwoInput = document.querySelector(`[data-input="player-two"]`)
        playerTwoInput.value = player.ID
        const playerTwoButton = document.querySelector(`[data-button="choose-player-two"]`)
        playerTwoButton.innerHTML = player.FirstName + " " + player.LastName
        dropdownPTwo.style.visibility = "hidden"
        dropdownPTwo.visible = false
      }
    })

    listOfPlayerButtons.push(playerButton)
  })

  listOfPlayerButtons.map(button => {
    dropdownEl.append(button)
  })
}

loadPlayerDropdown(dropdownPOne, 1)
const playerOneButton = document.querySelector(`[data-button="choose-player-one"]`)
playerOneButton.addEventListener("click", e => {
  e.preventDefault()

  if (dropdownPOne.visible) {
    dropdownPOne.style.visibility = "hidden";
    dropdownPOne.visible = false
  } else {
    dropdownPOne.style.visibility = "visible";
    dropdownPOne.visible = true
  }
})

loadPlayerDropdown(dropdownPTwo, 2)
const playerTwoButton = document.querySelector(`[data-button="choose-player-two"]`)
playerTwoButton.addEventListener("click", e => {
  e.preventDefault()
  dropdownPTwo.visibility = "hidden";
  if (dropdownPTwo.visible) {
    dropdownPTwo.style.visibility = "hidden";
    dropdownPTwo.visible = false
  } else {
    dropdownPTwo.style.visibility = "visible";
    dropdownPTwo.visible = true
  }
})

const createTeamForm = document.querySelector(`[data-form="create-team-form"]`)
createTeamForm.addEventListener("submit", async e => {
  e.preventDefault()
  const message = document.querySelector(`[data-message="create-team-message"]`)
  message.innerHTML = ""

  const data = new URLSearchParams(new FormData(e.target))

  const playerOneId = data.get("player-one")
  const playerTwoId = data.get("player-two")
  let correctInput = true

  if (playerOneId == "" || playerTwoId == "") {
    const payload = {
      message: "no player"
    }
    displayErrorMessage(payload)
    correctInput = false
  } else if (playerOneId == playerTwoId) {
    const payload = {
      message: "team"
    }
    displayErrorMessage(payload)
    correctInput = false
  }

  if (correctInput) {
    let teamName = ""
    const nameInput = data.get("team-name")
    if (nameInput) {
      teamName = nameInput
    }

    const userId = localStorage.getItem("userId")
    const body = {
      PlayerOne: playerOneId,
      PlayerTwo: playerTwoId,
      Name: teamName,
      UserID: userId
    }
    const headers = getHeaders()
    const res = await fetch("/api/teams", {
      method: "POST",
      body: JSON.stringify(body),
      headers
    })

    if (res.status == 201) {
      addMessage(message, "Team was successfully created", "Teams", "/teams")
    } else {
      message.innerHTML = ""
      displayErrorMessage(await res.json())
    }

  }

})
