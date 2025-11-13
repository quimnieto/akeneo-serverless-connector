import { useState, useEffect } from 'react';
import './App.css';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

function App() {
  const [subscribers, setSubscribers] = useState([]);
  const [subscriptions, setSubscriptions] = useState([]);
  const [eventTypes, setEventTypes] = useState([]);
  const [error, setError] = useState('');

  const [subscriberForm, setSubscriberForm] = useState({ 
    name: '', 
    contact: { technical_email: '' }
  });
  const [subscriptionForm, setSubscriptionForm] = useState({ 
    subscriber_id: '',
    source: 'pim',
    events: [],
    type: 'https',
    config: {
      url: '',
      secret: {
        primary: '',
        secondary: ''
      }
    }
  });

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
    try {
      const [subRes, subsRes] = await Promise.all([
        fetch(`${API_URL}/subscriber`).catch(() => ({ ok: false })),
        fetch(`${API_URL}/subscriptions`).catch(() => ({ ok: false }))
      ]);

      if (subRes.ok) {
        const subData = await subRes.json();
        setSubscribers(subData || []);
      }

      if (subsRes.ok) {
        const subsData = await subsRes.json();
        setSubscriptions(subsData || []);
      }
    } catch (err) {
      setError(err.message);
    }
  };

  const handleSubscriberSubmit = async (e) => {
    e.preventDefault();
    setError('');
    try {
      const res = await fetch(`${API_URL}/subscriber`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(subscriberForm)
      });

      if (!res.ok) throw new Error(await res.text());
      setSubscriberForm({ name: '', contact: { technical_email: '' } });
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
      setSubscriptionForm({ 
        subscriber_id: '',
        source: 'pim',
        events: [],
        type: 'https',
        config: {
          url: '',
          secret: {
            primary: '',
            secondary: ''
          }
        }
      });
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

  const deleteSubscription = async (id) => {
    if (!window.confirm('Delete this subscription?')) return;
    setError('');
    try {
      const res = await fetch(`${API_URL}/subscriptions/${id}`, { method: 'DELETE' });
      if (!res.ok) throw new Error(await res.text());
      await loadData();
    } catch (err) {
      setError(err.message);
    }
  };

  const deleteSubscriber = async (id) => {
    if (!window.confirm('Delete this subscriber?')) return;
    setError('');
    try {
      const res = await fetch(`${API_URL}/subscriber/${id}`, { method: 'DELETE' });
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
        <h2>Subscribers</h2>
        <form onSubmit={handleSubscriberSubmit}>
          <input
            type="text"
            placeholder="Subscriber Name"
            value={subscriberForm.name}
            onChange={(e) => setSubscriberForm({ ...subscriberForm, name: e.target.value })}
            required
          />
          <input
            type="email"
            placeholder="Technical Email"
            value={subscriberForm.contact.technical_email}
            onChange={(e) => setSubscriberForm({ 
              ...subscriberForm, 
              contact: { technical_email: e.target.value }
            })}
            required
          />
          <button type="submit">Create Subscriber</button>
        </form>

        <div className="list">
          {subscribers.length === 0 ? (
            <p>No subscribers</p>
          ) : (
            subscribers.map((sub) => (
              <div key={sub.id} className="item">
                <div>
                  <strong>{sub.name}</strong>
                  <span style={{ color: '#7f8c8d', fontSize: '0.9em', marginLeft: '10px' }}>
                    ID: {sub.id}
                  </span>
                </div>
                <div className="events">
                  Email: {sub.contact?.technical_email}
                </div>
                <div className="actions">
                  <button onClick={() => deleteSubscriber(sub.id)}>Delete</button>
                </div>
              </div>
            ))
          )}
        </div>
      </section>

      <section className="section">
        <h2>Subscriptions</h2>
        <form onSubmit={handleSubscriptionSubmit}>
          <select
            value={subscriptionForm.subscriber_id}
            onChange={(e) => setSubscriptionForm({ ...subscriptionForm, subscriber_id: e.target.value })}
            required
            style={{ padding: '10px', borderRadius: '4px', border: '1px solid #ddd', fontSize: '14px' }}
          >
            <option value="">Select Subscriber</option>
            {subscribers.map((sub) => (
              <option key={sub.id} value={sub.id}>
                {sub.name} ({sub.contact?.technical_email})
              </option>
            ))}
          </select>
          <input
            type="url"
            placeholder="Webhook URL"
            value={subscriptionForm.config.url}
            onChange={(e) => setSubscriptionForm({ 
              ...subscriptionForm, 
              config: { ...subscriptionForm.config, url: e.target.value }
            })}
            required
          />
          <input
            type="text"
            placeholder="Primary Secret Key"
            value={subscriptionForm.config.secret.primary}
            onChange={(e) => setSubscriptionForm({ 
              ...subscriptionForm, 
              config: { 
                ...subscriptionForm.config, 
                secret: { ...subscriptionForm.config.secret, primary: e.target.value }
              }
            })}
            required
          />
          <input
            type="text"
            placeholder="Secondary Secret Key (optional)"
            value={subscriptionForm.config.secret.secondary}
            onChange={(e) => setSubscriptionForm({ 
              ...subscriptionForm, 
              config: { 
                ...subscriptionForm.config, 
                secret: { ...subscriptionForm.config.secret, secondary: e.target.value }
              }
            })}
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
          <button type="submit">Create Subscription</button>
        </form>

        <div className="list">
          {subscriptions.length === 0 ? (
            <p>No subscriptions</p>
          ) : (
            subscriptions.map((sub) => (
              <div key={sub.id} className="item">
                <div>
                  <strong>{sub.config?.url || 'Webhook'}</strong>
                  <span style={{ color: '#7f8c8d', fontSize: '0.9em', marginLeft: '10px' }}>
                    Subscriber: {sub.subscriber_name}
                  </span>
                </div>
                <div className="events">{sub.events?.join(', ')}</div>
                <div className="actions">
                  <button onClick={() => deleteSubscription(sub.id)}>Delete</button>
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
