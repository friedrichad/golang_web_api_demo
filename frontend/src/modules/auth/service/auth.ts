const API_URL = 'http://localhost:8080/'

export const loginApi = async (
  username: string,
  password: string,
  grant_type: string = 'password'
) => {
  const formData = new URLSearchParams()

  formData.append('username', username)
  formData.append('password', password)
  formData.append('grant_type', grant_type)
  const res = await fetch(
    `${API_URL}auth/login`,
    {
      method: 'POST',

      headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
      body: formData,
    },
  )

  if (!res.ok) {
    throw new Error('Login failed')
  }

  const data = await res.json()

  localStorage.setItem(
    'token',
    data.data.access_token
  )

  return data
}