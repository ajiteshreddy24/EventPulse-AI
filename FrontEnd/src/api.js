function authHeaders() {
  const token = localStorage.getItem("token")
  return {
    "Content-Type": "application/json",
    ...(token && { Authorization: `Bearer ${token}` }),
  }
}

async function parseJsonResponse(res, fallbackMessage) {
  const text = await res.text()

  if (!res.ok) {
    throw new Error(text || fallbackMessage)
  }

  return text ? JSON.parse(text) : null
}

/* ================= EVENTS ================= */

export async function getEvents() {
  const res = await fetch("/api/events", { headers: authHeaders() })
  const data = await parseJsonResponse(res, "Failed to fetch events")
  return Array.isArray(data) ? data : []
}

export async function getEvent(id) {
  const res = await fetch(`/api/events/${id}`, { headers: authHeaders() })
  return parseJsonResponse(res, "Failed to fetch event")
}

export async function createEvent(data) {
  const res = await fetch("/api/events", {
    method: "POST",
    headers: authHeaders(),
    body: JSON.stringify(data),
  })
  return parseJsonResponse(res, "Failed to create event")
}

export async function updateEvent(id, data) {
  const res = await fetch(`/api/events/${id}`, {
    method: "PUT",
    headers: authHeaders(),
    body: JSON.stringify(data),
  })
  return parseJsonResponse(res, "Failed to update event")
}

export async function deleteEvent(id) {
  const res = await fetch(`/api/events/${id}`, {
    method: "DELETE",
    headers: authHeaders(),
  })
  if (!res.ok) throw new Error(await res.text() || "Delete failed")
}

export async function rsvpEvent(id) {
  const res = await fetch(`/api/events/${id}/rsvp`, {
    method: "POST",
    headers: authHeaders(),
  })
  if (!res.ok) throw new Error(await res.text() || "RSVP failed")
}

export async function cancelRSVP(id) {
  const res = await fetch(`/api/events/${id}/rsvp`, {
    method: "DELETE",
    headers: authHeaders(),
  })
  if (!res.ok) throw new Error(await res.text() || "Cancel RSVP failed")
}

/* ================= AUTH ================= */

export async function login(data) {
  const res = await fetch("/api/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  })
  return parseJsonResponse(res, "Login failed")
}

export async function signup(data) {
  const res = await fetch("/api/auth/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  })
  return parseJsonResponse(res, "Signup failed")
}

export async function getMe() {
  const res = await fetch("/api/auth/me", { headers: authHeaders() })
  if (!res.ok) return null
  return parseJsonResponse(res, "Failed to fetch current user")
}
