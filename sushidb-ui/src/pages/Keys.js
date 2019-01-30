import React from "react";

import { Icon, Button } from "@blueprintjs/core";
import { NavLink } from "react-router-dom";
import { fetchKeys, useResource, deleteMetric } from "../Api";

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
function genQueryLink(key) {
  return key.type === "message"
    ? `/query/message/${key.metric_id}`
    : `/query/single/${key.metric_id}`;
}

export function Keys() {
  const keys = useResource(() => fetchKeys(), []);

  async function onDeleteClick(key) {
    if (
      window.confirm(`DELETE THIS KEY \ntype:${key.type} id:${key.metric_id}`)
    ) {
      await deleteMetric(key.type, key.metric_id);
      keys.refresh();
    }
  }

  return (
    <div className="page keys">
      <h1>Metric Keys</h1>
      <table
        style={{ width: "100%", tableLayout: "fixed" }}
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
                <NavLink
                  to={genQueryLink(key)}
                  className="bp3-button bp3-minimal"
                >
                  <Icon icon="search" />
                  <span className="bp3-button-text">Query Metric</span>
                </NavLink>
                <Button
                  minimal={true}
                  icon="delete"
                  onClick={() => onDeleteClick(key)}
                >
                  Delete
                </Button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <div>
        <button onClick={keys.refresh}>refresh</button>
      </div>
    </div>
  );
}
