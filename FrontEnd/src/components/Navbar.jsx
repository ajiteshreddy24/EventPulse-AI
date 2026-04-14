import { Link } from 'react-router-dom'

export default function Navbar() {
  return (
    <nav className="navbar">
      <Link to="/">Home</Link>
      <Link to="/create">Create New Event</Link>
      <Link to="/events">Upcoming Events</Link>
    </nav>
  )
}
