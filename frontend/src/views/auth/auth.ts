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
  
  try {
    const res = await fetch(
      `${API_URL}auth/login`,
      {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: formData,
      },
    )

    const data = await res.json()
    if (data.code !== '0' && data.code !== '200' && data.code !== 0 && data.code !== 200) {
      throw new Error(data.message || 'Login failed')
    }
    if (!data.data || !data.data.access_token) {
      throw new Error('No access token in response')
    }
    localStorage.setItem('token', data.data.access_token)

    return data
  } catch (error: any) {
    // Distinguish between network error and API error
    if (error instanceof TypeError) {
      throw new Error(`Network error: ${error.message}. Check CORS settings on backend.`)
    }
    throw error
  }
}