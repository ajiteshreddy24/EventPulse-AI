const API_URL = "http://localhost:8080";

export async function createEvent(data) {
  return fetch(`${API_URL}/events`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data)
  });
}

export async function getEvents() {
  const res = await fetch(`${API_URL}/events`);
  return res.json();
}
