import axios from 'axios'

const service = axios.create({
  baseURL: process.env.VITE_APP_API_BASE_URL,
  withCredentials: true,
  timeout: 10000
})

service.defaults.headers.post['Content-Type'] = 'application/json;charset=UTF-8'

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
