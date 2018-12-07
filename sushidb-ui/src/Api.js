import React from "react";

export const API_BASE = process.env.NODE_ENV === "development" ? "/api" : "";

export function queryString(query = {}) {
  if (query instanceof Object) {
    const str = Object.entries(query)
      .map(param => `${param[0]}=${param[1]}`)
      .join("&");
    return str ? `?${str}` : "";
  }
  return "";
}

export function fetchKeys() {
  return fetch(`${API_BASE}/keys`);
}

export function fetchSingleMetric(metricId, options = {}) {
  return fetch(`${API_BASE}/metric/single/${metricId}${queryString(options)}`);
}

export function fetchMessageMetric(metricId, options = {}) {
  return fetch(`${API_BASE}/metric/message/${metricId}${queryString(options)}`);
}

export function useResource(fn, defaultValue, dependency = []) {
  const [body, setBody] = React.useState(defaultValue);

  async function refresh() {
    setBody(defaultValue);
    const res = await fn();
    const json = await res.json();
    setBody(json);
  }

  React.useEffect(() => {
    refresh();
    return () => {};
  }, dependency);

  return {
    body,
    refresh
  };
}
