import { Link } from "react-router-dom";

export default function Navbar() {
  return (
    <nav style={styles.nav}>
      <h3>EventPulse-AI</h3>
      <div>
        <Link to="/">Home</Link>
        <Link to="/create">Create Event</Link>
        <Link to="/events">Events</Link>
      </div>
    </nav>
  );
}

const styles = {
  nav: {
    display: "flex",
    justifyContent: "space-between",
    padding: "15px",
    background: "#0021A5",
    color: "white"
  }
};
