const API_BASE = '/api'

const request = async (url, options = {}) => {
  try {
    const response = await fetch(`${API_BASE}${url}`, {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers
      },
      ...options
    })
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    return await response.json()
  } catch (error) {
    console.error(`API request failed: ${url}`, error)
    throw error
  }
}

export const uploadImage = async (file) => {
  const formData = new FormData()
  formData.append('file', file)
  try {
    const response = await fetch(`${API_BASE}/images/upload`, {
      method: 'POST',
      body: formData
    })
    if (!response.ok) {
      throw new Error(`Upload failed: ${response.status}`)
    }
    return await response.json()
  } catch (error) {
    console.error('Image upload failed:', error)
    throw error
  }
}

export const uploadMultipleImages = async (files) => {
  const results = []
  for (const file of files) {
    try {
      const result = await uploadImage(file)
      results.push(result)
    } catch (error) {
      console.error(`Failed to upload ${file.name}:`, error)
    }
  }
  return results
}

export const saveEditLog = async (imageUuid, actionType, actionData = {}) => {
  return request('/edit-logs', {
    method: 'POST',
    body: JSON.stringify({
      image_uuid: imageUuid,
      action_type: actionType,
      action_data: actionData
    })
  })
}

export const getEditLogs = async (imageUuid) => {
  return request(`/edit-logs/${imageUuid}`)
}

export const getImages = async () => {
  return request('/images')
}

export const getImage = async (uuid) => {
  return request(`/images/${uuid}`, {
    headers: {
      'Accept': 'application/json'
    }
  })
}

export const getImageUrl = (uuid) => {
  return `/api/images/${uuid}`
}

export const deleteImage = async (uuid) => {
  return request(`/images/${uuid}`, {
    method: 'DELETE'
  })
}
