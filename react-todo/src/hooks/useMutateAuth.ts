import axios from 'axios'
import { useNavigate } from 'react-router-dom'
import useStore from '../store'
import { useError } from './useError'
import { useMutation } from '@tanstack/react-query'
import type { Credentials, RegisterUser } from '../types'

export const useMutateAuth = () => {
  const navigate = useNavigate()
  const resetEditedTask = useStore((state) => state.resetEditedTask)
  const { switchErrorHandling } = useError()

  const loginMutation = useMutation({
    mutationFn: async (user: Credentials) => {
      const res = await axios.post(
        `${import.meta.env.VITE_API_URL}/login`,
        user
      )
      return res.data
    },
    onSuccess: () => {
      navigate('/todo')
    },
    onError: (err: any) => {
      if (err.response.data.message) {
        switchErrorHandling(err.response.data.message)
      } else {
        switchErrorHandling(err.response.data)
      }
    },
  })

  const registerMutation = useMutation({
    mutationFn: async (user: RegisterUser) => {
      const res = await axios.post(
        `${import.meta.env.VITE_API_URL}/signup`,
        user
      )
      return res.data
    },
    onError: (err: any) => {
      if (err.response.data.message) {
        switchErrorHandling(err.response.data.message)
      } else {
        switchErrorHandling(err.response.data)
      }
    },
  })

  const logoutMutation = useMutation({
    mutationFn: async () => {
      const res = await axios.post(`${import.meta.env.VITE_API_URL}/logout`)
      return res.data
    },
    onSuccess: () => {
      resetEditedTask()
      navigate('/')
    },
    onError: (err: any) => {
      if (err.response.data.message) {
        switchErrorHandling(err.response.data.message)
      } else {
        switchErrorHandling(err.response.data)
      }
    },
  })

  return { loginMutation, registerMutation, logoutMutation }
}
