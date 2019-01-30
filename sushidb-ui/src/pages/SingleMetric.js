import React, { useMemo } from "react";

import { fetchSingleMetric, useResource } from "../Api";
import { dateFormat } from "../Formatter";

import { ParentSize } from "@vx/responsive";
import MetricGraph from "../component/MetricGraph";

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
  let sorted = useMemo(
    () =>
      metrics.body.rows &&
      metrics.body.rows.concat().sort((a, b) => a.time - b.time),
    [metrics]
  );
  return (
    <div className="page single-metric">
      <h1>Single Metric View</h1>
      {sorted && (
        <ParentSize className="graph-container">
          {({ width: w, height: h }) => (
            <MetricGraph
              data={sorted}
              width={w}
              height={400}
              margin={{ left: 80, right: 30, top: 30, bottom: 50 }}
            />
          )}
        </ParentSize>
      )}
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
