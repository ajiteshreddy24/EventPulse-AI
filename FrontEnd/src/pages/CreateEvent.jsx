import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { createEvent } from "../api";

export default function CreateEvent() {
  const navigate = useNavigate();

  const [form, setForm] = useState({
    title: "",
    description: "",
    location: "",
    event_date: ""
  });

  const [message, setMessage] = useState("");

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      const payload = {
        ...form,
        event_date: new Date(form.event_date).toISOString()
      };

      await createEvent(payload);
      setMessage("Event created successfully!");

      setTimeout(() => {
        navigate("/events");
      }, 1000);
    } catch (err) {
      setMessage("Failed to create event");
    }
  };

  return (
    <div className="page">
      <h2>Create Event</h2>

      <form onSubmit={handleSubmit}>
        <input name="title" placeholder="Title" onChange={handleChange} required />
        <input name="description" placeholder="Description" onChange={handleChange} required />
        <input name="location" placeholder="Location" onChange={handleChange} required />
        <input type="datetime-local" name="event_date" onChange={handleChange} required />
        <button type="submit">Create Event</button>
      </form>

      {message && <p>{message}</p>}
    </div>
  );
}
