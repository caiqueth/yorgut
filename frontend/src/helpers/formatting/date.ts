import dayjs from "dayjs";

export const dateFromText = (text: string, pattern: string): string =>
  dayjs(text).format(pattern);

export const dayjsFromDate = (d: Date | string, format: string) => {
  if (!d) return null;
  return dayjs(d.toString()).format(format);
};
