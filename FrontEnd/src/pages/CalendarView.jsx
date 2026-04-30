import FullCalendar from "@fullcalendar/react"
import dayGridPlugin from "@fullcalendar/daygrid"
import { useEffect, useState } from "react"
import { getEvents } from "../api"

export default function CalendarView() {
  const [events, setEvents] = useState([])
  const [error, setError] = useState("")

  useEffect(() => {
    async function loadCalendarEvents() {
      try {
        const data = await getEvents()
        const items = Array.isArray(data) ? data : []

        setEvents(
          items.map((e) => ({
            title: e.title,
            date: e.event_date,
          }))
        )
        setError("")
      } catch (err) {
        console.error("Failed to load calendar events:", err)
        setEvents([])
        setError("Unable to load the calendar right now.")
      }
    }

    loadCalendarEvents()
  }, [])

  return (
    <div className="container">
      <h2>Calendar</h2>
      {error && <p style={{ color: "crimson" }}>{error}</p>}
      <FullCalendar
        plugins={[dayGridPlugin]}
        initialView="dayGridMonth"
        events={events}
      />
    </div>
  )
}
