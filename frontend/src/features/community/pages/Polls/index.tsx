import React, { useState } from "react";
import { Box, Skeleton, Typography } from "@mui/material";
import { Link, useNavigate } from "react-router-dom";
import { useCommunityPolls } from "api/poll/queries";
import { PollSummary } from "api/generated/model";
import { BoxWithTitle, Count } from "components/BoxWithTitle";
import { Button } from "components/Button";
import { ListPagination } from "components/ListPagination";
import { StripedGrid } from "components/StripedGrid";
import { MembersOnly } from "features/community/components/MembersOnly";
import { useForumAccess } from "features/community/useForumAccess";
import { pollStatusText } from "features/community/pollStatus";
import { FORUM_PAGE_SIZE } from "helpers/paging";
import { useIdParam } from "hooks/useIdParam";

const PollRow: React.FC<{ poll: PollSummary }> = ({ poll }) => (
  <Box sx={{ px: 1, py: 0.5 }}>
    <Typography>
      <Link to={`${poll.id}`}>{poll.question}</Link>
    </Typography>
    <Typography variant={"body2"}>
      por <Link to={`/profile/${poll.authorId}`}>{poll.authorFirstName}</Link>
      {" · "}{poll.voteCount} {poll.voteCount === 1 ? "voto" : "votos"}
      {" · "}{pollStatusText(poll)}
      {poll.voted && " · ✓ você votou"}
    </Typography>
  </Box>
);

const PollRowSkeleton: React.FC = () => (
  <Box sx={{ px: 1, py: 0.5 }}>
    <Skeleton variant={"text"} width={220} />
    <Skeleton variant={"text"} width={160} />
  </Box>
);

// A community's polls, newest first. Same access rule as the forum: anyone in
// a public community, only members in a private one.
export const Polls: React.FC = () => {
  const communityId = useIdParam("communityId");
  const navigate = useNavigate();
  const { canRead, isLoading: isLoadingAccess } = useForumAccess(communityId);
  const [page, setPage] = useState(1);
  const { data, isLoading: isLoadingData, isFetching } = useCommunityPolls(communityId, page, FORUM_PAGE_SIZE, canRead);
  // The request waits for the access check; both count as loading.
  const isLoading = isLoadingAccess || isLoadingData;

  if (!isLoadingAccess && !canRead) {
    return (
      <BoxWithTitle title={"enquetes"}>
        <MembersOnly communityId={communityId} action={"ver as enquetes desta comunidade privada"} />
      </BoxWithTitle>
    );
  }

  const items = (data?.polls ?? []).map((p) => <PollRow key={p.id} poll={p} />);

  return (
    <BoxWithTitle title={<>enquetes <Count value={data?.count} /></>} showProgress={isFetching && !isLoading}>
      {items.length === 0 && !isLoading ? (
        <Box sx={{ p: 1 }}>
          <Typography>Nenhuma enquete ainda.</Typography>
        </Box>
      ) : (
        <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} skeleton={<PollRowSkeleton />} />
      )}
      <ListPagination page={page} pageSize={FORUM_PAGE_SIZE} count={data?.count} onChange={setPage} />
      <MembersOnly communityId={communityId} action={"criar enquetes"}>
        <Button onClick={() => navigate(`/community/${communityId}/polls/new`)}>+ nova enquete</Button>
      </MembersOnly>
    </BoxWithTitle>
  );
};
