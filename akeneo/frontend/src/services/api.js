const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

export const api = {
  getSubscriber: () => fetch(`${API_URL}/subscriber`),
  createSubscriber: (data) => fetch(`${API_URL}/subscriber`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  }),
  updateSubscriber: (data) => fetch(`${API_URL}/subscriber`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  }),
  getSubscriptions: () => fetch(`${API_URL}/subscriptions`),
  createSubscription: (data) => fetch(`${API_URL}/subscriptions`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  }),
  updateSubscription: (code, data) => fetch(`${API_URL}/subscriptions/${code}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  }),
  deleteSubscription: (code) => fetch(`${API_URL}/subscriptions/${code}`, {
    method: 'DELETE'
  })
};
