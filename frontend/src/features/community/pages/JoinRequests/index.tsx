import { useIdParam } from "hooks/useIdParam";
import { Box, Stack, Typography } from "@mui/material";
import { useCommunity, useCommunityJoinRequests, useReplyJoinRequest } from "api/community/queries";
import { useMe } from "api/me/queries";
import { Button } from "components/Button";
import { BoxWithTitle, Count } from "components/BoxWithTitle";
import { PicAndTextHorizontal, PicAndTextHorizontalSkeleton } from "components/PicAndTextHorizontal";
import { StripedGrid } from "components/StripedGrid";
import { dayjsFromDate } from "helpers/formatting/date";
import React from "react";
import { Navigate } from "react-router-dom";

// The owner's list of users asking to join their private community.
export const JoinRequests: React.FC = () => {
  const communityId = useIdParam("communityId");

  const { data: me } = useMe();
  const { data: community } = useCommunity(communityId);
  const isOwner = !!me && !!community && me.id === community.ownerId;
  const { data, isLoading: isLoadingRequests, isFetching } = useCommunityJoinRequests(communityId, isOwner && community.private);
  // The request waits until we know the user owns the community.
  const isLoading = !me || !community || isLoadingRequests;
  const { mutate: reply, isPending, variables: replying } = useReplyJoinRequest(communityId);

  // Only the owner of a private community has requests; the backend enforces it too.
  if (me && community && (!isOwner || !community.private)) {
    return <Navigate to={`/community/${communityId}`} replace />;
  }

  const items = (data?.requests ?? []).map((r) => (
    <Stack
      key={r.userId}
      direction={"row"}
      justifyContent={"space-between"}
      alignItems={"center"}
      gap={1}
      flexWrap={"wrap"}
    >
      <PicAndTextHorizontal
        title={`${r.firstName} ${r.lastName}`}
        text={`pediu em ${dayjsFromDate(r.createdAt, "D [de] MMMM [de] YYYY")}`}
        url={`/profile/${r.userId}`}
        imageSrc={r.profilePicture}
        imageSize={56}
      />
      <Stack direction={"row"} gap={1}>
        <Button
          onClick={() => reply({ userId: r.userId, accept: true })}
          loading={isPending && replying?.userId === r.userId && replying.accept}
        >
          aceitar
        </Button>
        <Button
          onClick={() => reply({ userId: r.userId, accept: false })}
          loading={isPending && replying?.userId === r.userId && !replying.accept}
        >
          recusar
        </Button>
      </Stack>
    </Stack>
  ));

  return (
    <BoxWithTitle title={<>pedidos para entrar <Count value={data?.requests.length} /></>} showProgress={isFetching && !isLoading}>
      {items.length === 0 && !isLoading ? (
        <Box sx={{ p: 1 }}>
          <Typography>Nenhum pedido pendente.</Typography>
        </Box>
      ) : (
        <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} skeleton={<PicAndTextHorizontalSkeleton imageSize={56} lineWidths={[140, 160]} />} />
      )}
    </BoxWithTitle>
  );
};
