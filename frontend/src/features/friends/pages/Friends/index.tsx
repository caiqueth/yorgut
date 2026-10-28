import React from "react";
import { Box, Skeleton, Stack, Typography } from "@mui/material";
import { Link } from "react-router-dom";
import { useMeFriends } from "api/me/queries";
import foreverAlone from "assets/images/foreveralone.png";
import { PagedList } from "components/PagedList";
import { PicAndTextHorizontal, PicAndTextHorizontalSkeleton } from "components/PicAndTextHorizontal";
import { FriendRating } from "features/friends/components/FriendRating";
import { LIST_PAGE_SIZE } from "helpers/paging";
import { usePage } from "hooks/usePage";

// The logged user's friends, each with their ratings and fan star.
export const Friends: React.FC = () => {
  const [page, setPage] = usePage();
  const { data, isLoading, isFetching } = useMeFriends(page, LIST_PAGE_SIZE);

  const items =
    data?.friends.map((f) => (
      <Stack
        key={f.id}
        direction={"row"}
        justifyContent={"space-between"}
        alignItems={"center"}
        gap={2}
        flexWrap={"wrap"}
      >
        <PicAndTextHorizontal
          imageSrc={f.profilePicture}
          title={f.firstName}
          url={`/profile/${f.id}`}
          imageSize={64}
        />
        <FriendRating userId={f.id} myReview={f.myReview} amIFan={f.amIFan} />
      </Stack>
    )) ?? [];

  return (
    // This page sits directly in the app layout, so it has to claim the full width.
    <Box width={"100%"}>
      <PagedList
        title={"meus amigos"}
        count={data?.count}
        items={items}
        isLoading={isLoading}
        isFetching={isFetching}
        skeleton={
          <Stack direction={"row"} justifyContent={"space-between"} alignItems={"center"} gap={2}>
            <PicAndTextHorizontalSkeleton imageSize={64} lineWidths={[96]} />
            <Typography variant={"body2"}><Skeleton variant={"text"} width={320} /></Typography>
          </Stack>
        }
        empty={
          <>
            <Typography>Parece que você ainda não possui amigos.</Typography>
            <span><img src={foreverAlone} /></span>
            <Typography><Link to="/search">Clique aqui</Link> e procure por conexões.</Typography>
          </>
        }
        page={page}
        onPageChange={setPage}
      />
    </Box>
  );
};
