import React from "react";

import { fetchSingleMetric, useResource } from "../Api";
import { dateFormat } from "../Formatter";

const styles = {
  td: {
    verticalAlign: "middle"
  }
};

export function SingleMetric(props) {
  const metricKey = props.match.params.key;
  const metrics = useResource(() => fetchSingleMetric(metricKey), {}, [
    metricKey
  ]);

  return (
    <div className="page single-metric">
      <h1>Single Metric View</h1>
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
                <td style={styles.td}>{row.value}</td>
              </tr>
            ))}
        </tbody>
      </table>
    </div>
  );
}
