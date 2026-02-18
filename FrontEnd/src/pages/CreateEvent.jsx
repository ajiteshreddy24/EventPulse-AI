import { useState } from "react";
import { createEvent } from "../api";

export default function CreateEvent() {
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
    await createEvent(form);
    setMessage("Event created successfully!");
    setForm({ title: "", description: "", location: "", event_date: "" });
  };

  return (
    <div className="page">
      <h2>Create Event</h2>

      <form onSubmit={handleSubmit}>
        <input name="title" placeholder="Title" value={form.title} onChange={handleChange} required />
        <input name="description" placeholder="Description" value={form.description} onChange={handleChange} required />
        <input name="location" placeholder="Location" value={form.location} onChange={handleChange} required />
        <input type="datetime-local" name="event_date" value={form.event_date} onChange={handleChange} required />
        <button type="submit">Create</button>
      </form>

      {message && <p className="success">{message}</p>}
    </div>
  );
}
