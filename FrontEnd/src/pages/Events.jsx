import { useEffect, useState, useContext } from "react"
import {
  getEvents,
  deleteEvent,
  rsvpEvent,
  cancelRSVP,
  joinWaitlist,
  getComments,
  createComment,
  deleteComment,
  getRecommendations,
} from "../api"
import { Link, useNavigate } from "react-router-dom"
import { AuthContext } from "../context/auth-context"

export default function Events() {
  const [events, setEvents] = useState([])
  const [error, setError] = useState("")
  const [searchTerm, setSearchTerm] = useState("")
  const [commentsByEvent, setCommentsByEvent] = useState({})
  const [visibleComments, setVisibleComments] = useState({})
  const [commentDrafts, setCommentDrafts] = useState({})
  const [recommendations, setRecommendations] = useState([])
  const [recommendationError, setRecommendationError] = useState("")
  const { user } = useContext(AuthContext)
  const navigate = useNavigate()

  async function fetchEventsData() {
    const data = await getEvents()
    return Array.isArray(data) ? data : []
  }

  async function refreshEvents() {
    try {
      const data = await fetchEventsData()
      setEvents(data)
      setError("")
    } catch (err) {
      console.error(err)
      setError(err.message || "Unable to load events")
      setEvents([])
    }
  }

  async function refreshRecommendations() {
    if (!user) {
      setRecommendations([])
      setRecommendationError("")
      return
    }

    try {
      const data = await getRecommendations()
      setRecommendations(data)
      setRecommendationError("")
    } catch (err) {
      console.error(err)
      setRecommendationError(err.message || "Unable to load recommendations")
      setRecommendations([])
    }
  }

  useEffect(() => {
    let isMounted = true

    async function loadInitialEvents() {
      try {
        const data = await fetchEventsData()
        if (!isMounted) return

        setEvents(data)
        setError("")
      } catch (err) {
        console.error(err)
        if (!isMounted) return

        setError(err.message || "Unable to load events")
        setEvents([])
      }
    }

    loadInitialEvents()

    return () => {
      isMounted = false
    }
  }, [])

  useEffect(() => {
    let isMounted = true

    async function loadRecommendations() {
      if (!user) {
        if (!isMounted) return
        setRecommendations([])
        setRecommendationError("")
        return
      }

      try {
        const data = await getRecommendations()
        if (!isMounted) return

        setRecommendations(data)
        setRecommendationError("")
      } catch (err) {
        console.error(err)
        if (!isMounted) return

        setRecommendationError(err.message || "Unable to load recommendations")
        setRecommendations([])
      }
    }

    loadRecommendations()

    return () => {
      isMounted = false
    }
  }, [user])

  async function handleDelete(id) {
    if (!user) return navigate("/login")
    try {
      await deleteEvent(id)
      await refreshEvents()
    } catch (err) {
      console.error(err)
      setError(err.message || "Delete failed")
    }
  }

  async function handleRSVP(event) {
    if (!user) return navigate("/login")
    try {
      if (event.userHasRSVP) {
        await cancelRSVP(event.id)
      } else {
        await rsvpEvent(event.id)
      }
      await refreshEvents()
      await refreshRecommendations()
    } catch (err) {
      console.error(err)
      setError(err.message || "RSVP failed")
    }
  }

  async function handleJoinWaitlist(event) {
    if (!user) return navigate("/login")

    try {
      await joinWaitlist(event.id)
      await refreshEvents()
      await refreshRecommendations()
    } catch (err) {
      console.error(err)
      setError(err.message || "Unable to join waitlist")
    }
  }

  async function toggleComments(eventId) {
    const isVisible = visibleComments[eventId]

    if (!isVisible && !commentsByEvent[eventId]) {
      try {
        const comments = await getComments(eventId)
        setCommentsByEvent((current) => ({ ...current, [eventId]: comments }))
      } catch (err) {
        console.error(err)
        setError(err.message || "Unable to load comments")
      }
    }

    setVisibleComments((current) => ({ ...current, [eventId]: !isVisible }))
  }

  async function handleCommentSubmit(eventId) {
    if (!user) {
      navigate("/login")
      return
    }

    const draft = (commentDrafts[eventId] || "").trim()
    if (!draft) {
      return
    }

    try {
      const comment = await createComment(eventId, draft)
      setCommentsByEvent((current) => ({
        ...current,
        [eventId]: [...(current[eventId] || []), comment],
      }))
      setCommentDrafts((current) => ({ ...current, [eventId]: "" }))
      setVisibleComments((current) => ({ ...current, [eventId]: true }))
    } catch (err) {
      console.error(err)
      setError(err.message || "Unable to post comment")
    }
  }

  async function handleDeleteComment(eventId, commentId) {
    try {
      await deleteComment(commentId)
      setCommentsByEvent((current) => ({
        ...current,
        [eventId]: (current[eventId] || []).filter((comment) => comment.id !== commentId),
      }))
    } catch (err) {
      console.error(err)
      setError(err.message || "Unable to delete comment")
    }
  }

  const filteredEvents = events.filter((event) => {
    const query = searchTerm.trim().toLowerCase()
    if (!query) {
      return true
    }

    return [event.title, event.location, event.description]
      .filter(Boolean)
      .some((value) => value.toLowerCase().includes(query))
  })

  return (
    <div className="container">
      <h2 style={{ textAlign: "center", marginBottom: "40px" }}>Campus Events</h2>

      {error && <p style={{ color: "crimson", textAlign: "center" }}>{error}</p>}

      {user && (
        <section className="recommendation-panel">
          <div className="section-header">
            <h3>AI Recommendations</h3>
            {recommendationError && <p style={{ color: "crimson" }}>{recommendationError}</p>}
          </div>

          {recommendations.length > 0 ? (
            <div className="recommendation-grid">
              {recommendations.map((event) => (
                <article key={`recommendation-${event.id}`} className="recommendation-card">
                  <p className="recommendation-tag">Suggested for you</p>
                  <h4>{event.title}</h4>
                  <p>{event.recommendationReason}</p>
                  <span>{new Date(event.event_date).toLocaleString()}</span>
                </article>
              ))}
            </div>
          ) : (
            <p className="recommendation-empty">
              {user?.interests?.length
                ? "No event matches your saved interests yet. Try checking back after more events are added."
                : "Add interests in your profile to unlock personalized recommendations."}
            </p>
          )}
        </section>
      )}

      <div className="event-toolbar">
        <input
          data-cy="search-input"
          className="search-input"
          placeholder="Search by title, location, or description"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
      </div>

      {filteredEvents.length === 0 && <p style={{ textAlign: "center" }}>No matching events found.</p>}

      <div className="event-grid">
        {filteredEvents.map((e) => {
          const isFull = e.capacity > 0 && e.rsvpCount >= e.capacity
          const comments = commentsByEvent[e.id] || []
          const commentsVisible = visibleComments[e.id]

          return (
          <div className="card" key={e.id}>
            <h3>{e.title}</h3>
            <p>{e.description}</p>
            <div>Location: {e.location}</div>
            <div>Date: {new Date(e.event_date).toLocaleString()}</div>
            <div>Capacity: {e.rsvpCount}/{e.capacity}</div>
            <div>Waitlist: {e.waitlistCount}</div>

            <div className="event-actions">
              {user && (
                <Link to={`/edit/${e.id}`}>
                  <button data-cy="edit-btn">Edit</button>
                </Link>
              )}
              {user && (
                <button data-cy="delete-btn" onClick={() => handleDelete(e.id)}>Delete</button>
              )}

              {e.userHasRSVP ? (
                <button data-cy="rsvp-btn" onClick={() => handleRSVP(e)}>
                  Cancel RSVP
                </button>
              ) : isFull ? (
                <button
                  data-cy="waitlist-btn"
                  onClick={() => handleJoinWaitlist(e)}
                  disabled={e.userOnWaitlist}
                >
                  {e.userOnWaitlist ? "Joined Waitlist" : "Join Waitlist"}
                </button>
              ) : (
                <button data-cy="rsvp-btn" onClick={() => handleRSVP(e)}>
                  RSVP
                </button>
              )}
            </div>

            <div className="comment-section">
              <button
                type="button"
                className="comment-toggle"
                data-cy="comments-toggle"
                onClick={() => toggleComments(e.id)}
              >
                {commentsVisible ? "Hide Comments" : "Show Comments"}
              </button>

              {commentsVisible && (
                <div className="comment-thread">
                  {comments.length === 0 && <p>No comments yet.</p>}

                  {comments.map((comment) => (
                    <div key={comment.id} className="comment-item">
                      <div className="comment-meta">
                        <strong>{comment.user?.name || "User"}</strong>
                        <span>{new Date(comment.created_at).toLocaleString()}</span>
                      </div>
                      <p>{comment.content}</p>
                      {(user?.id === comment.user_id || user) && (
                        <button
                          type="button"
                          className="comment-delete"
                          onClick={() => handleDeleteComment(e.id, comment.id)}
                        >
                          Delete Comment
                        </button>
                      )}
                    </div>
                  ))}

                  <div className="comment-compose">
                    <textarea
                      data-cy="comment-input"
                      placeholder="Join the discussion"
                      value={commentDrafts[e.id] || ""}
                      onChange={(event) =>
                        setCommentDrafts((current) => ({
                          ...current,
                          [e.id]: event.target.value,
                        }))
                      }
                    />
                    <button
                      type="button"
                      data-cy="comment-submit"
                      onClick={() => handleCommentSubmit(e.id)}
                    >
                      Post Comment
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        )})}
      </div>
    </div>
  )
}
