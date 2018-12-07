import React from "react";

import { fetchStoreList, useResource } from "../Api";
import { ProgressBar, Button, Checkbox, FormGroup } from "@blueprintjs/core";

import "./StoreInfo.css";

export function StoreInfo(props) {
  const stores = useResource(() => fetchStoreList(), {}, []);

  const [autoRefresh, setAutoRefresh] = React.useState(true);
  const inverseAutoRefresh = React.useCallback(() =>
    setAutoRefresh(!autoRefresh)
  );
  React.useEffect(
    () => {
      if (autoRefresh) {
        const id = setInterval(() => stores.refresh(), 3000);
        return () => clearInterval(id);
      }
      return () => {};
    },
    [autoRefresh]
  );

  return (
    <div className="page store-info">
      <h1>TiVK Store Info</h1>

      <FormGroup>
        <Button icon="refresh" onClick={stores.refresh}>
          Refresh
        </Button>
        <Checkbox checked={autoRefresh} onChange={inverseAutoRefresh}>
          Refresh every 3s
        </Checkbox>
      </FormGroup>

      <div className={`progress${stores.isLoading ? " loading" : ""}`}>
        fetching...
        <ProgressBar value={0.7} />
      </div>

      <div className="results">
        {stores.body.stores &&
          stores.body.stores.map(info => (
            <div key={info.store.id}>
              <div className="header">
                <span className="address">{info.store.address}</span>
                <span className="state">{info.store.state_name}</span>
              </div>
              <div className="describe">
                <table>
                  <tbody>
                    <tr>
                      <th>ID</th>
                      <td>{info.store.id}</td>
                    </tr>
                    <tr>
                      <th>Version</th>
                      <td>{info.store.version}</td>
                    </tr>
                    <tr className="spacing" />
                    <tr>
                      <th>Disk</th>
                      <td>
                        {info.status.available} / {info.status.capacity}
                      </td>
                    </tr>
                    <tr>
                      <th>Leader Weight</th>
                      <td>{info.status.leader_weight}</td>
                    </tr>
                    <tr>
                      <th>Region</th>
                      <td>
                        count: {info.status.region_count} / weight:{" "}
                        {info.status.region_weight} / score:{" "}
                        {info.status.region_score} / size:{" "}
                        {info.status.region_size}
                      </td>
                    </tr>
                    <tr className="spacing" />
                    <tr>
                      <th>Start</th>
                      <td>{info.status.start_ts}</td>
                    </tr>
                    <tr>
                      <th>Heartbeat</th>
                      <td>{info.status.last_heartbeat_ts}</td>
                    </tr>
                    <tr>
                      <th>Uptime</th>
                      <td>{info.status.uptime}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          ))}
      </div>
    </div>
  );
}
