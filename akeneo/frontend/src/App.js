import { useState, useEffect } from 'react';
import './App.css';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

function App() {
  const [subscriber, setSubscriber] = useState(null);
  const [subscriptions, setSubscriptions] = useState([]);
  const [eventTypes, setEventTypes] = useState([]);
  const [error, setError] = useState('');

  const [subscriberForm, setSubscriberForm] = useState({ url: '', active: true });
  const [subscriptionForm, setSubscriptionForm] = useState({ connection_code: '', events: [], active: true });

  useEffect(() => {
    loadData();
    loadEventTypes();
  }, []);

  const loadEventTypes = async () => {
    try {
      const res = await fetch(`${API_URL}/event-types`);
      if (res.ok) {
        const data = await res.json();
        setEventTypes(data || []);
      }
    } catch (err) {
      console.error('Failed to load event types:', err);
    }
  };

  const loadData = async () => {
    setLoading(true);
    try {
      const [subRes, subsRes] = await Promise.all([
        fetch(`${API_URL}/subscriber`).catch(() => ({ ok: false })),
        fetch(`${API_URL}/subscriptions`).catch(() => ({ ok: false }))
      ]);

      if (subRes.ok) {
        const subData = await subRes.json();
        setSubscriber(subData);
        setSubscriberForm({ url: subData.url || '', active: subData.active ?? true });
      }

      if (subsRes.ok) {
        const subsData = await subsRes.json();
        setSubscriptions(subsData || []);
      }
    } catch (err) {
      setError(err.message);
    }
    setLoading(false);
  };

  const handleSubscriberSubmit = async (e) => {
    e.preventDefault();
    setError('');
    try {
      const method = subscriber ? 'PATCH' : 'POST';
      const res = await fetch(`${API_URL}/subscriber`, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(subscriberForm)
      });

      if (!res.ok) throw new Error(await res.text());
      await loadData();
    } catch (err) {
      setError(err.message);
    }
  };

  const handleSubscriptionSubmit = async (e) => {
    e.preventDefault();
    setError('');
    if (subscriptionForm.events.length === 0) {
      setError('Please select at least one event type');
      return;
    }
    try {
      const res = await fetch(`${API_URL}/subscriptions`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(subscriptionForm)
      });

      if (!res.ok) throw new Error(await res.text());
      setSubscriptionForm({ connection_code: '', events: [], active: true });
      await loadData();
    } catch (err) {
      setError(err.message);
    }
  };

  const handleEventToggle = (event) => {
    setSubscriptionForm(prev => ({
      ...prev,
      events: prev.events.includes(event)
        ? prev.events.filter(e => e !== event)
        : [...prev.events, event]
    }));
  };

  const toggleSubscription = async (code, currentActive) => {
    setError('');
    try {
      const res = await fetch(`${API_URL}/subscriptions/${code}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ active: !currentActive })
      });

      if (!res.ok) throw new Error(await res.text());
      await loadData();
    } catch (err) {
      setError(err.message);
    }
  };

  const deleteSubscription = async (code) => {
    if (!window.confirm('Delete this subscription?')) return;
    setError('');
    try {
      const res = await fetch(`${API_URL}/subscriptions/${code}`, { method: 'DELETE' });
      if (!res.ok) throw new Error(await res.text());
      await loadData();
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div className="app">
      <h1>Akeneo Event Platform Config</h1>

      {error && <div className="error">{error}</div>}

      <section className="section">
        <h2>Subscriber</h2>
        <form onSubmit={handleSubscriberSubmit}>
          <input
            type="url"
            placeholder="Webhook URL"
            value={subscriberForm.url}
            onChange={(e) => setSubscriberForm({ ...subscriberForm, url: e.target.value })}
            required
          />
          <label>
            <input
              type="checkbox"
              checked={subscriberForm.active}
              onChange={(e) => setSubscriberForm({ ...subscriberForm, active: e.target.checked })}
            />
            Active
          </label>
          <button type="submit">{subscriber ? 'Update' : 'Create'} Subscriber</button>
        </form>
        {subscriber && (
          <div className="info">
            <p><strong>URL:</strong> {subscriber.url}</p>
            <p><strong>Status:</strong> {subscriber.active ? 'Active' : 'Inactive'}</p>
          </div>
        )}
      </section>

      <section className="section">
        <h2>Subscriptions</h2>
        <form onSubmit={handleSubscriptionSubmit}>
          <input
            type="text"
            placeholder="Connection Code"
            value={subscriptionForm.connection_code}
            onChange={(e) => setSubscriptionForm({ ...subscriptionForm, connection_code: e.target.value })}
            required
          />
          <div className="event-selector">
            <label>Select Events ({subscriptionForm.events.length} selected):</label>
            <div className="event-list">
              {eventTypes.length === 0 ? (
                <p>Loading event types...</p>
              ) : (
                eventTypes.map(event => (
                  <label key={event} className="event-item">
                    <input
                      type="checkbox"
                      checked={subscriptionForm.events.includes(event)}
                      onChange={() => handleEventToggle(event)}
                    />
                    {event}
                  </label>
                ))
              )}
            </div>
          </div>
          <label>
            <input
              type="checkbox"
              checked={subscriptionForm.active}
              onChange={(e) => setSubscriptionForm({ ...subscriptionForm, active: e.target.checked })}
            />
            Active
          </label>
          <button type="submit">Create Subscription</button>
        </form>

        <div className="list">
          {subscriptions.length === 0 ? (
            <p>No subscriptions</p>
          ) : (
            subscriptions.map((sub) => (
              <div key={sub.connection_code} className="item">
                <div>
                  <strong>{sub.connection_code}</strong>
                  <span className={sub.active ? 'active' : 'inactive'}>
                    {sub.active ? 'Active' : 'Inactive'}
                  </span>
                </div>
                <div className="events">{sub.events?.join(', ')}</div>
                <div className="actions">
                  <button onClick={() => toggleSubscription(sub.connection_code, sub.active)}>
                    {sub.active ? 'Deactivate' : 'Activate'}
                  </button>
                  <button onClick={() => deleteSubscription(sub.connection_code)}>Delete</button>
                </div>
              </div>
            ))
          )}
        </div>
      </section>
    </div>
  );
}

export default App;
