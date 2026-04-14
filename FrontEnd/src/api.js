export async function getEvents() {
    const res = await fetch('/api/events')
    if (!res.ok) throw new Error('Failed to fetch latest events')
    return res.json()
  }
  
  export async function createEvent(data) {
    const res = await fetch('/api/events', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    })
  
    if (!res.ok) throw new Error(await res.text())
    return res.json()
  }
  
  export async function updateEvent(id, data) {
    const res = await fetch(`/api/events/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    })
  
    if (!res.ok) throw new Error(await res.text())
    return res.json()
  }
  
  export async function deleteEvent(id) {
    const res = await fetch(`/api/events/${id}`, {
      method: 'DELETE',
    })
  
    if (!res.ok) throw new Error('Delete failed')
  }
