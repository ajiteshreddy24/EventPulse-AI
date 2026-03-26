import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { createEvent } from '../api'

export default function CreateEvent() {
  const navigate = useNavigate()

  const [form, setForm] = useState({
    title: '',
    description: '',
    location: '',
    event_date: '',
  })

  const handleChange = (e) =>
    setForm({ ...form, [e.target.name]: e.target.value })

  const handleSubmit = async (e) => {
    e.preventDefault()

    const payload = {
      ...form,
      event_date: new Date(form.event_date).toISOString(),
    }

    await createEvent(payload)
    navigate('/events')
  }

  return (
    <div className="container">
      <h2>Create Event</h2>

      <form onSubmit={handleSubmit}>
        <input
          data-cy="title-input"
          name="title"
          placeholder="Title"
          onChange={handleChange}
          required
        />

        <textarea
          data-cy="description-input"
          name="description"
          placeholder="Description"
          onChange={handleChange}
          required
        />

        <input
          data-cy="location-input"
          name="location"
          placeholder="Location"
          onChange={handleChange}
          required
        />

        <input
          data-cy="date-input"
          type="datetime-local"
          name="event_date"
          onChange={handleChange}
          required
        />

        <button data-cy="create-btn" type="submit">
          Create Event
        </button>
      </form>
    </div>
  )
}