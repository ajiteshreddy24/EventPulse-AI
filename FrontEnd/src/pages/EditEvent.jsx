import { useParams, useNavigate } from 'react-router-dom'
import { useState } from 'react'
import { updateEvent } from '../api'

export default function EditEvent() {
  const { id } = useParams()
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

    await updateEvent(id, payload)
    navigate('/events')
  }

  return (
    <div className="container">
      <h2>Edit Event</h2>

      <form onSubmit={handleSubmit}>
        <input name="title" onChange={handleChange} required />
        <textarea name="description" onChange={handleChange} required />
        <input name="location" onChange={handleChange} required />
        <input type="datetime-local" name="event_date" onChange={handleChange} required />
        <button type="submit">Update Event</button>
      </form>
    </div>
  )
}