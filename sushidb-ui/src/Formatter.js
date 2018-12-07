import dateTime from "date-and-time";

export function dateFormat(date) {
  return dateTime.format(date, "YYYY/MM/DD HH:mm:ss:SSS");
}
