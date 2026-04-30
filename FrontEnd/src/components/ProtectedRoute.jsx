import { useContext } from "react"
import { Navigate } from "react-router-dom"
import { AuthContext } from "../context/auth-context"

export default function ProtectedRoute({ children }) {
  const { user, loading } = useContext(AuthContext)

  if (loading) return null

  return user ? children : <Navigate to="/login" />
}
