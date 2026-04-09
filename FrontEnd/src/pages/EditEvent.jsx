import { useEffect, useState, useContext } from "react"
import { useNavigate, useParams } from "react-router-dom"
import { getEvent, updateEvent } from "../api"
import { AuthContext } from "../context/auth-context"

export default function EditEvent() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { user } = useContext(AuthContext)
  const [form, setForm] = useState(null)
  const [error, setError] = useState("")

  useEffect(() => {
    async function load() {
      try {
        const event = await getEvent(id)
        setForm({
          ...event,
          event_date: event.event_date.slice(0, 16),
        })
        setError("")
      } catch (err) {
        console.error(err)
        setError(err.message || "Unable to load event")
      }
    }

    if (user) {
      load()
    }
  }, [id, user])

  if (!form) {
    return (
      <div className="container">
        <h2>Edit Event</h2>
        {error ? <p style={{ color: "crimson" }}>{error}</p> : <p>Loading...</p>}
      </div>
    )
  }

  const handleChange = (e) =>
    setForm({ ...form, [e.target.name]: e.target.value })

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      await updateEvent(id, { ...form, event_date: new Date(form.event_date).toISOString() })
      navigate("/events")
    } catch (err) {
      console.error(err)
      setError(err.message || "Update failed")
    }
  }

  return (
    <div className="container">
      <h2>Edit Event</h2>
      {error && <p style={{ color: "crimson" }}>{error}</p>}
      <form onSubmit={handleSubmit}>
        <input data-cy="title-input" name="title" value={form.title} onChange={handleChange} required />
        <textarea name="description" value={form.description} onChange={handleChange} required />
        <input name="location" value={form.location} onChange={handleChange} required />
        <input type="datetime-local" name="event_date" value={form.event_date} onChange={handleChange} required />
        <button data-cy="update-btn">Update Event</button>
      </form>
    </div>
  )
}
