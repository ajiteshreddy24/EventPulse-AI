import { useEffect, useState } from "react";
import { getEvents } from "../api";

export default function Events() {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchEvents() {
      try {
        const data = await getEvents();
        setEvents(data);
      } catch (err) {
        console.error("Failed to fetch events", err);
      } finally {
        setLoading(false);
      }
    }

    fetchEvents();
  }, []);

  if (loading) return <p>Loading events...</p>;

  return (
    <div className="page">
      <h2>Upcoming Events</h2>

      {events.length === 0 && <p>No events found.</p>}

      {events.map((event) => (
        <div key={event.id} className="card">
          <h3>{event.title}</h3>
          <p>{event.description}</p>
          <p><b>Location:</b> {event.location}</p>
          <p>
            <b>Date:</b> {new Date(event.event_date).toLocaleString()}
          </p>
        </div>
      ))}
    </div>
  );
}
