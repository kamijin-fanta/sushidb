import React, { useState } from "react";

import { queryMessage, querySingle, useResource } from "../Api";
import { dateFormat } from "../Formatter";
import MonacoEditor from "react-monaco-editor";

const styles = {
  td: {
    verticalAlign: "middle"
  }
};

export function QueryMetric(props) {
  const metricType = props.match.params.type;
  const metricKey = props.match.params.key;
  const initial = `{\n  "filters": []\n}`;
  const [inputQuery, setInputQuery] = useState(initial);
  const [query, setQuery] = useState(initial);
  const metrics = useResource(
    () =>
      metricType === "message"
        ? queryMessage(metricKey, query)
        : querySingle(metricKey, query),
    {},
    [metricType, metricKey, query]
  );

  const onSubmit = () => {
    onFormats();
    setQuery(inputQuery);
  };

  const onFormats = () => {
    setInputQuery(JSON.stringify(JSON.parse(inputQuery), null, 2));
  };

  const options = {
    selectOnLineNumbers: true,
    quickSuggestions: true,
    wordBasedSuggestions: true
  };

  return (
    <div className="page query-metric">
      <h1>Querying View</h1>
      {/* <textarea
        value={inputQuery}
        onChange={e => setInputQuery(e.target.value)}
        rows="10"
      /> */}
      <MonacoEditor
        width="100%"
        height="300"
        language="json"
        theme="vs-dark"
        value={inputQuery}
        options={options}
        onChange={(newValue, e) => setInputQuery(newValue)}
        editorDidMount={() => {}}
      />
      <button onClick={onSubmit}>Search</button>
      <button onClick={onFormats}>Format</button>

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
