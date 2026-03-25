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
    loadEvents()
  }

  return (
    <div className="container">
      <h2 style={{ textAlign: 'center', marginBottom: '30px' }}>
        🚀 Upcoming Events
      </h2>

      {events.length === 0 && (
        <p style={{ textAlign: 'center', opacity: 0.6 }}>
          No events yet. Create something awesome ✨
        </p>
      )}

      {events.map((e) => (
        <div className="card" key={e.id}>
          <h3>{e.title}</h3>

          <p style={{ marginTop: '8px' }}>{e.description}</p>

          <p style={{ marginTop: '6px' }}>
            📍 {e.location}
          </p>

          <p style={{ marginTop: '6px' }}>
            🗓 {new Date(e.event_date).toLocaleString()}
          </p>

          <div style={{ marginTop: '15px', display: 'flex', gap: '15px' }}>
            <Link to={`/edit/${e.id}`}>
              ✏️ Edit
            </Link>

            <button
              style={{
                background: 'linear-gradient(135deg, #ef4444, #dc2626)',
              }}
              onClick={() => handleDelete(e.id)}
            >
              🗑 Delete
            </button>
          </div>
        </div>
      ))}
    </div>
  )
}