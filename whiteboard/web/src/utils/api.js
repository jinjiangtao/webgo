const BASE_URL = '/api'

function _request(url, options = {}) {
  const fullUrl = BASE_URL + url
  const config = {
    headers: {
      'Content-Type': 'application/json'
    },
    ...options
  }
  if (config.body && typeof config.body !== 'string') {
    config.body = JSON.stringify(config.body)
  }
  return fetch(fullUrl, config).then((res) => {
    if (!res.ok) {
      throw new Error(`HTTP ${res.status}: ${res.statusText}`)
    }
    if (res.status === 204) {
      return null
    }
    return res.json()
  })
}

export function createWhiteboard(name, background = '#ffffff') {
  return _request('/whiteboards', {
    method: 'POST',
    body: { name, background }
  })
}

export function listWhiteboards() {
  return _request('/whiteboards', {
    method: 'GET'
  })
}

export function getWhiteboard(id) {
  return _request(`/whiteboards/${id}`, {
    method: 'GET'
  })
}

export function updateWhiteboard(id, name, background) {
  return _request(`/whiteboards/${id}`, {
    method: 'PUT',
    body: { name, background }
  })
}

export function deleteWhiteboard(id) {
  return _request(`/whiteboards/${id}`, {
    method: 'DELETE'
  })
}

export function saveOperations(id, operations) {
  return _request(`/whiteboards/${id}/operations`, {
    method: 'POST',
    body: { operations }
  })
}

export function createSnapshot(id, name, operations) {
  return _request(`/whiteboards/${id}/snapshots`, {
    method: 'POST',
    body: { name, operations }
  })
}

export function listSnapshots(id) {
  return _request(`/whiteboards/${id}/snapshots`, {
    method: 'GET'
  })
}

export function clearWhiteboard(id) {
  return _request(`/whiteboards/${id}/clear`, {
    method: 'POST'
  })
}
