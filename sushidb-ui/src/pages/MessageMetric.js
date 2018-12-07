import React from "react";

import { fetchMessageMetric, useResource } from "../Api";
import { dateFormat } from "../Formatter";

const styles = {
  td: {
    verticalAlign: "middle"
  }
};

export function MessageMetric(props) {
  const metricKey = props.match.params.key;
  const metrics = useResource(() => fetchMessageMetric(metricKey), {}, [
    metricKey
  ]);

  return (
    <div className="page message-metric">
      <h1>Message Metric View</h1>
      <table
        style={{ width: "100%" }}
        className="bp3-html-table bp3-html-table-condensed bp3-html-table-striped"
      >
        <thead>
          <tr>
            <th>Metric ID</th>
            <th>Time</th>
            <th>Value</th>
          </tr>
        </thead>
        <tbody>
          {metrics.body.rows &&
            metrics.body.rows.map(row => (
              <tr key={row.time}>
                <td style={styles.td}>{metricKey}</td>
                <td style={styles.td}>
                  {dateFormat(new Date(row.time / 1000))}
                </td>
                <td style={styles.td}>{JSON.stringify(row.value, null, 2)}</td>
              </tr>
            ))}
        </tbody>
      </table>
    </div>
  );
}
