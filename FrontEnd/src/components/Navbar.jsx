import { Link } from "react-router-dom"
import { useContext } from "react"
import { AuthContext } from "../context/auth-context"

export default function Navbar() {
  const { user, logout } = useContext(AuthContext)

  return (
    <nav className="navbar">
      <Link className="brand-mark" to="/">
        <span className="brand-dot" />
        <span>EventPulse AI</span>
      </Link>

      <div className="nav-links">
        <Link to="/events">Events</Link>
        <Link to="/calendar">Calendar</Link>
        <Link to="/create">Create Event</Link>
      </div>

      <div className="nav-actions">
        {user ? (
          <button type="button" className="ghost-button" onClick={logout}>
            Logout
          </button>
        ) : (
          <>
            <Link to="/login">Login</Link>
            <Link className="nav-pill" to="/signup">
              Signup
            </Link>
          </>
        )}
      </div>
    </nav>
  )
}
