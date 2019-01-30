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

export function fetchSingleMetric(metricKey, options = {}) {
  return fetch(`${API_BASE}/metric/single/${metricKey}${queryString(options)}`);
}

export function fetchMessageMetric(metricKey, options = {}) {
  return fetch(
    `${API_BASE}/metric/message/${metricKey}${queryString(options)}`
  );
}

export function queryMessage(data = {}, options = {}) {
  return fetch(`${API_BASE}/query/message${queryString(options)}`, {
    method: "POST",
    body: data
  });
}

export function querySingle(data = "", options = {}) {
  return fetch(`${API_BASE}/query/single${queryString(options)}`, {
    method: "POST",
    body: data
  });
}

export function deleteMetric(type, metricKey, options = {}) {
  return fetch(
    `${API_BASE}/metric/${type}/${metricKey}${queryString(options)}`,
    {
      method: "DELETE"
    }
  );
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
  const [error, setError] = React.useState();
  const [isLoading, setLoading] = React.useState(false);

  async function refresh() {
    try {
      setError();
      setLoading(true);
      const res = await fn();
      const json = await res.json();
      setBody(json);
      setLoading(false);
    } catch (e) {
      setError(e);
    }
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
    isLoading,
    error
  };
}
