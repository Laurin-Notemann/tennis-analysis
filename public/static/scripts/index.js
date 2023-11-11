async function loadNavBar() {
  const { success, payload } = isLoggedIn()
  if (success) {
    console.log("test1")
  } else {
    console.log("test2")
  }
}

loadNavBar()

async function isLoggedIn() {
  const accessToken = localStorage.getItem("access-token")

  const body = {
    accessToken,
    refreshToken
  }

  const res = await fetch("/api/refresh", {
    method: "POST",
    body: JSON.stringify(body),
    headers: {
      "Content-Type": "application/json"
    }
  })

  const payload = await res.json()

  if (res.status == 400) {
    return {
      success: false,
      payload,
    }
  } else if (res.status == 401) {
    return {
      success: false,
      payload,
    }
  } else if (res.status == 500) {
    return {
      success: false,
      payload,
    }
  }
  return {
    success: true,
    payload
  }
}
