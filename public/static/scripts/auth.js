export async function isLoggedIn() {
  const accessToken = localStorage.getItem("access-token")

  const body = {
    accessToken,
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

export function setAccessTokenAndUser({ success, payload }) {
  if (success) {
    localStorage.setItem("access-token", payload.accessToken);
    localStorage.setItem("userId", payload.user.ID)
    localStorage.setItem("username", payload.user.Username)
    return true
  } else {
    localStorage.clear("access-token")
    localStorage.clear("userId")
    localStorage.clear("username")
    return false
  }
}
