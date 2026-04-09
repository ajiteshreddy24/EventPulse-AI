import { useEffect, useState, useContext } from "react"
import { getEvents, deleteEvent, rsvpEvent, cancelRSVP } from "../api"
import { Link, useNavigate } from "react-router-dom"
import { AuthContext } from "../context/auth-context"

export default function Events() {
  const [events, setEvents] = useState([])
  const [error, setError] = useState("")
  const { user } = useContext(AuthContext)
  const navigate = useNavigate()

  async function fetchEventsData() {
    const data = await getEvents()
    return Array.isArray(data) ? data : []
  }

  async function refreshEvents() {
    try {
      const data = await fetchEventsData()
      setEvents(data)
      setError("")
    } catch (err) {
      console.error(err)
      setError(err.message || "Unable to load events")
      setEvents([])
    }
  }

  useEffect(() => {
    let isMounted = true

    async function loadInitialEvents() {
      try {
        const data = await fetchEventsData()
        if (!isMounted) return

        setEvents(data)
        setError("")
      } catch (err) {
        console.error(err)
        if (!isMounted) return

        setError(err.message || "Unable to load events")
        setEvents([])
      }
    }

    loadInitialEvents()

    return () => {
      isMounted = false
    }
  }, [])

  async function handleDelete(id) {
    if (!user) return navigate("/login")
    try {
      await deleteEvent(id)
      await refreshEvents()
    } catch (err) {
      console.error(err)
      setError(err.message || "Delete failed")
    }
  }

  async function handleRSVP(event) {
    if (!user) return navigate("/login")
    try {
      if (event.userHasRSVP) {
        await cancelRSVP(event.id)
      } else {
        await rsvpEvent(event.id)
      }
      await refreshEvents()
    } catch (err) {
      console.error(err)
      setError(err.message || "RSVP failed")
    }
  }

  return (
    <div className="container">
      <h2 style={{ textAlign: "center", marginBottom: "40px" }}>Campus Events</h2>

      {error && <p style={{ color: "crimson", textAlign: "center" }}>{error}</p>}

      {events.length === 0 && <p style={{ textAlign: "center" }}>No events yet.</p>}

      <div className="event-grid">
        {events.map((e) => (
          <div className="card" key={e.id}>
            <h3>{e.title}</h3>
            <p>{e.description}</p>
            <div>Location: {e.location}</div>
            <div>Date: {new Date(e.event_date).toLocaleString()}</div>

            <div style={{ marginTop: "15px", display: "flex", gap: "10px" }}>
              {user && (
                <Link to={`/edit/${e.id}`}>
                  <button data-cy="edit-btn">Edit</button>
                </Link>
              )}
              {user && (
                <button data-cy="delete-btn" onClick={() => handleDelete(e.id)}>Delete</button>
              )}
              <button onClick={() => handleRSVP(e)}>
                {e.userHasRSVP ? "Cancel RSVP" : "RSVP"}
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
