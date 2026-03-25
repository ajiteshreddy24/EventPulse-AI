import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Navbar from './components/Navbar'
import Home from './pages/Home'
import CreateEvent from './pages/CreateEvent'
import Events from './pages/Events'
import EditEvent from './pages/EditEvent'

export default function App() {
  return (
    <BrowserRouter>
      <Navbar />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/create" element={<CreateEvent />} />
        <Route path="/events" element={<Events />} />
        <Route path="/edit/:id" element={<EditEvent />} />
      </Routes>
    </BrowserRouter>
  )
}