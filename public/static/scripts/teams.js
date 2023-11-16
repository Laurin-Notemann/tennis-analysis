import { loadNavBar } from "./navbar.js";
import { getHeaders, fetchAllPlayers } from "./utils.js";

loadNavBar()

async function showAllTeams() {
  const htmlBody = document.querySelector("body")
  const teams = document.querySelector(`[data-all-teams="all-teams-div"]`)
  teams.innerHTML = ""

  const createTeam = document.createElement("button")
  createTeam.classList.add("create-player-team-button")
  createTeam.innerHTML = "Create New Team"
  createTeam.addEventListener("click", e => {
    e.preventDefault()

    window.location.href = "/create-team"
  })
  teams.append(createTeam)

  const headers = getHeaders()
  const userId = localStorage.getItem("userId")
  const res = await fetch("/api/teams/" + userId, {
    headers: {
      Authorization: headers.Authorization
    }
  })
  if (res.status == 500) {
    const errorEl = document.createElement("p")
    errorEl.innerHTML = "Couldn't fetch teams"
  } else if (res.status == 200) {
    const teamsOb = await res.json()
    const teamsObj = teamsOb.filter(team => team.PlayerTwo != null)
    if (teamsObj == null || teamsObj.length === 0) {
      const noTeamMessage = document.createElement("p")
      noTeamMessage.innerHTML = "No teams created yet"
      htmlBody.append(noTeamMessage)
    } else {
      teamsObj.map(async team => {
        const teamId = team.ID
        const teamEl = document.createElement("div")
        teamEl.classList.add("player-team-obj")

        const teamEditButton = document.createElement("button")
        teamEditButton.innerHTML = "Edit"
        teamEditButton.addEventListener("click", e => {
          e.preventDefault()
          localStorage.setItem("team-player-one", team.PlayerOne.ID)
          localStorage.setItem("team-player-two", team.PlayerTwo.ID)
          window.location.href = "/edit-team/" + teamId
        })

        const teamDeleteButton = document.createElement("button")
        teamDeleteButton.innerHTML = "x"
        teamDeleteButton.addEventListener("click", async e => {
          e.preventDefault()

          const headers = getHeaders()
          const res = await fetch("/api/teams/" + teamId, {
            method: "DELETE",
            headers: headers
          })
          if (res.status == 200) {
            await showAllTeams()
          } else if (res.status == 500) {

          }
        })

        const players = await fetchAllPlayers()

        if (team.PlayerTwo) {
          const playerOne = players.find(player => player.ID == team.PlayerOne)
          const playerTwo = players.find(player => player.ID == team.PlayerTwo)

          let teamNaming = "[No Team name]"
          if (team.Name != "") {
            teamNaming = team.Name
          }
          const teamName = document.createElement("p")

          teamName.innerHTML = `Team Name: "${teamNaming}" Player One: "${playerOne.FirstName} ${playerOne.LastName}" Player Two: "${playerTwo.FirstName} ${playerTwo.LastName}"`
          teamEl.appendChild(teamName)

          const editDelete = document.createElement("div")
          editDelete.classList.add("edit-delete-button")
          editDelete.appendChild(teamEditButton)
          editDelete.appendChild(teamDeleteButton)

          teamEl.appendChild(editDelete)

          teams.appendChild(teamEl)
        }
      })
    }
  }
}

showAllTeams()
