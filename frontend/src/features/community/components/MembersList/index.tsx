import React from "react";
import { useCommunityMembers } from "api/community/queries";
import { Button } from "components/Button";
import { useNavigate } from "react-router-dom";
import { Box, Stack } from "@mui/material";
import { BoxWithTitle } from "components/BoxWithTitle";
import { StripedGrid } from "components/StripedGrid";
import { PicAndText } from "components/PicAndText";

interface IMembersListProps {
  communityId: number;
}
export const MembersList: React.FC<IMembersListProps> = ({ communityId }) => {
  const { data: members, isLoading } = useCommunityMembers(communityId);
  const navigate = useNavigate();

  const handleGoToMembers = () => {
    navigate("members");
  };

  const items =
    members?.members.map((m) => (
      <PicAndText
        imageSrc={m.profilePicture}
        url={`/profile/${m.id}`}
        text={`${m.firstName} (${m.friendCount})`}
      />
    )) ?? [];

  return (
    <BoxWithTitle title={`membros (${members?.memberCount})`}>
      <Stack direction={"column"} alignItems={"center"}>
        <StripedGrid items={items} elementsPerRow={4} isLoading={isLoading} />
        <Box alignSelf={"center"}>
          <Button onClick={handleGoToMembers}>ver membros</Button>
        </Box>
      </Stack>
    </BoxWithTitle>
  );
};
