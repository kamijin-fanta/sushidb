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

/********** SushiDB API **********/
export function fetchKeys() {
  return fetch(`${API_BASE}/keys`);
}

export function fetchSingleMetric(metricId, options = {}) {
  return fetch(`${API_BASE}/metric/single/${metricId}${queryString(options)}`);
}

export function fetchMessageMetric(metricId, options = {}) {
  return fetch(`${API_BASE}/metric/message/${metricId}${queryString(options)}`);
}

/********** PD API **********/
export function fetchPdList() {
  return fetch(`${API_BASE}/pd/`);
}
export function fetchStoreList() {
  return fetch(`${API_BASE}/pd/api/v1/stores`);
}

export function useResource(fn, defaultValue, dependency = []) {
  const [body, setBody] = React.useState(defaultValue);
  const [isLoading, setLoading] = React.useState(false);

  async function refresh() {
    setLoading(true);
    const res = await fn();
    const json = await res.json();
    setBody(json);
    setLoading(false);
  }
  async function clearAndRefresh() {
    setBody(defaultValue);
    await refresh();
  }

  React.useEffect(() => {
    clearAndRefresh();
    return () => {};
  }, dependency);

  return {
    body,
    refresh,
    isLoading
  };
}