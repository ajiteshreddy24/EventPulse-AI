import { useEffect, useState } from 'react'
import { getEvents, deleteEvent } from '../api'
import { Link } from 'react-router-dom'

export default function Events() {
  const [events, setEvents] = useState([])

  useEffect(() => {
    loadEvents()
  }, [])

  async function loadEvents() {
    const data = await getEvents()
    setEvents(data)
  }

  async function handleDelete(id) {
    await deleteEvent(id)
    await loadEvents()
  }

  return (
    <div className="container">
      <h2>Upcoming Events</h2>

      {events.length === 0 && (
        <p>No events yet. Create something awesome ✨</p>
      )}

      {events.map((e) => (
        <div className="card" key={e.id}>
          <h3>{e.title}</h3>
          <p>{e.description}</p>
          <p>{e.location}</p>
          <p>{new Date(e.event_date).toLocaleString()}</p>

          <div style={{ display: 'flex', gap: '10px' }}>
            <Link
              data-cy="edit-btn"
              to={`/edit/${e.id}`}
            >
              Edit
            </Link>

            <button
              data-cy="delete-btn"
              onClick={() => handleDelete(e.id)}
            >
              Delete
            </button>
          </div>
        </div>
      ))}
    </div>
  )
}