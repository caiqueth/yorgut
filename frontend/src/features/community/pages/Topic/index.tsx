import { useIdParam } from "hooks/useIdParam";
import { Box, Skeleton, Typography } from "@mui/material";
import { useTopicReplies } from "api/topic/queries";
import { PicAndTextHorizontal, PicAndTextHorizontalSkeleton } from "components/PicAndTextHorizontal";
import { BoxWithTitle, Count } from "components/BoxWithTitle";
import { StripedGrid } from "components/StripedGrid";
import { ListPagination } from "components/ListPagination";
import { FORUM_PAGE_SIZE } from "helpers/paging";
import { dayjsFromDate } from "helpers/formatting/date";
import React, { useState } from "react";
import { Reply } from "features/community/components/TopicReply";
import { MembersOnly } from "features/community/components/MembersOnly";
import { useForumAccess } from "features/community/useForumAccess";

export const Topic: React.FC = () => {
  const topicId = useIdParam("topicId");
  const communityId = useIdParam("communityId");

  const { canRead, isLoading: isLoadingAccess } = useForumAccess(communityId);
  const [page, setPage] = useState(1);
  const { data: replies, isLoading: isLoadingData, isFetching } = useTopicReplies(topicId, page, FORUM_PAGE_SIZE, canRead);
  // The request waits for the access check; both count as loading.
  const isLoading = isLoadingAccess || isLoadingData;

  // Show the new reply: it is the last one, on the last page.
  const handleReplied = () =>
    setPage(Math.max(1, Math.ceil(((replies?.count ?? 0) + 1) / FORUM_PAGE_SIZE)));

  // The topic itself opens the first page, followed by its replies.
  const posts = replies
    ? [...(replies.topic && page === 1 ? [replies.topic] : []), ...(replies.replies ?? [])]
    : [];
  const items = posts.map((r, idx) => (
      <PicAndTextHorizontal
        key={idx}
        title={r.authorFirstName}
        titleSuffix={` em ${dayjsFromDate(r.createdAt, "D [de] MMMM [de] YYYY [às] HH[h]mm")}:`}
        text={r.text}
        url={`/profile/${r.authorId}`}
        imageSrc={r.authorProfilePicture}
        imageSize={56}
      />
    ));

  if (!isLoadingAccess && !canRead) {
    return (
      <BoxWithTitle title="respostas">
        <MembersOnly communityId={communityId} action={"ver o fórum desta comunidade privada"} />
      </BoxWithTitle>
    );
  }

  return (
    <Box>
      <BoxWithTitle title={<>respostas <Count value={replies?.count} /></>} showProgress={isFetching && !isLoading}>
        <Box sx={{ px: 1 }}>
          <Typography variant={"h6"}>{replies ? replies.topic?.title : <Skeleton variant={"text"} width={"40%"} />}</Typography>
        </Box>
        <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} skeleton={<PicAndTextHorizontalSkeleton imageSize={56} lineWidths={[220, 260]} />} />
        <ListPagination page={page} pageSize={FORUM_PAGE_SIZE} count={replies?.count} onChange={setPage} />
        <Reply topicId={topicId} onReplied={handleReplied} />
      </BoxWithTitle>
    </Box>
  )
}