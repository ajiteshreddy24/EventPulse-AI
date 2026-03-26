import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { getEvents, updateEvent } from '../api'

export default function EditEvent() {
  const { id } = useParams()
  const navigate = useNavigate()

  const [form, setForm] = useState({
    title: '',
    description: '',
    location: '',
    event_date: '',
  })

  useEffect(() => {
    async function loadEvent() {
      const events = await getEvents()
      const event = events.find((e) => e.id === parseInt(id))

      if (event) {
        setForm({
          ...event,
          event_date: event.event_date.slice(0, 16),
        })
      }
    }

    loadEvent()
  }, [id])

  const handleChange = (e) =>
    setForm({ ...form, [e.target.name]: e.target.value })

  const handleSubmit = async (e) => {
    e.preventDefault()

    await updateEvent(id, {
      ...form,
      event_date: new Date(form.event_date).toISOString(),
    })

    navigate('/events')
  }

  return (
    <div className="container">
      <h2>Edit Event</h2>

      <form onSubmit={handleSubmit}>
        <input
          data-cy="title-input"
          name="title"
          value={form.title}
          onChange={handleChange}
          required
        />

        <textarea
          data-cy="description-input"
          name="description"
          value={form.description}
          onChange={handleChange}
          required
        />

        <input
          data-cy="location-input"
          name="location"
          value={form.location}
          onChange={handleChange}
          required
        />

        <input
          data-cy="date-input"
          type="datetime-local"
          name="event_date"
          value={form.event_date}
          onChange={handleChange}
          required
        />

        <button data-cy="update-btn">
          Update Event
        </button>
      </form>
    </div>
  )
}