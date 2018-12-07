import React from "react";

import { Icon } from "@blueprintjs/core";
import { NavLink } from "react-router-dom";
import { fetchKeys, useResource } from "../Api";

const styles = {
  td: {
    verticalAlign: "middle"
  }
};

function genLink(key) {
  return key.type === "message"
    ? `/metric/message/${key.metric_id}`
    : `/metric/single/${key.metric_id}`;
}

export function Keys() {
  const keys = useResource(() => fetchKeys(), []);

  return (
    <div className="page keys">
      <h1>Metric Keys</h1>
      <table
        style={{ width: "100%" }}
        className="bp3-html-table bp3-html-table-condensed bp3-html-table-striped"
      >
        <thead>
          <tr>
            <th>Metric ID</th>
            <th>Type</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {keys.body.map(key => (
            <tr key={key.metric_id}>
              <td style={styles.td}>{key.metric_id}</td>
              <td style={styles.td}>{key.type}</td>
              <td style={styles.td}>
                <NavLink to={genLink(key)} className="bp3-button bp3-minimal">
                  <Icon icon="chart" />
                  <span className="bp3-button-text">View Metric</span>
                </NavLink>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <div>
        <button onClick={keys.refresh}>refresh</button>
        <pre>{JSON.stringify(keys.body, null, 2)}</pre>
      </div>
    </div>
  );
}