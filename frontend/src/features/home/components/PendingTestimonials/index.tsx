import { Stack } from "@mui/material";
import { useMe } from "api/me/queries";
import {
  useApproveTestimonial,
  useDeleteTestimonial,
  useTestimonials,
} from "api/testimonial/queries";
import { BoxWithTitle } from "components/BoxWithTitle";
import { Button } from "components/Button";
import { PicAndTextHorizontal } from "components/PicAndTextHorizontal";
import { StripedGrid } from "components/StripedGrid";
import { dayjsFromDate } from "helpers/formatting/date";
import React from "react";

// Testimonials written to the logged user that wait for their approval.
// Hidden when there are none.
export const PendingTestimonials: React.FC = () => {
  const { data: me } = useMe();
  const userId = me?.id ?? 0;
  // The API sends the pending testimonials only to their recipient, next to
  // a page of approved ones; one approved is the smallest page.
  const { data, isFetching } = useTestimonials(userId, 1, 1);
  const { mutate: approve, isPending: isApproving, variables: approving } = useApproveTestimonial(userId);
  const { mutate: reject, isPending: isRejecting, variables: rejecting } = useDeleteTestimonial(userId);

  const pending = data?.pending ?? [];
  if (pending.length === 0) return null;

  const items = pending.map((t) => (
    <Stack key={t.id} direction={"column"} gap={1}>
      <PicAndTextHorizontal
        imageSrc={t.authorProfilePicture}
        title={t.authorName}
        titleSuffix={` em ${dayjsFromDate(t.createdAt, "D [de] MMMM [de] YYYY")}:`}
        text={t.content}
        url={`/profile/${t.authorId}`}
      />
      <Stack direction={"row"}>
        <Button onClick={() => approve(t.id)} loading={isApproving && approving === t.id}>
          aceitar
        </Button>
        <Button onClick={() => reject(t.id)} loading={isRejecting && rejecting === t.id}>
          recusar
        </Button>
      </Stack>
    </Stack>
  ));

  return (
    <BoxWithTitle title={`depoimentos aguardando aprovação (${pending.length})`} showProgress={isFetching}>
      <StripedGrid items={items} elementsPerRow={1} />
    </BoxWithTitle>
  );
};
