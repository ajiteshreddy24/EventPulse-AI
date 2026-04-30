import { Link } from "react-router-dom"

export default function Home() {
  return (
    <div className="home">
      <div className="home-backdrop" />
      <div className="home-gator-stripes" />
      <div className="home-gator-mark" aria-hidden="true">
        <span className="gator-eye" />
        <span className="gator-jaw gator-jaw-top" />
        <span className="gator-jaw gator-jaw-bottom" />
      </div>

      <section className="hero-shell">
        <div className="hero-copy">
          <p className="eyebrow">Smart campus event discovery</p>
          <h1>Find the right events without the noise.</h1>
          <p className="hero-text">
            EventPulse student community brings student events into one calm,
            searchable place so people can discover, plan, and engage faster.
          </p>

          <div className="hero-actions">
            <Link className="hero-primary" to="/events">
              Explore Events
            </Link>
            <Link className="hero-secondary" to="/create">
              Create an Event
            </Link>
          </div>

          <div className="hero-stats">
            <div className="stat-card">
              <strong>All in one</strong>
              <span>Events, calendar, RSVP, and updates in one flow.</span>
            </div>
            <div className="stat-card">
              <strong>Built for clarity</strong>
              <span>Less clutter, better focus, and faster decisions.</span>
            </div>
            <div className="stat-card">
              <strong>Campus-ready</strong>
              <span>Made for clubs, organizers, students, and communities.</span>
            </div>
          </div>
        </div>

        <div className="hero-panel">
          <div className="panel-window">
            <div className="panel-row">
              <span className="panel-label">Today</span>
              <span className="panel-chip">3 live picks</span>
            </div>

            <div className="panel-event featured">
              <p>Design Jam</p>
              <span>Innovation Lab • 6:00 PM</span>
            </div>

            <div className="panel-event">
              <p>Startup Mixer</p>
              <span>North Hall • 7:30 PM</span>
            </div>

            <div className="panel-event">
              <p>Open Mic Night</p>
              <span>Student Center • 8:00 PM</span>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}
