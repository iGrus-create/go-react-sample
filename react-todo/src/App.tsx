import { useEffect } from 'react'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Auth from './components/Auth'
import Todo from './components/Todo'
import axios from 'axios'
import type { CsrfToken } from './types'

const App = () => {
  useEffect(() => {
    axios.defaults.withCredentials = true
    const getCsrfToken = async () => {
      const { data } = await axios.get<CsrfToken>(
        `${import.meta.env.VITE_API_URL}/csrf`
      )
      axios.defaults.headers.common['X-CSRF-TOKEN'] = data.csrf_token
    }
    getCsrfToken()
  }, [])

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Auth />} />
        <Route path="/todo" element={<Todo />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
