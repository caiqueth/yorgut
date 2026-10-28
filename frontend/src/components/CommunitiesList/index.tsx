import React from "react";
import { useNavigate } from "react-router-dom";
import { Box, Stack } from "@mui/material";
import { Button } from "components/Button";
import { BoxWithTitle, Count } from "components/BoxWithTitle";
import { StripedGrid } from "components/StripedGrid";
import { PicAndText, PicAndTextSkeleton } from "components/PicAndText";
import { useMeMemberships } from "api/me/queries";
import { useUserMemberships, useUserProfileSection } from "api/user/queries";
import { PrivateSection } from "components/PrivateSection";
import { FindUserMembershipsResponse } from "api/generated/model";
import { COMMUNITIES_BOX_SIZE } from "helpers/paging";

interface ICommunitiesListProps {
  // Whose communities to show; omit for the logged user ("minhas comunidades").
  userId?: number;
}

// The communities box of the home page (logged user) and of a profile page:
// a random few of their communities.
export const CommunitiesList: React.FC<ICommunitiesListProps> = ({ userId }) =>
  userId === undefined ? <MyCommunitiesList /> : <UserCommunitiesList userId={userId} />;

const MyCommunitiesList: React.FC = () => {
  const { data: memberships, isLoading, isFetching } = useMeMemberships(1, COMMUNITIES_BOX_SIZE, true);
  return (
    <CommunitiesGrid
      title={<>minhas comunidades <Count value={memberships?.count} /></>}
      memberships={memberships}
      isLoading={isLoading}
      isRefreshing={isFetching && !isLoading}
    />
  );
};

const UserCommunitiesList: React.FC<{ userId: number }> = ({ userId }) => {
  const section = useUserProfileSection(userId, "communities");
  const { data: memberships, isLoading, isFetching } = useUserMemberships(userId, 1, COMMUNITIES_BOX_SIZE, true, section.canSee);

  if (!section.isLoading && !section.canSee) {
    return (
      <BoxWithTitle title={"comunidades"}>
        <PrivateSection firstName={section.firstName} what={"as comunidades"} />
      </BoxWithTitle>
    );
  }
  return (
    <CommunitiesGrid
      title={<>comunidades <Count value={memberships?.count} /></>}
      memberships={memberships}
      // The request waits for the profile, which says whether it may be made.
      isLoading={section.isLoading || isLoading}
      isRefreshing={isFetching && !isLoading}
    />
  );
};

interface ICommunitiesGridProps {
  title: React.ReactNode;
  memberships?: FindUserMembershipsResponse;
  isLoading: boolean;
  isRefreshing: boolean;
}

const CommunitiesGrid: React.FC<ICommunitiesGridProps> = ({ title, memberships, isLoading, isRefreshing }) => {
  const navigate = useNavigate();

  const items =
    memberships?.communities.map((c) => (
      <PicAndText
        imageSrc={c.profilePicture}
        url={`/community/${c.id}`}
        title={c.name}
        titleSuffix={` (${c.membershipCount})`}
      />
    )) ?? [];

  return (
    <BoxWithTitle title={title} showProgress={isRefreshing}>
      <Stack direction={"column"} alignItems={"center"}>
        <StripedGrid items={items} elementsPerRow={3} length={COMMUNITIES_BOX_SIZE} isLoading={isLoading} skeleton={<PicAndTextSkeleton />} />
        <Box alignSelf={"center"}>
          <Button onClick={() => navigate("communities")}>ver comunidades</Button>
        </Box>
      </Stack>
    </BoxWithTitle>
  );
};
