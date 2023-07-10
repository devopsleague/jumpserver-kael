import axios from 'axios'

const service = axios.create({
  baseURL: import.meta.env.VITE_APP_API_BASE_URL,
  withCredentials: true,
  timeout: 10000
})

service.interceptors.request.use(
  (config) => {
    return config
  },
  (error) => {
    return Promise.reject(error.response)
  },
)

service.interceptors.response.use(
  (response) => {
    if (response.status === 200)
      return response

    throw new Error(response.status.toString())
  },
  (error) => {
    return Promise.reject(error)
  },
)

export default service
