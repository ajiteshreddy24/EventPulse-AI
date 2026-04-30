import { useContext, useMemo, useState } from "react"
import { updateInterests } from "../api"
import { AuthContext } from "../context/auth-context"
import { INTEREST_OPTIONS } from "../constants/interests"

export default function Profile() {
  const { user, updateUser } = useContext(AuthContext)
  const [selectedInterests, setSelectedInterests] = useState(user?.interests || [])
  const [status, setStatus] = useState("")
  const [error, setError] = useState("")
  const title = useMemo(() => user?.name || "Your profile", [user?.name])

  function toggleInterest(interest) {
    setSelectedInterests((current) =>
      current.includes(interest)
        ? current.filter((item) => item !== interest)
        : [...current, interest]
    )
  }

  async function handleSave() {
    try {
      const updatedUser = await updateInterests(selectedInterests)
      updateUser(updatedUser)
      setStatus("Interests saved. Recommendations will update from these choices.")
      setError("")
    } catch (err) {
      setError(err.message || "Unable to save interests")
      setStatus("")
    }
  }

  return (
    <div className="container">
      <div className="card profile-card">
        <h2>{title}</h2>
        <p>Choose the interests you want GatorHive to use for your event recommendations.</p>

        <div className="interest-picker">
          <div className="interest-grid">
            {INTEREST_OPTIONS.map((interest) => (
              <button
                key={interest}
                type="button"
                className={`interest-pill ${selectedInterests.includes(interest) ? "active" : ""}`}
                onClick={() => toggleInterest(interest)}
              >
                {interest}
              </button>
            ))}
          </div>
        </div>

        {error && <p style={{ color: "crimson" }}>{error}</p>}
        {status && <p style={{ color: "green" }}>{status}</p>}

        <button type="button" onClick={handleSave}>
          Save Interests
        </button>
      </div>
    </div>
  )
}
