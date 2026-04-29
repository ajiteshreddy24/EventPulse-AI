import { useState, useContext } from "react"
import { signup } from "../api"
import { AuthContext } from "../context/auth-context"
import { useNavigate } from "react-router-dom"
import { INTEREST_OPTIONS } from "../constants/interests"

export default function Signup() {
  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [interests, setInterests] = useState([])
  const { login } = useContext(AuthContext)
  const navigate = useNavigate()

  function toggleInterest(interest) {
    setInterests((current) =>
      current.includes(interest)
        ? current.filter((item) => item !== interest)
        : [...current, interest]
    )
  }

  async function handleSubmit(e) {
    e.preventDefault()
    try {
      const data = await signup({ name, email, password, interests })
      login(data)
      navigate("/events")
    } catch {
      alert("Signup failed")
    }
  }

  return (
    <div className="container">
      <h2>Signup</h2>
      <form onSubmit={handleSubmit}>
        <input
          placeholder="Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <input
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <div className="interest-picker">
          <p>Select your interests for smarter recommendations</p>
          <div className="interest-grid">
            {INTEREST_OPTIONS.map((interest) => (
              <button
                key={interest}
                type="button"
                className={`interest-pill ${interests.includes(interest) ? "active" : ""}`}
                onClick={() => toggleInterest(interest)}
              >
                {interest}
              </button>
            ))}
          </div>
        </div>
        <button type="submit">Signup</button>
      </form>
    </div>
  )
}
