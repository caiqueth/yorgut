import React from "react";
import { useCommunityMembers } from "api/community/queries";
import { Button } from "components/Button";
import { useNavigate } from "react-router-dom";
import { Box, Stack } from "@mui/material";
import { BoxWithTitle, Count } from "components/BoxWithTitle";
import { StripedGrid } from "components/StripedGrid";
import { PicAndText, PicAndTextSkeleton } from "components/PicAndText";
import { BOX_SIZE } from "helpers/paging";

interface IMembersListProps {
  communityId: number;
}
export const MembersList: React.FC<IMembersListProps> = ({ communityId }) => {
  const { data: members, isLoading, isFetching } = useCommunityMembers(communityId, 1, BOX_SIZE);
  const navigate = useNavigate();

  const handleGoToMembers = () => {
    navigate("members");
  };

  const items =
    members?.members.map((m) => (
      <PicAndText
        imageSrc={m.profilePicture}
        url={`/profile/${m.id}`}
        title={m.firstName}
        titleSuffix={` (${m.friendCount})`}
      />
    )) ?? [];

  return (
    <BoxWithTitle title={<>membros <Count value={members?.count} /></>} showProgress={isFetching && !isLoading}>
      <Stack direction={"column"} alignItems={"center"}>
        <StripedGrid items={items} elementsPerRow={4} isLoading={isLoading} skeleton={<PicAndTextSkeleton />} />
        <Box alignSelf={"center"}>
          <Button onClick={handleGoToMembers}>ver membros</Button>
        </Box>
      </Stack>
    </BoxWithTitle>
  );
};
