import React, { useState, useEffect } from "react";

import { queryMessage, querySingle, useResource } from "../Api";
import { dateFormat } from "../Formatter";
import MonacoEditor from "react-monaco-editor";
import * as monaco from "monaco-editor/esm/vs/editor/editor.api";

const styles = {
  td: {
    verticalAlign: "middle"
  }
};

function editorDidMount(editor) {
  console.log(editor.getModel());
  monaco.languages.json.jsonDefaults.setDiagnosticsOptions({
    validate: true,
    schemas: [
      {
        uri: "https://github.com/kamijin-fanta/sushidb/schema/root.json",
        fileMatch: ["*"],
        schema: {
          type: "object",
          properties: {
            filters: {
              type: "array",
              items: {
                $ref:
                  "https://github.com/kamijin-fanta/sushidb/schema/filter.json"
              }
            },
            lower: { type: "integer" },
            upper: { type: "integer" },
            sort: { enum: ["asc", "desc"] },
            limit: { type: "integer" }
          }
        }
      },
      {
        uri: "https://github.com/kamijin-fanta/sushidb/schema/filter.json",
        schema: {
          anyOf: [
            {
              $ref:
                "https://github.com/kamijin-fanta/sushidb/schema/filter-leef.json"
            },
            {
              $ref:
                "https://github.com/kamijin-fanta/sushidb/schema/filter-has-child.json"
            }
          ]
        }
      },
      {
        uri: "https://github.com/kamijin-fanta/sushidb/schema/filter-leef.json",
        schema: {
          type: "object",
          properties: {
            type: {
              enum: ["eq", "gte", "gt", "lte", "lt"]
            },
            path: {
              type: "string"
            },
            value: {
              type: ["integer", "string"]
            }
          },
          additionalProperties: false,
          required: ["type", "path", "value"]
        }
      },
      {
        uri:
          "https://github.com/kamijin-fanta/sushidb/schema/filter-has-child.json",
        schema: {
          type: "object",
          properties: {
            type: {
              enum: ["and", "or"]
            },
            children: {
              type: "array",
              items: {
                $ref:
                  "https://github.com/kamijin-fanta/sushidb/schema/filter.json"
              }
            }
          },
          additionalProperties: false,
          required: ["type", "path", "value"]
        }
      }
    ]
  });
  monaco.languages.registerCompletionItemProvider("json", {
    provideCompletionItems: function(model, position) {
      var textUntilPosition = model.getValueInRange({
        startLineNumber: position.lineNumber,
        startColumn: 1,
        endLineNumber: position.lineNumber,
        endColumn: position.column
      });
      const suggestions = [];

      const dateTime = new Date();
      const ns = dateTime.valueOf() * 1000;
      if (
        textUntilPosition.includes("lower") ||
        textUntilPosition.includes("upper") ||
        textUntilPosition.includes("value")
      ) {
        suggestions.push({
          label: `0 current ns time: ${ns}`,
          kind: monaco.languages.CompletionItemKind.Value,
          documentation: "Current time stamp",
          insertText: ns.toString()
        });
      }

      return {
        suggestions
      };
    }
  });
}

export function QueryMetric(props) {
  const metricType = props.match.params.type;
  const metricKey = props.match.params.key;
  const initial = `{\n  "filters": [\n    \n  ]\n}`;
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
    setQuery(inputQuery);
    metrics.refresh();
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
      <MonacoEditor
        width="100%"
        height="300"
        language="json"
        theme="vs-dark"
        value={inputQuery}
        options={options}
        onChange={(newValue, e) => setInputQuery(newValue)}
        editorDidMount={editorDidMount}
      />
      <button onClick={onSubmit}>Search</button>
      <button onClick={onFormats}>Format</button>

      {metrics.isLoading ? (
        <>loading</>
      ) : (
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
                  <td style={styles.td}>
                    {JSON.stringify(row.value, null, 2)}
                  </td>
                </tr>
              ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
