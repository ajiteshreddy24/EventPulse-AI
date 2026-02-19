import { BrowserRouter, Routes, Route, Link } from "react-router-dom";
import Home from "./pages/Home";
import CreateEvent from "./pages/CreateEvent";
import Events from "./pages/Events";

export default function App() {
  return (
    <BrowserRouter>
      <nav style={{ padding: "10px", background: "#f4f4f4" }}>
        <Link to="/" style={{ marginRight: "10px" }}>Home</Link>
        <Link to="/create" style={{ marginRight: "10px" }}>Create Event</Link>
        <Link to="/events">Upcoming Events</Link>
      </nav>

      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/create" element={<CreateEvent />} />
        <Route path="/events" element={<Events />} />
      </Routes>
    </BrowserRouter>
  );
}
