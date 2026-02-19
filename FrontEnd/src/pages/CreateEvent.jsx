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
    await createEvent(form);
    setMessage("Event created successfully!");

    setTimeout(() => {
      navigate("/events");
    }, 1000);
  };

  return (
    <div className="page">
      <h2>Create Event</h2>
      <form onSubmit={handleSubmit}>
        <input name="title" value={form.title} onChange={handleChange} />
        <input name="description" value={form.description} onChange={handleChange} />
        <input name="location" value={form.location} onChange={handleChange} />
        <input type="datetime-local" name="event_date" value={form.event_date} onChange={handleChange} />
        <button type="submit">Create Event</button>
      </form>
      {message && <p className="success">{message}</p>}
    </div>
  );
}
