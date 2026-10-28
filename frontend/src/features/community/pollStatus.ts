import { dayjsFromDate } from "helpers/formatting/date";

// "encerrada", "aberta até 20 de outubro às 18h00" or "aberta", for a poll's
// list line and its page.
export const pollStatusText = (poll: { closed: boolean; closesAt: string | null }) => {
  if (poll.closed) return "encerrada";
  if (poll.closesAt) return `aberta até ${dayjsFromDate(poll.closesAt, "D [de] MMMM [às] HH[h]mm")}`;
  return "aberta";
};
