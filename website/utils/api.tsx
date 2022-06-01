import axios from 'axios';

export const baseURL = process.env.NEXT_PUBLIC_FORM_API_ENDPOINT

const api = axios.create({
  baseURL,
  headers: {
    'Content-Type': 'application/json'
  }
})

export default api
