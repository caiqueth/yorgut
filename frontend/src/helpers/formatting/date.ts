import dayjs, { Dayjs } from "dayjs";

export const dateFromText = (text: string, pattern: string): string =>
  dayjs(text).format(pattern);

export const dayjsFromDate = (d: Date, format: string) => {
  return dayjs(d.toString()).format(format);
};
