import { useMutation, useQueryClient } from '@tanstack/react-query'
import axios from 'axios'
import { useError } from './useError'
import useStore from '../store'
import type { Task } from '../types'

export const useMutateTask = () => {
  const queryClient = useQueryClient()
  const { switchErrorHandling } = useError()
  const resetEditedTask = useStore((state) => state.resetEditedTask)

  const createTaskMutation = useMutation({
    mutationFn: async (
      task: Omit<Task, 'id' | 'created_at' | 'updated_at'>
    ) => {
      const res = await axios.post(
        `${import.meta.env.VITE_API_URL}/tasks`,
        task
      )
      return res.data
    },
    onSuccess: (data) => {
      const previousTasks = queryClient.getQueryData<Task[]>(['tasks'])
      if (previousTasks) {
        queryClient.setQueryData(['tasks'], [...previousTasks, data])
      }
      resetEditedTask()
    },
    onError: (err: any) => {
      if (err.response.data.message) {
        switchErrorHandling(err.response.data.message)
      } else {
        switchErrorHandling(err.response.data)
      }
    },
  })

  const updateTaskMutation = useMutation({
    mutationFn: async (task: Omit<Task, 'created_at' | 'updated_at'>) => {
      const res = await axios.put(
        `${import.meta.env.VITE_API_URL}/tasks/${task.id}`,
        task
      )
      return res.data
    },
    onSuccess: (data, variables) => {
      const previousTasks = queryClient.getQueryData<Task[]>(['tasks'])
      if (previousTasks) {
        queryClient.setQueryData<Task[]>(
          ['tasks'],
          previousTasks.map((task) => (task.id === variables.id ? data : task))
        )
      }
      resetEditedTask()
    },
    onError: (err: any) => {
      if (err.response.data.message) {
        switchErrorHandling(err.response.data.message)
      } else {
        switchErrorHandling(err.response.data)
      }
    },
  })

  const deleteTaskMutation = useMutation({
    mutationFn: async (id: number) => [
      `${import.meta.env.VITE_API_URL}/tasks/${id}`,
    ],
    onSuccess: (_, variables) => {
      const previousTasks = queryClient.getQueryData<Task[]>(['tasks'])
      if (previousTasks) {
        queryClient.setQueryData<Task[]>(
          ['tasks'],
          previousTasks.filter((task) => task.id !== variables)
        )
      }
      resetEditedTask()
    },
    onError: (err: any) => {
      if (err.response.data.message) {
        switchErrorHandling(err.response.data.message)
      } else {
        switchErrorHandling(err.response.data)
      }
    },
  })
  return {
    createTaskMutation,
    updateTaskMutation,
    deleteTaskMutation,
  }
}
