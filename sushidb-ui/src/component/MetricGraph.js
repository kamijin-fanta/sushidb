import React from "react";

import { AreaClosed, Line, Bar } from "@vx/shape";
import { curveMonotoneX } from "@vx/curve";
import { GridRows, GridColumns } from "@vx/grid";
import { Group } from "@vx/group";
import { scaleLinear } from "@vx/scale";
import { withTooltip, TooltipWithBounds } from "@vx/tooltip";
import { localPoint } from "@vx/event";
import { extent, max, bisector } from "d3-array";
import { timeFormat } from "d3-time-format";
import { AxisLeft, AxisBottom } from "@vx/axis";
import { GlyphDot } from "@vx/glyph";

const formatDate = timeFormat("%b %d, '%y");
// const min = (arr, fn) => Math.min(...arr.map(fn));
// const max = (arr, fn) => Math.max(...arr.map(fn));
// const extent = (arr, fn) => [min(arr, fn), max(arr, fn)];

// accessors
const xSelector = d => d.time;
const ySelector = d => d.value;
const bisectDate = bisector(xSelector).left;

class Area extends React.Component {
  constructor(props) {
    super(props);
    this.handleTooltip = this.handleTooltip.bind(this);
  }
  handleTooltip({ margin, event, data, xSelector, xScale, yScale }) {
    const { showTooltip } = this.props;
    const { x } = localPoint(event);
    const time = xScale.invert(x - margin.left);
    const index = bisectDate(data, time, 1);
    const d0 = data[index - 1];
    const d1 = data[index];
    const d = d1 && time - xSelector(d0) > xSelector(d1) - time ? d1 : d0;
    const fixedX = xScale(xSelector(d));
    showTooltip({
      tooltipData: d,
      tooltipLeft: fixedX + margin.left,
      tooltipTop: yScale(ySelector(d))
    });
  }
  render() {
    const {
      width,
      height,
      margin,
      hideTooltip,
      tooltipData,
      tooltipTop,
      tooltipLeft,
      data
    } = this.props;
    if (width < 10) return null;

    // bounds
    const xMax = width - margin.left - margin.right;
    const yMax = height - margin.top - margin.bottom;

    // scales
    const xScale = scaleLinear({
      range: [0, xMax],
      domain: extent(data, xSelector)
    });
    console.log('extent(data, xSelector)', extent(data, xSelector))
    const yScale = scaleLinear({
      range: [yMax, 0],
      domain: [0, max(data, ySelector)],
      nice: true
    });

    return (
      <div>
        <svg ref={s => (this.svg = s)} width={width} height={height}>
          <rect
            x={0}
            y={0}
            width={width}
            height={height}
            fill="#4b6496"
            rx={14}
          />
          <defs>
            <linearGradient id="gradient" x1="0%" y1="0%" x2="0%" y2="100%">
              <stop offset="0%" stopColor="#FFFFFF" stopOpacity={1} />
              <stop offset="100%" stopColor="#FFFFFF" stopOpacity={0.2} />
            </linearGradient>
          </defs>
          <Group top={margin.top} left={margin.left}>
            <GridRows
              numTicks={5}
              lineStyle={{ pointerEvents: "none" }}
              scale={yScale}
              width={xMax}
              strokeDasharray="2,2"
              stroke="rgba(255,255,255,0.3)"
            />
            <GridColumns
              lineStyle={{ pointerEvents: "none" }}
              scale={xScale}
              height={yMax}
              strokeDasharray="2,2"
              stroke="rgba(255,255,255,0.3)"
            />
            <AxisLeft
              top={0}
              left={0}
              scale={yScale}
              numTicks={5}
              label="Value"
              labelProps={{
                fill: "#fff",
                textAnchor: "middle",
                fontSize: 12
              }}
              stroke="#fff"
              tickStroke="#fff"
              tickLabelProps={(value, index) => ({
                fill: "#fff",
                textAnchor: "end",
                fontSize: 10,
                dx: "-0.25em",
                dy: "0.25em"
              })}
              tickComponent={({ formattedValue, ...tickProps }) => (
                <text {...tickProps}>{formattedValue}</text>
              )}
            />
            <AxisBottom
              top={yMax}
              left={0}
              scale={xScale}
              label="Time"
              labelProps={{
                fill: "#fff",
                textAnchor: "middle",
                fontSize: 12
              }}
              stroke="#fff"
              tickStroke="#fff"
              tickLabelProps={(value, index) => ({
                fill: "#fff",
                textAnchor: "middle",
                fontSize: 10,
                dx: "-0.25em",
                dy: "0.25em"
              })}
              tickComponent={({ formattedValue, ...tickProps }) => (
                <text {...tickProps}>{formattedValue}</text>
              )}
            />
            <AreaClosed
              data={data}
              x={d => xScale(xSelector(d))}
              y={d => yScale(ySelector(d))}
              yScale={yScale}
              strokeWidth={1}
              stroke={"url(#gradient)"}
              fill={"url(#gradient)"}
              curve={curveMonotoneX}
            />
            <Group>
              {data.map((d, i) => {
                const x = xScale(xSelector(d));
                const y = yScale(ySelector(d));
                return (
                  <GlyphDot
                    key={i}
                    cx={x}
                    cy={y}
                    r={4}
                    fill="#fff"
                    stroke="#4b6496"
                    strokeWidth={1}
                  />
                );
              })}
            </Group>
          </Group>
          <Bar
            x={margin.left}
            y={margin.top}
            width={xMax}
            height={yMax}
            fill="transparent"
            rx={14}
            data={data}
            onTouchStart={event =>
              this.handleTooltip({
                event,
                xSelector,
                xScale,
                yScale,
                margin,
                data: data
              })
            }
            onTouchMove={event =>
              this.handleTooltip({
                event,
                xSelector,
                xScale,
                yScale,
                margin,
                data: data
              })
            }
            onMouseMove={event =>
              this.handleTooltip({
                event,
                xSelector,
                xScale,
                yScale,
                margin,
                data: data
              })
            }
            onMouseLeave={event => hideTooltip()}
          />
          {tooltipData && (
            <Group top={margin.top}>
              <Line
                from={{ x: tooltipLeft, y: 0 }}
                to={{ x: tooltipLeft, y: yMax }}
                stroke="rgba(92, 119, 235, 1.000)"
                strokeWidth={2}
                style={{ pointerEvents: "none" }}
                strokeDasharray="2,2"
              />
              <circle
                cx={tooltipLeft}
                cy={tooltipTop + 1}
                r={4}
                fill="black"
                fillOpacity={0.1}
                stroke="black"
                strokeOpacity={0.1}
                strokeWidth={2}
                style={{ pointerEvents: "none" }}
              />
              <circle
                cx={tooltipLeft}
                cy={tooltipTop}
                r={4}
                fill="rgba(92, 119, 235, 1.000)"
                stroke="white"
                strokeWidth={2}
                style={{ pointerEvents: "none" }}
              />
            </Group>
          )}
        </svg>
        {tooltipData && (
          <div>
            tooltipTop: {tooltipTop} / tooltipLeft: {tooltipLeft} / tooltipData:{" "}
            {JSON.stringify(tooltipData)}
            <TooltipWithBounds
              top={tooltipTop - 12}
              left={tooltipLeft + 12}
              style={{
                backgroundColor: "rgba(92, 119, 235, 1.000)",
                color: "white"
              }}
            >
              {`${ySelector(tooltipData)}`}
              <br />
              {formatDate(xSelector(tooltipData) / 1000)}
            </TooltipWithBounds>
          </div>
        )}
      </div>
    );
  }
}

export default withTooltip(Area);
